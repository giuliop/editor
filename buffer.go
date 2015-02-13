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

// newLine returns a new line ending with a newline char
func newLine() line {
	ln := make([]rune, 1, 100)
	ln[0] = '\n'
	return ln
}

//mode represents an editing mode for the buffer
type mode int

// the different modes for a buffer
const (
	insertMode mode = iota
	normalMode
	commandMode
	visualMode
)

// content returns the slice containing the buffer text lines
func (b *buffer) content() []line {
	return b.text
}

// cursorline returns the line number of the buffer cursor
func (b *buffer) cursorLine() int {
	return b.cs.line
}

// cursorPos returns the pos number of the buffer cursor
func (b *buffer) cursorPos() int {
	return b.cs.pos
}

// statusLine returns the buffer statusline
func (b *buffer) statusLine() []interface{} {
	cs := b.cs
	return []interface{}{cs.pos + 1, fmt.Sprintf("%q", b.text[cs.line]),
		cs.lastCharPos() + 1, cs.lastLine() + 1}
}

// insertChar inserts the passed in rune after the mark
func (m mark) insertChar(ch rune) {
	b := m.buf
	b.text[m.line] = append(b.text[m.line], 0)
	copy(b.text[m.line][m.pos+1:], b.text[m.line][m.pos:])
	b.text[m.line][m.pos] = ch
}

// insertNewLineChar inserts a new line after the mark
func (m mark) insertNewLineChar() {
	b := m.buf
	m.insertLineBelow()
	b.text[m.line+1] = append(line(nil), b.text[m.line][m.pos:]...)
	b.text[m.line] = append(b.text[m.line][:m.pos], '\n')
}

// insertLineBelow inserts a line belor the mark
func (m mark) insertLineBelow() {
	b := m.buf
	b.text = append(b.text, nil)
	m2 := mark{m.line + 1, 0, m.buf}
	copy(b.text[m2.line+1:], b.text[m2.line:])
	b.text[m2.line] = newLine()
}

// deleteCharBackward deletes the character before the mark and returns
// the new postion of the mark to be used to move the cursor if needed
func (m mark) deleteCharBackward() mark {
	b := m.buf
	if m.atLineStart() {
		if m.atFirstLine() {
			return m
		}
		m.line -= 1
		m.pos = m.lastCharPos() + 1
		m.joinLineBelow()
	} else {
		m.pos -= 1
		b.text[m.line] = append(b.text[m.line][:m.pos], b.text[m.line][m.pos+1:]...)
	}
	return m
}

// deleteCharForward deletes the character under the mark
func (m mark) deleteCharForward() {
	b := m.buf
	if m.atLineEnd() {
		if m.atLastLine() {
			return
		}
		m.joinLineBelow()
	} else {
		b.text[m.line] = append(b.text[m.line][:m.pos], b.text[m.line][m.pos+1:]...)
	}
}

// joinLineBelow joins the mark's line with the line below
func (m mark) joinLineBelow() {
	if m.atLastLine() {
		return
	}
	m.buf.text[m.line] = append(m.buf.text[m.line][:m.lastCharPos()+1],
		m.buf.text[m.line+1]...)
	mark{m.line + 1, 0, m.buf}.deleteLine()
}

// deleteLines deletes the mark's line
func (m mark) deleteLine() {
	b := m.buf
	if len(b.text) == 1 {
		b.text[0] = newLine()
		return
	}
	b.text = append(b.text[:m.line], b.text[m.line+1:]...)
}

// deleteLines deletes the lines between the two marks including marks' lines
func (b *buffer) deleteLines(m1, m2 mark) int {
	if m1.atFirstLine() && m2.atLastLine() {
		b.text[0] = newLine()
		b.text = b.text[:1]
		return m2.line
	}
	b.text = append(b.text[:m1.line], b.text[m2.line+1:]...)
	return m2.line - m1.line + 1
}

// deleteRegion deletes the text between the two region's marks and returns a mark
// to be used to position the cursor if needed
func (r region) delete(dir direction) mark {
	var fr, to = orderMarks(r.start, r.end)
	// if we delete towards right we also want to delete the end mark's char
	if dir == right && !to.atEmptyLine() {
		to.pos++
	}
	b := fr.buf
	if fr.line == to.line {
		b.text[fr.line] = append(b.text[fr.line][:fr.pos], b.text[fr.line][to.pos:]...)
	} else {
		if to.line > fr.line+1 {
			to.line -= b.deleteLines(mark{fr.line + 1, 0, b}, mark{to.line - 1, 0, b})
		}
		//delete required chars from fr and to lines; if then empty delete them
		// making sure at least one line is left in the buffer
		b.text[to.line] = b.text[to.line][to.pos:]
		if to.atEmptyLine() {
			to.deleteLine()
		}
		switch {
		case fr.pos > 0:
			b.text[fr.line] = append(b.text[fr.line][:fr.pos], '\n')
		case fr.totalLines() == 1:
			b.text[fr.line] = newLine()
		default:
			fr.deleteLine()
			if fr.line == fr.totalLines() {
				fr.line -= 1
			}
		}
	}
	fr.fixPos()
	return fr
}

func (m mark) lastTextCharPos() mark {
	m2 := mark{m.lastLine(), 0, m.buf}
	m2.pos = m2.lastCharPos()
	m2.fixPos()
	return m2
}

func (m mark) firstTextCharPos() mark {
	m2 := mark{0, 0, m.buf}
	return m2
}
