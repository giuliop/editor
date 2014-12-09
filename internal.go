package main

import (
	"fmt"
	"time"
)

// interalEditor specifies the api of the internal representation of buffers
type internalEditor interface {
	init()
	newBuffer(name string) *buffer
	//openFile()
	//closeBuffer()
	//saveBuffer()

	// operations on current buffer
	insertCh(ch rune)
	insertNewLineCh()
	deleteChBackward()
	//replaceCh(ch rune)
	//insertString(s []rune)
	//insertLineAbove()
	insertLineBelow()
	deleteLine()
	text() []line
	statusLine() []interface{}

	moveCursorUp(steps int)
	moveCursorDown(steps int)
	moveCursorLeft(steps int)
	moveCursorRight(steps int)
	setCursor(line int, pos int)
	cursorLine() int
}

// internalModel is the internal representation of internal editor
type internal struct {
	bufs []buffer // the open buffers
	cb   *buffer  // the current buffer
}

// buffer is the internal representation of a buffer
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

// initInternal returns the internal editor after having initialized it (for now with one empty buffer)
func initInternal() internalEditor {
	in := &internal{}
	in.init()
	return in
}

// init initializes the internal editor (for now with one empty buffer)
func (in *internal) init() {
	in.cb = in.newBuffer("")
}

// newBuffer adds a new empty buffer to internal and returns a pointer to it
func (i *internal) newBuffer(name string) *buffer {
	b := &buffer{
		text: make([]line, 1, 20),
		name: name,
	}
	b.cs = newMark(b)
	b.text[0] = newLine()
	i.bufs = append(i.bufs, *b)
	return b
}

func (in *internal) text() []line {
	return in.cb.text
}

func (in *internal) cursorLine() int {
	return in.cb.cs.line
}

func (in *internal) statusLine() []interface{} {
	cs := in.cb.cs
	return []interface{}{cs.pos + 1, fmt.Sprintf("%q", in.cb.text[cs.line]),
		cs.lastChPos() + 1, cs.maxLine() + 1}
}

func (in *internal) moveCursorUp(steps int) {
	in.cb.cs.moveUp(steps)
}

func (in *internal) moveCursorDown(steps int) {
	in.cb.cs.moveDown(steps)
}

func (in *internal) moveCursorRight(steps int) {
	in.cb.cs.moveRight(steps)
}

func (in *internal) moveCursorLeft(steps int) {
	in.cb.cs.moveLeft(steps)
}

func (in *internal) setCursor(line, pos int) {
	in.cb.cs.set(line, pos)
}

func (in *internal) insertCh(ch rune) {
	b := in.cb
	cs := b.cs
	b.text[cs.line] = append(b.text[cs.line], 0)
	copy(b.text[cs.line][cs.pos+1:], b.text[cs.line][cs.pos:])
	b.text[cs.line][cs.pos] = ch
	b.cs.pos++
}

func (in *internal) insertNewLineCh() {
	cs := in.cb.cs
	in.insertLineBelow()
	copy(in.cb.text[cs.line+1], in.cb.text[cs.line][cs.pos:])
	in.cb.text[cs.line] = append(in.cb.text[cs.line][:cs.pos], '\n')
	in.cb.cs.pos = 0
	in.cb.cs.line += 1
}

func newLine() line {
	return make([]rune, 0, 100)
}

func (in *internal) insertLineBelow() {
	b := in.cb
	cs := b.cs
	if cs.atLastLine() {
		b.text = append(b.text, newLine())
	} else {
		b.text = append(b.text, nil)
		copy(b.text[cs.line+1:], b.text[cs.line:])
		b.text[cs.line] = newLine()
	}
	b.cs = cs
}

func (in *internal) deleteChBackward() {
	b := in.cb
	cs := b.cs
	// if empty line delete it (unless first line in buffer)
	if cs.atLineStart() {
		if cs.atFirstLine() {
			return
		}
		in.deleteLine()
		cs.line -= 1
		// if last line delete newline char
		if cs.atLastLine() {
			b.text[cs.line] = b.text[cs.line][:cs.lastChPos()]
		}
		//cs.pos = cs.lastChPos()
		cs.pos = cs.lastChPos() + 1
	} else {
		cs.pos -= 1
		b.text[cs.line] = append(b.text[cs.line][:cs.pos], b.text[cs.line][cs.pos+1:]...)
	}
	b.cs = cs
}

func (in *internal) deleteLine() {
	b := in.cb
	b.text = append(b.text[:b.cs.line], b.text[b.cs.line+1:]...)
}

func (in *internal) DeleteChForward() {
}
