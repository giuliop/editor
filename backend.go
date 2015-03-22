package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"time"
	"unicode/utf8"
)

const defaultFileName = "newfile"

// backend holds the buffers open in the editor
type backend struct {
	bufs        []*buffer // the open buffers
	msgLine     line      // to hold messages to display to user
	commandMode bool      // wether we are in command mode
}

// initBackend returns the backend after having initialized it
func initBackend() backend {
	be := backend{}
	be.msgLine = line{}
	return be
}

type filetype int

const (
	any filetype = iota
	_go
)

var filetypes = map[string]filetype{
	".go": _go,
}

// open takes a list of filenames and open a buffer for each returning the
// first one as current buffer; if the list in empty it returns a new buffer.
// Non-existing filenames will also open new buffers
func (be *backend) open(filenames []string) *buffer {
	if len(filenames) == 0 {
		return be.newBuffer("")
	}
	i, fn := 0, ""
	for i, fn = range filenames {
		b := be.newBuffer("")
		err := be.openFile(b, fn)
		if err != nil {
			b.text[0] = line(fmt.Sprint(err))
		}
		be.bufs = append(be.bufs, b)
	}
	return be.bufs[len(be.bufs)-i-1]
}

// newBuffer adds a new empty buffer to the backend and returns a pointer to it
// Note that the last line of a buffer ends with a newline which is removed before
// saving to file
func (be *backend) newBuffer(name string) *buffer {
	if name == "" {
		name = defaultFileName
	}
	b := &buffer{
		text:       make([]line, 1, 20),
		name:       name,
		filename:   filename(name),
		changeList: changeList{ops: make([]bufferChange, 1)},
	}
	newMark(b).initLastInsert()
	b.text[0] = newLine()
	be.bufs = append(be.bufs, b)
	return b
}

// reopen refresh the buffer contenct from the file, useful if an external command
// changed the file
func (b *buffer) reopen() {
	err := be.openFile(b, b.filename)
	if err != nil {
		debug.Printf("buffer reopen error: %v", err)
	}
}

// save saves the buffer
func (b *buffer) save() error {
	return b.saveAs(b.filename)
}

// saveAs saves the buffer in a file named after the passed parameter
func (b *buffer) saveAs(filename string) error {
	for _, h := range beforeSaveHooks[any] {
		h(b)
	}
	for _, h := range beforeSaveHooks[b.filetype] {
		h(b)
	}

	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	br := &bufReader{b, 0, 0}
	_, err = io.Copy(f, br)
	if err == nil {
		b.fileSync = time.Now().UTC()
	}
	return err
}

// bufReader is used to implement the Reader interface and help copy
// buffers to files
type bufReader struct {
	buf  *buffer
	line int
	pos  int
}

// Read implements the Reader interface for bufReader
func (br *bufReader) Read(p []byte) (n int, err error) {
	b := br.buf
	for linenum, ln := range b.text[br.line:] {
		for posnum, r := range ln[br.pos:] {
			if len(p) < n+utf8.UTFMax {
				br.line += linenum
				br.pos += posnum
				return n, nil
			}
			n += utf8.EncodeRune(p[n:], r)
		}
		br.pos = 0
	}
	return n, io.EOF
}

func filename(name string) string {
	fp, err := filepath.Abs(name)
	if err != nil {
		fatalError(err)
	}
	return fp
}

// openFile opens the file in filename and adds it to passed in buffer b
func (be *backend) openFile(b *buffer, name string) error {
	fp := filename(name)
	b.filename = fp
	b.name = path.Base(fp)
	b.filetype = filetypes[path.Ext(fp)]
	errPrefix := "Hmpf, I cannot open the file '%v', "
	f, err := os.Open(fp)
	switch {
	case os.IsNotExist(err):
		return nil
	case os.IsPermission(err):
		return fmt.Errorf(errPrefix+"we do not have access rights", fp)
	case err != nil:
		return fmt.Errorf(errPrefix+"got this error:\n%v\n", fp, err)
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	b.text = b.text[:0]
	b.mod = normalMode
	for sc.Scan() {
		b.text = append(b.text, []rune(sc.Text()+"\n"))
	}
	if err := sc.Err(); err != nil {
		return fmt.Errorf(errPrefix+"got this error:\n%v\n", fp, err)
	}
	if len(b.text) == 0 {
		b.text[0] = newLine()
		b.mod = insertMode
	}
	b.fileSync = time.Now().UTC()
	return nil
}

func (be *backend) CommandMode() bool {
	return be.commandMode
}

func (be *backend) MsgLine() line {
	return be.msgLine
}
