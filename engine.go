package main

import (
	"fmt"
	"time"
)

// interalEditor specifies the api of the engine representation of buffers
type textEngine interface {
	newBuffer(name string) *buffer
	//openFile()
	//closeBuffer()
	//saveBuffer()

	insertChar(m mark, ch rune)
	insertNewLineChar(m mark)
	deleteCharBackward(m mark) mark
	//replaceChar(ch rune)
	//insertString(s []rune)
	//insertLineAbove()
	insertLineBelow(m mark)
	deleteLine(m mark)
	text(b *buffer) []line
	statusLine(b *buffer) []interface{}

	cursorLine(b *buffer) int
	cursorPos(b *buffer) int
}

// engineModel is the engine representation of engine editor
type engine struct {
	bufs []buffer // the open buffers
}

// buffer is the engine representation of a buffer
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

//mode represents an editing mode for the editor
type mode int

const (
	insertMode mode = iota
	normalMode
	commandMode
)

// initengine returns the engine editor after having initialized it (for now with one empty buffer)
func initEngine() textEngine {
	eng := &engine{}
	return eng
}

// newBuffer adds a new empty buffer to engine and returns a pointer to it
func (eng *engine) newBuffer(name string) *buffer {
	b := &buffer{
		text: make([]line, 1, 20),
		name: name,
	}
	b.cs = newMark(b)
	b.text[0] = newLine()
	eng.bufs = append(eng.bufs, *b)
	return b
}

func (eng *engine) text(b *buffer) []line {
	return b.text
}

func (eng *engine) cursorLine(b *buffer) int {
	return b.cs.line
}

func (eng *engine) cursorPos(b *buffer) int {
	return b.cs.pos
}

func (eng *engine) statusLine(b *buffer) []interface{} {
	cs := b.cs
	return []interface{}{cs.pos + 1, fmt.Sprintf("%q", b.text[cs.line]),
		cs.lastCharPos() + 1, cs.maxLine() + 1}
}

func (eng *engine) insertChar(m mark, ch rune) {
	b := m.buf
	b.text[m.line] = append(b.text[m.line], 0)
	copy(b.text[m.line][m.pos+1:], b.text[m.line][m.pos:])
	b.text[m.line][m.pos] = ch
}

func (eng *engine) insertNewLineChar(m mark) {
	b := m.buf
	eng.insertLineBelow(m)
	b.text[m.line+1] = append(b.text[m.line+1], b.text[m.line][m.pos:m.lastCharPos()+1]...)
	b.text[m.line] = append(b.text[m.line][:m.pos], '\n')
}

func newLine() line {
	return make([]rune, 0, 100)
}

func (eng *engine) insertLineBelow(m mark) {
	b := m.buf
	b.text = append(b.text, newLine())
	m2 := mark{m.line + 1, 0, m.buf}
	if !(m2.line == len(b.text)-1) {
		copy(b.text[m2.line+1:], b.text[m2.line:])
		b.text[m2.line] = newLine()
	}
}

// deleteCharBackward deleted the character before the mark and returns the new postion of the mark
// to be used to move the cursor if needed
func (eng *engine) deleteCharBackward(m mark) mark {
	b := m.buf
	// if empty line delete it (unless first line in buffer)
	if m.atLineStart() {
		if m.atFirstLine() {
			return m
		}
		m.line -= 1
		eng.joinLineBelow(m)
		// if last line delete newline char
		if m.atLastLine() {
			b.text[m.line] = b.text[m.line][:m.lastCharPos()+1]
		}
		m.pos = m.lastCharPos() + 1
	} else {
		m.pos -= 1
		b.text[m.line] = append(b.text[m.line][:m.pos], b.text[m.line][m.pos+1:]...)
	}
	return m
}

func (eng *engine) joinLineBelow(m mark) {
	if m.atLastLine() {
		return
	}
	m.buf.text[m.line] = append(m.buf.text[m.line][:m.lastCharPos()+1],
		m.buf.text[m.line+1]...)
	eng.deleteLine(mark{m.line + 1, 0, m.buf})
}

func (eng *engine) deleteLine(m mark) {
	b := m.buf
	b.text = append(b.text[:m.line], b.text[m.line+1:]...)
}

func (eng *engine) DeleteCharForward() {
}
