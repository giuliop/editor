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
	name     string
	filename string
	fileSync time.Time
	modified bool
}

// line represent a line in a buffer
type line []rune

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
	copy(b.text[m.line+1], b.text[m.line][m.pos:])
	b.text[m.line] = append(b.text[m.line][:m.pos], '\n')
}

func newLine() line {
	return make([]rune, 0, 100)
}

func (eng *engine) insertLineBelow(m mark) {
	b := m.buf
	if m.atLastLine() {
		b.text = append(b.text, newLine())
	} else {
		b.text = append(b.text, nil)
		copy(b.text[m.line+1:], b.text[m.line:])
		b.text[m.line] = newLine()
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
		eng.deleteLine(m)
		m.line -= 1
		// if last line delete newline char
		if m.atLastLine() {
			b.text[m.line] = b.text[m.line][:m.lastCharPos()]
		}
		m.pos = m.lastCharPos() + 1
	} else {
		m.pos -= 1
		b.text[m.line] = append(b.text[m.line][:m.pos], b.text[m.line][m.pos+1:]...)
	}
	return m
}

func (eng *engine) deleteLine(m mark) {
	b := m.buf
	b.text = append(b.text[:m.line], b.text[m.line+1:]...)
}

func (eng *engine) DeleteCharForward() {
}
