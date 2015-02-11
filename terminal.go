package main

import (
	"fmt"

	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

const defCol = termbox.ColorDefault

type terminal struct {
	name    string
	curBuf  *buffer
	message string // to hold messages to display to user
}

func (t *terminal) Init(b *buffer) error {
	t.curBuf = b
	return termbox.Init()
}

func (t *terminal) Close() {
	termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)
}

func (t *terminal) userMessage(s string) {
	t.message = s
}

func (t *terminal) Draw() {
	var _debug bool
	//_debug = true
	t.clear()
	// viPos tracks the visual position of chars in the line since some chars
	// might take two spaces on screen
	b := t.CurrentBuffer()
	for i, line := range b.content() {
		viPos := 0
		for _, ch := range line {
			t.setCell(viPos, i, ch)
			viPos += runewidth.RuneWidth(ch)
		}
	}
	if _debug {
		offset := 80
		for l, line := range b.content() {
			s := fmt.Sprintf("%q", line)
			for c, ch := range s {
				t.setCell(c+offset, l, ch)
			}
		}
	}
	t.statusLine()
	t.messageLine()
	stringBeforeCs := string(b.content()[b.cursorLine()][:b.cursorPos()])
	t.setCursor(runewidth.StringWidth(stringBeforeCs), b.cursorLine())
	t.flush()
}

func (t *terminal) statusLine() {
	termw, termh := termbox.Size()
	termw += 0
	line := termh - 4
	args := t.CurrentBuffer().statusLine()
	s := fmt.Sprintf("Line %v, char %v, raw line %v, total chars %v, total lines %v",
		t.CurrentBuffer().cursorLine()+1, args[0], args[1], args[2], args[3])
	for i, ch := range s {
		t.setCell(i, line, ch)
	}
}

func (t *terminal) messageLine() {
	termw, termh := termbox.Size()
	termw += 0
	line := termh - 2
	s := t.curBuf.name + " - " + t.curBuf.filename + " - " + t.message
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
		t.curBuf,
		UIEventType(ev.Type),
		UIModifier(ev.Mod),
		Keypress{Key(ev.Key), ev.Ch, ev.Ch == 0},
		ev.Width,
		ev.Height,
		ev.Err,
		ev.MouseX,
		ev.MouseY,
	}
}

func (t *terminal) CurrentBuffer() *buffer {
	return t.curBuf
}
