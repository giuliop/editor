package main

import (
	"fmt"
	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

const defCol = termbox.ColorDefault

type terminal struct {
	name string
}

func (t *terminal) Init() error {
	return termbox.Init()
}

func (t *terminal) Close() {
	termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)
}

func (t *terminal) Draw() {
	debug := true
	t.clear()
	// viPos tracks the visual position of chars in the line since some chars
	// might take two spaces on screen
	var viPos int
	for i, line := range in.text() {
		viPos = 0
		for _, ch := range line {
			t.setCell(viPos, i, ch)
			viPos += runewidth.RuneWidth(ch)
		}
	}
	if debug {
		offset := 80
		for l, line := range in.text() {
			s := fmt.Sprintf("%q", line)
			for c, ch := range s {
				t.setCell(c+offset, l, ch)
			}
		}
	}
	t.statusLine()
	t.setCursor(viPos, in.cursorLine())
	t.flush()
}

func (t *terminal) statusLine() {
	termw, termh := termbox.Size()
	termw += 0
	line := termh - 2
	args := in.statusLine()
	s := fmt.Sprintf("Line %v, char %v, raw line %v, total chars %v, total lines %v",
		in.cursorLine()+1, args[0], args[1], args[2], args[3])
	for i, ch := range s {
		t.setCell(i, line, ch)
	}
}

func (t *terminal) clear() error {
	return termbox.Clear(defCol, defCol)
}

func (t *terminal) flush() error {
	return termbox.Flush()
}

func (t *terminal) setCursor(x, y int) {
	termbox.SetCursor(x, y)
}

func (t *terminal) setCell(x, y int, ch rune) {
	termbox.SetCell(x, y, ch, defCol, defCol)
}

func (t *terminal) hideCursor() {
	termbox.HideCursor()
}

func (t *terminal) PollEvent() UIEvent {
	ev := termbox.PollEvent()
	return UIEvent{
		UIEventType(ev.Type),
		UIModifier(ev.Mod),
		Key(ev.Key),
		ev.Ch,
		ev.Width,
		ev.Height,
		ev.Err,
		ev.MouseX,
		ev.MouseY,
	}
}
