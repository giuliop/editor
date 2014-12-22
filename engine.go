package main

import (
	"fmt"
	"time"
)

// interalEditor specifies the api of the engine representation of buffers
type textEngine interface {
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

// engineModel is the engine representation of engine editor
type engine struct {
	bufs []buffer // the open buffers
	cb   *buffer  // the current buffer
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
	eng.init()
	return eng
}

// init initializes the engine editor (for now with one empty buffer)
func (eng *engine) init() {
	eng.cb = eng.newBuffer("")
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

func (eng *engine) text() []line {
	return eng.cb.text
}

func (eng *engine) cursorLine() int {
	return eng.cb.cs.line
}

func (eng *engine) statusLine() []interface{} {
	cs := eng.cb.cs
	return []interface{}{cs.pos + 1, fmt.Sprintf("%q", eng.cb.text[cs.line]),
		cs.lastChPos() + 1, cs.maxLine() + 1}
}

func (eng *engine) moveCursorUp(steps int) {
	eng.cb.cs.moveUp(steps)
}

func (eng *engine) moveCursorDown(steps int) {
	eng.cb.cs.moveDown(steps)
}

func (eng *engine) moveCursorRight(steps int) {
	eng.cb.cs.moveRight(steps)
}

func (eng *engine) moveCursorLeft(steps int) {
	eng.cb.cs.moveLeft(steps)
}

func (eng *engine) setCursor(line, pos int) {
	eng.cb.cs.set(line, pos)
}

func (eng *engine) insertCh(ch rune) {
	b := eng.cb
	cs := b.cs
	b.text[cs.line] = append(b.text[cs.line], 0)
	copy(b.text[cs.line][cs.pos+1:], b.text[cs.line][cs.pos:])
	b.text[cs.line][cs.pos] = ch
	b.cs.pos++
}

func (eng *engine) insertNewLineCh() {
	cs := eng.cb.cs
	eng.insertLineBelow()
	copy(eng.cb.text[cs.line+1], eng.cb.text[cs.line][cs.pos:])
	eng.cb.text[cs.line] = append(eng.cb.text[cs.line][:cs.pos], '\n')
	eng.cb.cs.pos = 0
	eng.cb.cs.line += 1
}

func newLine() line {
	return make([]rune, 0, 100)
}

func (eng *engine) insertLineBelow() {
	b := eng.cb
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

func (eng *engine) deleteChBackward() {
	b := eng.cb
	cs := b.cs
	// if empty line delete it (unless first line in buffer)
	if cs.atLineStart() {
		if cs.atFirstLine() {
			return
		}
		eng.deleteLine()
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

func (eng *engine) deleteLine() {
	b := eng.cb
	b.text = append(b.text[:b.cs.line], b.text[b.cs.line+1:]...)
}

func (eng *engine) DeleteChForward() {
}
