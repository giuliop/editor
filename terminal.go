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
	msgLine struct {
		msg         line  // to hold messages to display to user
		commandMode bool  // wether we are in command mode
		lastPane    *pane // last active pane before entering commandMode
	}
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

func (t *terminal) Init(b *buffer) error {
	v := &view{buf: b, cs: &mark{0, 0, b}, startline: 0}
	t.window = pane{nosplit, v, nil, nil, nil}
	t.curPane = &t.window
	t.msgLine.msg = line{}
	return termbox.Init()
}

func (t *terminal) Close() {
	termbox.Close()
	termbox.SetInputMode(termbox.InputEsc)
}

func (t *terminal) ReadMessageLine() line {
	return t.msgLine.msg
}

func (t *terminal) SetMessageLine(l line) {
	t.msgLine.msg = l
}

func (t *terminal) Draw() {
	t.clear()

	w, h := termbox.Size()
	h = h - 1
	t.drawMessageLine(h)

	t.window.draw(0, h-1, 0, w-1, t.curPane)
	t.flush()
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

func drawLine(dir splitType, from, to, at int, color termbox.Attribute) {
	switch dir {
	case vertical:
		for i := from; i <= to; i++ {
			setCellWithColor(at, i, ' ', defCol, color)
		}
	case horizontal:
		for i := from; i <= to; i++ {
			setCellWithColor(i, at, ' ', defCol, color)
		}
	}
}

func (p *pane) draw(lineFrom, lineTo, colFrom, colTo int, curPane *pane) {
	switch p.split {
	case vertical:
		midCol := (colTo-colFrom)/2 + 1 // we'll draw a separation line at midCol
		p.first.draw(lineFrom, lineTo, colFrom, midCol-1, curPane)
		p.second.draw(lineFrom, lineTo, midCol+1, colTo, curPane)
		drawLine(vertical, lineFrom, lineTo, midCol, termbox.ColorBlack)
	case horizontal:
		midLine := (lineTo-lineFrom)/2 + 1
		p.first.draw(lineFrom, midLine, colFrom, colTo, curPane)
		p.second.draw(midLine+1, lineTo, colFrom, colTo, curPane)
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
				setCellWithColor(j+colFrom, i+lineFrom, ch, termbox.ColorBlack,
					termbox.ColorWhite)
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

		if p == curPane {
			lineBeforeCs := v.buf.content()[v.cursorLine()][:v.cursorPos()]
			setCursor(lineVisualWidth(lineBeforeCs)+len(lineNumString)+colFrom,
				v.cursorLine()-v.startline+lineFrom)
		}
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

func (t *terminal) drawMessageLine(ln int) {
	msg := t.msgLine.msg
	if t.msgLine.commandMode == true {
		setCursor(lineVisualWidth(msg), ln)
	} else {
		msg = append(line(t.curPane.view.buf.name+" - "+
			t.curPane.view.buf.filename+" - "), msg...)
	}
	for i, ch := range msg {
		setCell(i, ln, ch)
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
	debug.Println(t.curPane)
	debug.Println(t.curPane.nextPane(dir))
	if p := t.curPane.nextPane(dir); p != nil {
		t.curPane = p
		t.Draw()
	}
}

func (p *pane) nextPane(dir direction) *pane {
	pr := p.parent
	if pr == nil {
		return nil
	}
	if (dir == right && pr.split == vertical && p == pr.first) ||
		(dir == down && pr.split == horizontal && p == pr.first) {
		if pr.second.split == nosplit {
			return pr.second
		}
		return pr.second.firstView()
	}
	if (dir == left && pr.split == vertical && p == pr.second) ||
		(dir == up && pr.split == horizontal && p == pr.second) {
		if pr.first.split == nosplit {
			return pr.first
		}
		return pr.first.firstView()
	}
	return pr.nextPane(dir)
}

func (p *pane) firstView() *pane {
	if p.first.split == nosplit {
		return p.first
	}
	return p.first.firstView()
}

func (t *terminal) enterCommandMode() {
	t.msgLine.commandMode = true
	t.msgLine.lastPane = t.curPane
	t.curPane = nil
}

func (t *terminal) exitCommandMode() {
	t.msgLine.commandMode = false
	t.curPane = t.msgLine.lastPane
	t.msgLine.lastPane = nil
}
