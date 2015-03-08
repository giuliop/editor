package main

import (
	"fmt"
	"strconv"

	"github.com/mattn/go-runewidth"
	"github.com/nsf/termbox-go"
)

const (
	defCol        = termbox.ColorDefault
	lineNumString = "    " // we use four chars for line numbers
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

	w, h := termbox.Size()
	var textLines int
	switch h {
	case 1:
		textLines = 1
	case 2:
		textLines = 1
		t.statusLine(h-1, w)
	default:
		textLines = h - 2
		t.statusLine(h-2, w)
		t.messageLine(h - 1)
	}

	v := t.curView
	text := t.curView.buf.content()
	v.fixScroll(textLines)
	endline := v.startline + textLines - 1
	if endline > len(text)-1 {
		endline = len(text) - 1
	}

	for i, line := range text[v.startline : endline+1] {
		// viPos tracks the visual position of chars in the line since some chars
		// might take more than one space on screen
		lineNum := strconv.Itoa(v.relativeLineNumber(v.startline + i))
		lineNum = lineNumString[:len(lineNumString)-len(lineNum)-1] + lineNum + " "
		for j, ch := range lineNum {
			t.setCellWithColor(j, i, ch, termbox.ColorBlack, termbox.ColorWhite)
		}
		viPos := len(lineNumString)
		for _, ch := range line {
			t.setCell(viPos, i, ch)
			viPos += runeWidth(ch)
		}
	}

	lineBeforeCs := text[v.cursorLine()][:v.cursorPos()]
	t.setCursor(lineVisualWidth(lineBeforeCs)+len(lineNumString),
		v.cursorLine()-v.startline)

	t.flush()
}

func (t *terminal) statusLine(line, width int) {
	args := t.curView.statusLine()
	s := fmt.Sprintf("Line %v, char %v, raw line %v, total chars %v, total lines %v",
		t.curView.cursorLine()+1, args[0], args[1], args[2], args[3])
	for i, ch := range s {
		t.setCellWithColor(i, line, ch, termbox.ColorBlack, termbox.ColorWhite)
	}
	for i := len(s); i < width; i++ {
		t.setCellWithColor(i, line, ' ', termbox.ColorBlack, termbox.ColorWhite)
	}
}

func (t *terminal) messageLine(line int) {
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

func (t *terminal) setCellWithColor(x, y int, ch rune, fg, bg termbox.Attribute) {
	termbox.SetCell(x, y, ch, fg, bg)
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

// runeWidth returns the number of visual spaces the rune takes on screen
func runeWidth(r rune) int {
	switch r {
	case '\t':
		return tabStop
	default:
		return runewidth.RuneWidth(r)
	}
}

// lineVisualWidth returns the number of visual spaces taken by the line ln
func lineVisualWidth(ln line) (i int) {
	for _, r := range ln {
		i += runeWidth(r)
	}
	return i
}
