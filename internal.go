package main

import "time"

type internal struct {
	bufs []buffer // the open buffers
	cb   *buffer  // the current buffer
}

type buffer struct {
	text     []line
	cs       cursor
	name     string
	filename string
	fileSync time.Time
	modified bool
}

type cursor struct {
	line  int // line number; starting from 1
	chPos int // char offset in the line
}

type line []rune

func initInternal() internal {
	in := internal{}
	in.cb = in.newBuffer("")
	return in
}

// newBuffer adds a new empty buffer to internal and returns a pointer to it
func (i *internal) newBuffer(name string) *buffer {
	b := buffer{
		text: make([]line, 1, 20),
		cs:   cursor{0, 0},
		name: name,
	}
	b.text[0] = newLine()
	in.bufs = append(in.bufs, b)
	return &b
}

func (b *buffer) insertChar(ch rune) {
	line := b.cs.line
	pos := b.cs.chPos
	b.text[line] = append(b.text[line], 0)
	copy(b.text[line][pos+1:], b.text[line][pos:])
	b.text[line][pos] = ch
	b.cs.chPos++
}

func (b *buffer) insertNewLineChar() {
	b.insertChar('\n')
	b.addNewLine(b.cs.line + 1)
	b.cs.chPos = 0
	b.cs.line += 1
}

func newLine() line {
	return make([]rune, 0, 100)
}

func (b *buffer) addNewLine(line int) {
	oldCs := b.cs
	if line > 0 && b.text[line-1][len(b.text[line-1])-1] != '\n' {
		b.cs.line = line - 1
		b.cs.chPos = len(b.text[line-1])
		b.insertChar('\n')
	}
	// if line to be added at the end
	if line == oldCs.line+1 {
		b.text = append(b.text, newLine())
	} else {
		b.text = append(b.text, nil)
		copy(b.text[line+1:], b.text[line:])
		b.text[line] = newLine()
	}
	b.cs = oldCs
}

func (b *buffer) deleteChBackward() {
	line := b.cs.line
	pos := b.cs.chPos
	// if empty line delete it (unless first line in buffer)
	if pos == 0 {
		if line == 0 {
			return
		}
		b.deleteLine(line)
		// if last line delete newline char
		line -= 1
		if line == len(b.text)-1 {
			b.text[line] = b.text[line][:len(b.text[line])-1]
		}
		// reposition cursor
		b.cs.line -= 1
		b.cs.chPos = len(b.text[line])
	} else {
		pos -= 1
		b.text[line] = append(b.text[line][:pos], b.text[line][pos+1:]...)
		// reposition cursor
		b.cs.chPos = pos
	}
}

func (b *buffer) deleteLine(line int) {
	b.text = append(b.text[:line], b.text[line+1:]...)
}

func (b *buffer) DeleteChForward() {
}
