package main

import (
	"fmt"
	"time"
)

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

func (be *backend) text(b *buffer) []line {
	return b.text
}

func (be *backend) cursorLine(b *buffer) int {
	return b.cs.line
}

func (be *backend) cursorPos(b *buffer) int {
	return b.cs.pos
}

func (be *backend) statusLine(b *buffer) []interface{} {
	cs := b.cs
	return []interface{}{cs.pos + 1, fmt.Sprintf("%q", b.text[cs.line]),
		cs.lastCharPos() + 1, cs.lastLine() + 1}
}

func (be *backend) insertChar(m mark, ch rune) {
	b := m.buf
	b.text[m.line] = append(b.text[m.line], 0)
	copy(b.text[m.line][m.pos+1:], b.text[m.line][m.pos:])
	b.text[m.line][m.pos] = ch
}

func (be *backend) insertNewLineChar(m mark) {
	b := m.buf
	be.insertLineBelow(m)
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

func (be *backend) insertLineBelow(m mark) {
	b := m.buf
	b.text = append(b.text, nil)
	m2 := mark{m.line + 1, 0, m.buf}
	copy(b.text[m2.line+1:], b.text[m2.line:])
	b.text[m2.line] = newLine()
}

// deleteCharBackward deletes the character before the mark and returns
// the new postion of the mark to be used to move the cursor if needed
func (be *backend) deleteCharBackward(m mark) mark {
	b := m.buf
	if m.atLineStart() {
		if m.atFirstLine() {
			return m
		}
		m.line -= 1
		m.pos = m.lastCharPos() + 1
		be.joinLineBelow(m)
	} else {
		m.pos -= 1
		b.text[m.line] = append(b.text[m.line][:m.pos], b.text[m.line][m.pos+1:]...)
	}
	return m
}

// deleteCharForward deletes the character under the mark
func (be *backend) deleteCharForward(m mark) {
	b := m.buf
	if m.atLineEnd() {
		if m.atLastLine() {
			return
		}
		be.joinLineBelow(m)
	} else {
		b.text[m.line] = append(b.text[m.line][:m.pos], b.text[m.line][m.pos+1:]...)
	}
}

func (be *backend) joinLineBelow(m mark) {
	if m.atLastLine() {
		return
	}
	m.buf.text[m.line] = append(m.buf.text[m.line][:m.lastCharPos()+1],
		m.buf.text[m.line+1]...)
	be.deleteLine(mark{m.line + 1, 0, m.buf})
}

func (be *backend) deleteLine(m mark) {
	b := m.buf
	b.text = append(b.text[:m.line], b.text[m.line+1:]...)
}

func (be *backend) deleteLines(m1, m2 mark) int {
	b := m1.buf
	b.text = append(b.text[:m1.line], b.text[m2.line+1:]...)
	return m1.line - m2.line
}

func (be *backend) deleteRegion(r region) mark {
	var fr, to = orderMarks(r.start, r.end)
	b := fr.buf
	if fr.line == to.line {
		b.text[fr.line] = append(b.text[fr.line][:fr.pos], b.text[fr.line][to.pos+1:]...)
	} else {
		if to.line > fr.line+1 {
			to.line -= e.deleteLines(mark{fr.line + 1, 0, b}, mark{to.line - 1, 0, b})
		}
		//delete required chars from fr and to lines
		b.text[fr.line] = b.text[fr.line][:fr.pos]
		b.text[to.line] = b.text[to.line][to.pos+1:]
	}
	fr.fixPos()
	return fr
}

func (be *backend) lastTextCharPos(m mark) mark {
	m2 := mark{m.lastLine(), 0, m.buf}
	m2.pos = m2.lastCharPos()
	m2.fixPos()
	return m2
}

func (be *backend) firstTextCharPos(m mark) mark {
	m2 := mark{0, 0, m.buf}
	return m2
}