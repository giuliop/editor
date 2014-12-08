package main

import "fmt"

type cursor struct {
	line  int // line number; starting from 1
	chPos int // char offset in the line
}

type line []rune

func insertChar(ch rune) {
	text[cs.line] = append(text[cs.line], 0)
	copy(text[cs.line][cs.chPos+1:], text[cs.line][cs.chPos:])
	text[cs.line][cs.chPos] = ch
	cs.chPos++
}

func insertNewLineChar() {
	insertChar('\n')
	addNewLine(cs.line + 1)
	cs.chPos = 0
	cs.line += 1
}

func newLine() line {
	return make([]rune, 0, 100)
}

func addNewLine(line int) {
	oldCs := *cs
	if line > 0 && text[line-1][len(text[line-1])-1] != '\n' {
		cs.line = line - 1
		cs.chPos = len(text[line-1])
		insertChar('\n')
	}
	// if line to be added at the end
	if line == oldCs.line+1 {
		text = append(text, newLine())
	} else {
		text = append(text, nil)
		copy(text[line+1:], text[line:])
		text[line] = newLine()
	}
	*cs = oldCs
}

func deleteChBackward() {
	if cs.chPos == 0 {
		if cs.line == 0 {
			return
		} else {
			cs.line -= 1
			deleteLine(cs.line + 1)
			if text[cs.line][len(text[cs.line])-1] != '\n' {
				debug(true, fmt.Sprintf("Last char of line %v is %v, was expecting \\n", cs.line, text[cs.line]))
			}
			text[cs.line] = text[cs.line][:len(text[cs.line])-1]
			cs.chPos = len(text[cs.line])
		}
	} else {
		cs.chPos -= 1
		text[cs.line] = append(text[cs.line][:cs.chPos], text[cs.line][cs.chPos+1:]...)
	}
}

func deleteLine(line int) {
	text = append(text[:line], text[line+1:]...)
}

func DeleteRuneForward() {
}
