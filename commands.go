package main

type direction int

const (
	right direction = iota
	left
	up
	down
)

type cmdContext struct {
	times  int
	action string
	object string
	char   rune
	dir    direction
	point  mark
}

type cmdFunc func()

var ctx = new(cmdContext)

var cmdNames = map[string]cmdFunc{}
var cmdKeys = map[Key]cmdFunc{
	KeyEsc:        exitProgram,
	KeyBackspace:  deleteCharBackward,
	KeyBackspace2: deleteCharBackward,
	KeyTab:        insertTab,
	KeySpace:      insertSpace,
	KeyEnter:      insertNewLine,
	KeyCtrlJ:      insertNewLine,
}

func exitProgram() {
	exitSignal <- true
}

func deleteCharBackward() {
	b := ui.CurrentBuffer()
	b.cs = eng.deleteCharBackward(b.cs)
}

func insertTab() {
	b := ui.CurrentBuffer()
	eng.insertChar(b.cs, '\t')
	b.cs.pos++
}

func insertSpace() {
	b := ui.CurrentBuffer()
	eng.insertChar(b.cs, ' ')
	b.cs.pos++
}

func insertNewLine() {
	b := ui.CurrentBuffer()
	eng.insertNewLineChar(b.cs)
	b.cs.pos = 0
	b.cs.line += 1
}

func insertChar(ch rune) {
	b := ui.CurrentBuffer()
	eng.insertChar(b.cs, ch)
	b.cs.pos++
}
