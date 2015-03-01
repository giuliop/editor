package main

import (
	"fmt"

	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

const (
	defCol      = termbox.ColorDefault
	statusLines = 5
)

type terminal struct {
	curView *view
	panes   pane
	message string // to hold messages to display to user
}

type splitType int

const (
	nosplit splitType = iota
	horizontal
	vertical
)

type pane struct {
	split  splitType
	view   *view // if split != nosplit this is nil
	first  *pane // left or top split, or nil
	second *pane // right or bottom split, or nil
}

func (t *terminal) Init(b *buffer) error {
	t.curView = &view{buf: b, cs: &mark{0, 0, b}, startline: 0}
	t.panes = pane{nosplit, t.curView, nil, nil}
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
	t.clear()
	v := t.curView
	text := t.curView.buf.content()
	_, h := termbox.Size()
	linesVisible := h - statusLines
	v.fixScroll(linesVisible)
	endline := v.startline + linesVisible - 1
	if endline > len(text)-1 {
		endline = len(text) - 1
	}
	debug.Println(v.startline, endline, len(text))
	for i, line := range text[v.startline : endline+1] {
		// viPos tracks the visual position of chars in the line since some chars
		// might take more than one space on screen
		viPos := 0
		for _, ch := range line {
			t.setCell(viPos, i, ch)
			viPos += runewidth.RuneWidth(ch)
		}
	}
	t.statusLine()
	t.messageLine()
	//debug.Printf("line %v, maxline %v", b.cursorLine(), len(b.content())-1)
	//debug.Printf("pos %v, maxpos %v", b.cursorPos(), len(b.content()[b.cursorLine()])-1)
	stringBeforeCs := string(text[v.cursorLine()][:v.cursorPos()])
	t.setCursor(runewidth.StringWidth(stringBeforeCs), v.cursorLine()-v.startline)
	t.flush()
}

func (t *terminal) statusLine() {
	termw, termh := termbox.Size()
	termw += 0
	line := termh - 4
	args := t.curView.statusLine()
	s := fmt.Sprintf("Line %v, char %v, raw line %v, total chars %v, total lines %v",
		t.curView.cursorLine()+1, args[0], args[1], args[2], args[3])
	for i, ch := range s {
		t.setCell(i, line, ch)
	}
}

func (t *terminal) messageLine() {
	termw, termh := termbox.Size()
	termw += 0
	line := termh - 2
	s := t.curView.buf.name + " - " + t.curView.buf.filename + " - " + t.message
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
		t.curView,
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
