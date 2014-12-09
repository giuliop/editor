package main

import (
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

func (t *terminal) Clear() error {
	return termbox.Clear(defCol, defCol)
}

func (t *terminal) Flush() error {
	return termbox.Flush()
}

func (t *terminal) SetCursor(x, y int) {
	termbox.SetCursor(x, y)
}

func (t *terminal) SetCell(x, y int, ch rune) {
	termbox.SetCell(x, y, ch, defCol, defCol)
}

func (t *terminal) HideCursor() {
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
