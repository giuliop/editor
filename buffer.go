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
	// TODO make it file based
	changeList changeList
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
	b.text[fr.line] = append(b.text[fr.line][:fr.pos], b.text[to.line][to.pos:]...)
	if to.line > fr.line {
		to.line -= b.deleteLines(mark{fr.line + 1, 0, b}, to)
		if fr.atEmptyLine() && fr.maxLine() > 0 {
			fr.deleteLine()
		}
	}
	fr.fixPos()
	return fr
}

func (m mark) insertText(text []line) {
	if len(text) == 0 {
		return
	}
	b := m.buf
	emptyBuf := len(b.text) == 1 && b.text[0][0] == '\n'
	if m.line > m.maxLine() {
		b.text = append(b.text, line{})
	}
	suffix := line{}
	if !emptyBuf {
		suffix = append(suffix, b.text[m.line][m.pos:]...)
	}
	b.text[m.line] = append(b.text[m.line][:m.pos], text[0]...)
	seg1 := b.text[:m.line+1]
	seg2 := text[1:]
	seg3 := b.text[m.line+1:]
	//debug.Printf("seg3 %q, len %v", seg3, len(seg3))
	lastline := text[len(text)-1]
	if len(suffix) > 0 {
		if lastline[len(lastline)-1] != '\n' {
			seg2[len(seg2)-1] = append(seg2[len(seg2)-1], suffix...)
		} else {
			seg3 = append([]line{suffix}, seg3...)
		}
	}
	//debug.Printf("seg1 %q, len %v", seg1, len(seg1))
	//debug.Printf("seg2 %q, len %v", seg2, len(seg2))
	//debug.Printf("seg3 %q, len %v", seg3, len(seg3))
	b.text = append(seg1, append(seg2, seg3...)...)
}

// copy copies and return the text between the two marks included
func (from mark) copy(to mark) (text []line) {
	start, end := orderMarks(from, to)
	if start.line == end.line {
		text = append(text, start.buf.text[start.line][start.pos:end.pos+1])
		return text
	}
	text = append(text, start.buf.text[start.line][start.pos:])
	for i := start.line + 1; i < end.line; i++ {
		text = append(text, start.buf.text[i])
	}
	text = append(text, start.buf.text[end.line][:end.pos+1])
	return text
}
