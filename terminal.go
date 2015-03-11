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
	curPane *pane
	window  pane
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
	parent *pane // nil for window (main pane)
}

func (t *terminal) split(s splitType) {
	t.curPane.split = s
	t.curPane.first = &pane{nosplit, t.curPane.view, nil, nil, t.curPane}
	t.curPane.second = &pane{nosplit, copyView(t.curPane.view), nil, nil, t.curPane}
	t.curPane.view = nil
	t.curPane = t.curPane.second
}

func (t *terminal) SplitHorizontal() {
	t.split(horizontal)
}

func (t *terminal) SplitVertical() {
	t.split(vertical)
}

func (t *terminal) Init(b *buffer) error {
	v := &view{buf: b, cs: &mark{0, 0, b}, startline: 0}
	t.window = pane{nosplit, v, nil, nil, nil}
	t.curPane = &t.window
	return termbox.Init()
}

func (t *terminal) Close() {
	termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)
}

func (t *terminal) UserMessage(s string) {
	t.message = s
}

func (t *terminal) Draw() {
	t.clear()

	w, h := termbox.Size()

	// if we have at least two line we display message line
	if h > 1 {
		h = h - 1
		t.messageLine(h)
	}

	t.window.draw(0, h-1, 0, w-1)
	t.flush()
}

func (p *pane) draw(lineFrom, lineTo, colFrom, colTo int) {
	switch p.split {
	case vertical:
		p.first.draw(lineFrom, lineTo, colFrom, (colTo-colFrom)/2-1)
		p.second.draw(lineFrom, lineTo, (colTo-colFrom)/2+1, colTo)
	case horizontal:
		p.first.draw(lineFrom, (lineTo-lineFrom)/2-1, colFrom, colTo)
		p.second.draw((lineTo-lineFrom)/2+1, lineTo, colFrom, colTo)
	default:
		v := p.view
		text := v.buf.content()
		h := lineTo - lineFrom + 1
		//w := colTo - colFrom
		v.fixScroll(h)
		endline := v.startline + h - 1
		if endline > len(text)-1 {
			endline = len(text) - 1
		}

		for i, line := range text[v.startline : endline+1] {
			// draw the (relative) line numbers
			lineNum := strconv.Itoa(v.relativeLineNumber(v.startline + i))
			lineNum = lineNumString[:len(lineNumString)-len(lineNum)-1] + lineNum + " "
			for j, ch := range lineNum {
				setCellWithColor(j+colFrom, i+lineFrom, ch, termbox.ColorBlack, termbox.ColorWhite)
			}
			// viPos tracks the visual position of chars in the line since some chars
			// might take more than one space on screen
			viPos := len(lineNumString)
			for _, ch := range line {
				setCell(viPos+colFrom, i+lineFrom, ch)
				viPos += runeWidth(ch)
			}
		}
		// if we have at least two lines we dispaly the status line
		if h > 1 {
			p.statusLine(lineTo, colFrom, colTo)
		}

		lineBeforeCs := v.buf.content()[v.cursorLine()][:v.cursorPos()]
		setCursor(lineVisualWidth(lineBeforeCs)+len(lineNumString)+colFrom,
			v.cursorLine()-v.startline+lineFrom)

	}
}

func (p *pane) statusLine(line, colFrom, colTo int) {
	args := p.view.statusLine()
	s := fmt.Sprintf("Line %v, char %v, raw line %v, total chars %v, total lines %v",
		p.view.cursorLine()+1, args[0], args[1], args[2], args[3])
	for i, ch := range s {
		setCellWithColor(i+colFrom, line, ch, termbox.ColorBlack, termbox.ColorWhite)
	}
	for i := len(s) + colFrom; i < (colTo - colFrom + 1); i++ {
		setCellWithColor(i, line, ' ', termbox.ColorBlack, termbox.ColorWhite)
	}
}

func (t *terminal) messageLine(line int) {
	s := t.curPane.view.buf.name + " - " + t.curPane.view.buf.filename + " - " + t.message
	for i, ch := range s {
		setCell(i, line, ch)
	}
}

func (t *terminal) clear() error {
	return termbox.Clear(defCol, defCol)
}

func (t *terminal) flush() error {
	return termbox.Flush()
}

func setCursor(x, y int) {
	termbox.SetCursor(x, y)
}

func setCell(x, y int, ch rune) {
	termbox.SetCell(x, y, ch, defCol, defCol)
}

func setCellWithColor(x, y int, ch rune, fg, bg termbox.Attribute) {
	termbox.SetCell(x, y, ch, fg, bg)
}

func (t *terminal) hideCursor() {
	termbox.HideCursor()
}

func (t *terminal) PollEvent() UIEvent {
	ev := termbox.PollEvent()
	return UIEvent{
		t.curPane.view,
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

func (t *terminal) ToPane(dir direction) {
	if p := t.curPane.nextPane(dir); p != nil {
		t.curPane = p
	}
}

func (p *pane) nextPane(dir direction) *pane {
	pr := p.parent
	if pr == nil {
		return nil
	}
	if (dir == right && p.split == horizontal) ||
		(dir == down && p.split == vertical) {
		return pr.second
	}
	if (dir == left && p.split == horizontal) ||
		(dir == up && p.split == vertical) {
		return pr.first
	}
	return pr.nextPane(dir)
}
