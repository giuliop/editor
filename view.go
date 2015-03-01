package main

import "fmt"

type view struct {
	buf       *buffer
	cs        *mark
	startline int
}

const cursorLinesToMargin = 5

// cursorline returns the line number of the buffer cursor
func (v *view) cursorLine() int {
	return v.cs.line
}

// cursorPos returns the pos number of the buffer cursor
func (v *view) cursorPos() int {
	return v.cs.pos
}

// statusLine returns the buffer statusline
func (v *view) statusLine() []interface{} {
	cs := v.cs
	return []interface{}{cs.pos + 1, fmt.Sprintf("%q", v.buf.text[cs.line]),
		cs.lastCharPos() + 1, cs.lastLine() + 1}
}

// fixScroll modifies the startline of view v to make sure the cursors line
// fits in the passed number of lines
func (v *view) fixScroll(lines int) {
	switch {
	case lines < cursorLinesToMargin+1:
		v.startline = v.cs.line - (lines / 2)
	case v.cs.line < v.startline+cursorLinesToMargin:
		v.startline = v.cs.line - cursorLinesToMargin
	case v.cs.line-v.startline > lines-cursorLinesToMargin-1:
		v.startline = v.cs.line - lines + cursorLinesToMargin
	}
	if v.startline < 0 {
		v.startline = 0
	}
}

func (v *view) relativeLineNumber(line int) int {
	rel := line - v.cs.line
	if rel < 0 {
		rel *= -1
	}
	return rel
}
