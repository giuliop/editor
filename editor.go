package main

import (
	"fmt"
	"io/ioutil"

	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

type cursor struct {
	line  int // line number; starting from 1
	chPos int // char offset in the line
	viPos int // visual offset in the line
}

type line []rune
type direction int

const (
	up direction = iota
	down
	left
	rigth
)

const defCol = termbox.ColorDefault

var (
	cs   *cursor
	text []line
)

func insertChar(ch rune) {
	text[cs.line] = append(text[cs.line], 0)
	copy(text[cs.line][cs.chPos+1:], text[cs.line][cs.chPos:])
	text[cs.line][cs.chPos] = ch
	cs.chPos++
	cs.viPos += runewidth.RuneWidth(ch)
}

func insertNewLineChar() {
	insertChar('\n')
	addNewLine(cs.line + 1)
	cs.chPos = 0
	cs.viPos = 0
	cs.line += 1
}

func newLine() line {
	return make([]rune, 0, 100)
}

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func debug(args ...interface{}) {
	msg := []byte(fmt.Sprintln(args))
	err := ioutil.WriteFile("debug", msg, 0644)
	check(err)
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

func newViPos() int {
	v := 0
	for _, c := range text[cs.line] {
		v += runewidth.RuneWidth(c)
	}
	return v
}

func deleteChBackward() {
	if cs.chPos == 0 {
		if cs.line == 0 {
			return
		} else {
			cs.line -= 1
			deleteLine(cs.line + 1)
			if text[cs.line][len(text[cs.line])-1] != '\n' {
				panic(fmt.Sprintf("Last char of line %v is %v, was wxpecting \\n", cs.line, text[cs.line]))
			}
			text[cs.line] = text[cs.line][:len(text[cs.line])-1]
			cs.chPos = len(text[cs.line])
			cs.viPos = newViPos()
		}
	} else {
		cs.chPos -= 1
		cs.viPos -= runewidth.RuneWidth(text[cs.line][cs.chPos])
		text[cs.line] = append(text[cs.line][:cs.chPos], text[cs.line][cs.chPos+1:]...)
	}
}

func deleteLine(line int) {
	text = append(text[:line], text[line+1:]...)
}

func DeleteRuneForward() {
}

func draw() {
	termbox.Clear(defCol, defCol)
	for i, line := range text {
		pos := 0
		for _, ch := range line {
			termbox.SetCell(pos, i, ch, defCol, defCol)
			pos += runewidth.RuneWidth(ch)
		}
	}
	termbox.SetCursor(cs.viPos, cs.line)
	termbox.Flush()
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)

	text = make([]line, 1, 20)
	text[0] = newLine()
	cs = &cursor{0, 0, 0}
	draw()

mainloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				break mainloop
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				deleteChBackward()
			case termbox.KeyTab:
				insertChar('\t')
			case termbox.KeySpace:
				insertChar(' ')
			case termbox.KeyEnter, termbox.KeyCtrlJ:
				insertNewLineChar()
			default:
				if ev.Ch != 0 {
					insertChar(ev.Ch)
				}
			}
		case termbox.EventError:
			panic(ev.Err)
		}
		draw()
	}
}
