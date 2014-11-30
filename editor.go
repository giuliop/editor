package main

import (
	"github.com/nsf/termbox-go"
)

type cursor struct {
	x int
	y int
}

var cs *cursor

const defCol = termbox.ColorDefault

func InsertRune(ch rune) {
	termbox.SetCell(cs.x, cs.y, ch, defCol, defCol)
	cs.x++
}

func DeleteRuneBackward() {
	termbox.SetCell(cs.x-1, cs.y, ' ', defCol, defCol)
	cs.x--
}

func DeleteRuneForward() {
	termbox.SetCell(cs.x, cs.y, ' ', defCol, defCol)
	cs.x--
}

func draw() {
	termbox.SetCursor(cs.x, cs.y)
	termbox.Flush()
}

func main() {
	err := termbox.Init()
	if err != nil {
		panic(err)
	}
	defer termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)

	cs = &cursor{0, 0}
	draw()

mainloop:
	for {
		switch ev := termbox.PollEvent(); ev.Type {
		case termbox.EventKey:
			switch ev.Key {
			case termbox.KeyEsc:
				break mainloop
			case termbox.KeyBackspace, termbox.KeyBackspace2:
				DeleteRuneBackward()
			case termbox.KeyTab:
				InsertRune('\t')
			case termbox.KeySpace:
				InsertRune(' ')
			default:
				if ev.Ch != 0 {
					InsertRune(ev.Ch)
				}
			}
		case termbox.EventError:
			panic(ev.Err)
		}
		draw()
	}
}
