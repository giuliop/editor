package main

import (
	"fmt"
	"time"
)

// textEngine holds the buffers open in the editor and offers text manipulation
// primitives
type textEngine struct {
	bufs []buffer // the open buffers
}

// buffer is the representation of an open buffer
type buffer struct {
	text     []line
	cs       mark
	marks    []mark
	mod      mode
	name     string
	filename string
	fileSync time.Time
	modified bool
}

// line represent a line in a buffer
type line []rune

//mode represents an editing mode for the buffer
type mode int

const (
	insertMode mode = iota
	normalMode
	commandMode
	visualMode
)

// initTextEngine returns the textEngine editor after having initialized it
func initTextEngine() *textEngine {
	e := &textEngine{}
	return e
}

// newBuffer adds a new empty buffer to textEngine and returns a pointer to it
// Note that the last line of a buffer ends with a newline which is removed before
// saving to file
func (e *textEngine) newBuffer(name string) *buffer {
	b := &buffer{
		text: make([]line, 1, 20),
		name: name,
	}
	b.cs = newMark(b)
	b.text[0] = newLine()
	e.bufs = append(e.bufs, *b)
	return b
}

func (e *textEngine) text(b *buffer) []line {
	return b.text
}

func (e *textEngine) cursorLine(b *buffer) int {
	return b.cs.line
}

func (e *textEngine) cursorPos(b *buffer) int {
	return b.cs.pos
}

func (e *textEngine) statusLine(b *buffer) []interface{} {
	cs := b.cs
	return []interface{}{cs.pos + 1, fmt.Sprintf("%q", b.text[cs.line]),
		cs.lastCharPos() + 1, cs.maxLine() + 1}
}

func (e *textEngine) insertChar(m mark, ch rune) {
	b := m.buf
	b.text[m.line] = append(b.text[m.line], 0)
	copy(b.text[m.line][m.pos+1:], b.text[m.line][m.pos:])
	b.text[m.line][m.pos] = ch
}

func (e *textEngine) insertNewLineChar(m mark) {
	b := m.buf
	e.insertLineBelow(m)
	b.text[m.line+1] = append(line(nil), b.text[m.line][m.pos:]...)
	b.text[m.line] = append(b.text[m.line][:m.pos], '\n')
}

// newLine returns a new line ending with last which should be either a newline char
// or the endOfText char
func newLine() line {
	ln := make([]rune, 1, 100)
	ln[0] = '\n'
	return ln
}

func (e *textEngine) insertLineBelow(m mark) {
	b := m.buf
	b.text = append(b.text, nil)
	m2 := mark{m.line + 1, 0, m.buf}
	copy(b.text[m2.line+1:], b.text[m2.line:])
	b.text[m2.line] = newLine()
}

// deleteCharBackward deletes the character before the mark and returns
// the new postion of the mark to be used to move the cursor if needed
func (e *textEngine) deleteCharBackward(m mark) mark {
	b := m.buf
	if m.atLineStart() {
		if m.atFirstLine() {
			return m
		}
		m.line -= 1
		m.pos = m.lastCharPos() + 1
		e.joinLineBelow(m)
	} else {
		m.pos -= 1
		b.text[m.line] = append(b.text[m.line][:m.pos], b.text[m.line][m.pos+1:]...)
	}
	return m
}

// deleteCharForward deletes the character under the mark
func (e *textEngine) deleteCharForward(m mark) {
	b := m.buf
	if m.atLineEnd() {
		if m.atLastLine() {
			return
		}
		e.joinLineBelow(m)
	} else {
		b.text[m.line] = append(b.text[m.line][:m.pos], b.text[m.line][m.pos+1:]...)
	}
}

func (e *textEngine) joinLineBelow(m mark) {
	if m.atLastLine() {
		return
	}
	m.buf.text[m.line] = append(m.buf.text[m.line][:m.lastCharPos()+1],
		m.buf.text[m.line+1]...)
	e.deleteLine(mark{m.line + 1, 0, m.buf})
}

func (e *textEngine) deleteLine(m mark) {
	b := m.buf
	b.text = append(b.text[:m.line], b.text[m.line+1:]...)
}

func (e *textEngine) deleteRegion(r region) mark {
	var fr, to = orderMarks(r.start, r.end)
	b := fr.buf
	if fr.line == to.line {
		b.text[fr.line] = append(b.text[fr.line][:fr.pos], b.text[fr.line][to.pos+1:]...)
	} else {
		// delete all lines between the two marks
		m := mark{fr.line + 1, fr.pos, fr.buf}
		for ; m.line < to.line; m.line++ {
			e.deleteLine(m)
		}
		//delete required chars from fr and to lines
		b.text[fr.line] = b.text[fr.line][:fr.pos]
		b.text[to.line] = b.text[to.line][to.pos+1:]
	}
	fr.fixPos()
	return fr
}

func (e *textEngine) lastTextCharPos(m mark) mark {
	m2 := mark{m.maxLine(), 0, m.buf}
	m2.pos = m2.lastCharPos()
	m2.fixPos()
	return m2
}

func (e *textEngine) firstTextCharPos(m mark) mark {
	m2 := mark{0, 0, m.buf}
	return m2
}
