package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"path/filepath"
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
// first one as current buffer; if the list in empty it returns a new buffer
// not existing filenames will also open new buffers
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
		text: make([]line, 1, 20),
		name: name,
	}
	b.cs = newMark(b)
	b.text[0] = newLine()
	be.bufs = append(be.bufs, b)
	return b
}

func (b *buffer) save() error {
	if b.filename == "" {
		b.filename = "newfile"
	}
	return b.saveAs(b.filename)
}

func (b *buffer) saveAs(filename string) error {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	br := &bufReader{b, 0, 0}
	_, err = io.Copy(f, br)
	return err
}

type bufReader struct {
	buf  *buffer
	line int
	pos  int
}

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
	return b, nil
}
