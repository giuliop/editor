package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"
	"unicode/utf8"
)

// backend holds the buffers open in the editor
type backend struct {
	bufs []*buffer // the open buffers
}

// initBackend returns the backend after having initialized it
func initBackend() backend {
	be := backend{}
	return be
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
		b, err := be.openFile(fn)
		if err != nil {
			b = be.newBuffer("")
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
	b := &buffer{
		text:       make([]line, 1, 20),
		name:       name,
		changeList: changeList{ops: make([]bufferChange, 1)},
	}
	newMark(b).initLastInsert()
	b.text[0] = newLine()
	be.bufs = append(be.bufs, b)
	return b
}

// save saves the buffer; if it has no filename associated to it, it will be
// saved as "newfile"
func (b *buffer) save() error {
	if b.filename == "" {
		b.filename = "newfile"
	}
	return b.saveAs(b.filename)
}

// saveAs saves the buffer in a file named after the passed parameter
func (b *buffer) saveAs(filename string) error {
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

// openFile open the file named after the passed paramenter and adds it to
// backend buffer list, returning a pointer to its buffer or an error
func (be *backend) openFile(filename string) (*buffer, error) {
	b := be.newBuffer(filename)
	path, err := filepath.Abs(filename)
	if err != nil {
		fatalError(err)
	}
	b.filename = path
	errPrefix := "Hmpf, I cannot open the file '%v', "
	f, err := os.Open(path)
	switch {
	case os.IsNotExist(err):
		return b, nil
	case os.IsPermission(err):
		return nil, fmt.Errorf(errPrefix+"we do not have access rights", path)
	case err != nil:
		return nil, fmt.Errorf(errPrefix+"got this error:\n%v\n", path, err)
	}
	defer f.Close()

	sc := bufio.NewScanner(f)
	b.text = b.text[:0]
	b.mod = normalMode
	for sc.Scan() {
		b.text = append(b.text, []rune(sc.Text()+"\n"))
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf(errPrefix+"got this error:\n%v\n", path, err)
	}
	if len(b.text) == 0 {
		b.text[0] = newLine()
		b.mod = insertMode
	}
	b.fileSync = time.Now().UTC()
	return b, nil
}
