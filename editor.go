package main

import (
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

type cursor struct {
	line  int // line number; starting from 1
	chPos int // char offset in the line
	viPos int // visual offset in the line
}

type line []rune

var (
	cs   *cursor
	text []line
)

const defCol = termbox.ColorDefault

func insertChar(ch rune) {
	text[cs.line] = append(text[cs.line], 0)
	copy(text[cs.line][cs.chPos+1:], text[cs.line][cs.chPos:])
	text[cs.line][cs.chPos] = ch
	cs.chPos++
	cs.viPos += runewidth.RuneWidth(ch)
}

func deleteChBackward() {
	if cs.chPos == 0 {
		return
	}
	cs.chPos -= 1
	cs.viPos -= runewidth.RuneWidth(text[cs.line][cs.chPos])
	text[cs.line] = append(text[cs.line][:cs.chPos], text[cs.line][cs.chPos+1:]...)
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
	text[0] = make([]rune, 0, 100)
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
