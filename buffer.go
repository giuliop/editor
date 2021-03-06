package main

import "time"

// buffer is the representation of an open buffer
type buffer struct {
	text        text
	marks       []mark
	savedCursor mark // to save the cursor when the buffer has no view attached
	mod         mode
	name        string
	filename    string
	filetype    filetype
	fileSync    time.Time
	modified    bool       // true if not synched with file
	changeList  changeList // for undo / redo (TODO make it file based)
	lastInsert  insertText // text added in last insertMode session
}

// insertText represents the change to the buffer's text since insertMode was
// last entered
type insertText struct {
	newText text  // the new text inserted
	oldText text  // the old text deleted
	start   *mark // where the change starts
}

type text []line

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

// insertChar inserts the passed in rune after the mark
func (m mark) insertChar(ch rune) {
	if ch == '\n' {
		panic("Wrong function to insert newline")
	}
	b := m.buf
	b.text[m.line] = append(b.text[m.line], 0)
	copy(b.text[m.line][m.pos+1:], b.text[m.line][m.pos:])
	b.text[m.line][m.pos] = ch

	// add undo info
	b.lastInsert.newText.appendChar(ch)
}

// insertNewLineChar inserts a new line after the mark
func (m mark) insertNewLineChar() {
	b := m.buf
	m.insertLineBelow()
	b.text[m.line+1] = append(line(nil), b.text[m.line][m.pos:]...)
	b.text[m.line] = append(b.text[m.line][:m.pos], '\n')

	// add undo info
	b.lastInsert.newText.appendChar('\n')
}

func (m mark) insertTab() {
	for _, r := range tab {
		m.insertChar(r)
	}
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
	var deleted rune // for undo info

	if m.atLineStart() {
		if m.atFirstLine() {
			return m
		}
		m.line -= 1
		m.pos = m.lastCharPos() + 1
		m.joinLineBelow()
		deleted = '\n'
	} else {
		m.pos -= 1
		deleted = m.char()
		b.text[m.line] = append(b.text[m.line][:m.pos], b.text[m.line][m.pos+1:]...)
	}

	// add undo info
	b.lastInsert.oldText.prependChar(deleted)
	if m.isBefore(*b.lastInsert.start) {
		b.lastInsert.start = &m
	}

	return m
}

// deleteCharForward deletes the character after the mark
func (m mark) deleteCharForward() {
	b := m.buf
	var deleted rune // for undo info

	if m.atLineEnd() {
		if m.atLastLine() {
			return
		}
		m.joinLineBelow()
		deleted = '\n'
	} else {
		deleted = m.char()
		b.text[m.line] = append(b.text[m.line][:m.pos], b.text[m.line][m.pos+1:]...)
	}

	// add undo info
	b.lastInsert.oldText.appendChar(deleted)
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

// delete deletes the text between the two region's marks and returns a mark
// to be used to position the cursor if needed
func (r region) delete() mark {
	var fr, to = orderMarks(r.start, r.end)
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

// replace replaces the region with newText
func (r region) replace(newText text) {
	r.delete()
	r.start.insertText(newText)
}

// insertText inserts the text at mark
func (m mark) insertText(text text) {
	if text.empty() {
		return
	}

	b := m.buf
	if m.line > m.maxLine() {
		b.text = append(b.text, line{})
	}
	suffix := append(line{}, b.text[m.line][m.pos:]...)
	b.text[m.line] = append(b.text[m.line][:m.pos], text[0]...)
	seg1 := b.text[:m.line+1]
	seg2 := text[1:]
	seg3 := b.text[m.line+1:]
	if len(suffix) > 0 {
		switch {
		case seg1.lastChar() != '\n':
			seg1.appendChars(suffix)
		case len(seg2) > 0 && seg2.lastChar() != '\n':
			seg2.appendChars(suffix)
		default:
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
		text = append(text, start.buf.text[start.line][start.pos:end.pos])
		return text
	}
	text = append(text, start.buf.text[start.line][start.pos:])
	for i := start.line + 1; i < end.line; i++ {
		text = append(text, start.buf.text[i])
	}
	text = append(text, start.buf.text[end.line][:end.pos])
	return text
}

// setMode sets the buffer mode. Before doing that it calls the functions
// exitingMode and enteringMode to hook up actions. It returns a closure
// with the exitedMode and enteredMode functions to be called by the caller
// after setting the cursor
func (m mark) setMode(newM mode) func(newCursor *mark) {
	oldM := m.buf.mod
	change := oldM != newM
	if !change {
		return nil
	}
	m.exitingMode()
	m.enteringMode(newM)

	m.buf.mod = newM

	return func(newCursor *mark) {
		m.exitedMode(oldM)
		newCursor.enteredMode(newM)
	}
}

func (m mark) enteredMode(mod mode) {
	switch mod {
	case insertMode:
		m.initLastInsert()
		//debug.Println("entered insertMode\n")
	case normalMode:
		//debug.Println("entered normalMode\n")
	}
}

func (m mark) exitedMode(mod mode) {
	switch mod {
	case insertMode:
		//debug.Println("exited insertMode")
	case normalMode:
		//debug.Println("exited normalMode")
	}
}

func (m mark) exitingMode() {
	switch m.buf.mod {
	case insertMode:
		//debug.Println("exiting insertMode")
		m.addUndoRedoLastInsert()
	case normalMode:
		//debug.Println("exiting normalMode")
	}
}

func (m mark) enteringMode(mod mode) {
	switch mod {
	case insertMode:
		//debug.Println("entering insertMode")
	case normalMode:
		//debug.Println("entering normalMode")
	}
}

func (m mark) addUndoRedoLastInsert() {
	if m.buf.lastInsert.newText.empty() && m.buf.lastInsert.oldText.empty() {
		return
	}

	start := *m.buf.lastInsert.start
	end := mark{m.line, m.pos, m.buf}
	undoCtx := undoContext{
		text:  m.buf.lastInsert.oldText,
		start: start,
		end:   end,
	}
	undoEnd := start.toEndofText(m.buf.lastInsert.oldText)
	regF := func(m mark) (region, direction) {
		return region{start: start,
			end: undoEnd}, right
	}
	redoCtx := &cmdContext{
		num:      1,
		cmd:      replace,
		point:    m.buf.lastInsert.start,
		text:     m.buf.lastInsert.newText,
		reg:      regF,
		cmdChans: cmdStack{commands, make(chan struct{}, 1)},
	}
	m.buf.changeList.add(newRedoCtx(redoCtx), undoCtx)
}

func textToString(t []line) string {
	s := ""
	for _, line := range t {
		s += string(line)
	}
	return s
}
