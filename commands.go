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
	point  *mark
}

type cmdFunc func(ctx *cmdContext)

var cmdKeys [2]map[Key]cmdFunc

func initCmdTables() {
	cmdKeys[insertMode] = cmdKeysInsertMode
	cmdKeys[normalMode] = cmdKeysNormalMode
}

var cmdKeysNormalMode = map[Key]cmdFunc{
	KeyEsc: exitProgram,
}

var cmdCharsNormalMode = map[rune]cmdFunc{
	'i': enterInsertMode,
	'a': enterInsertModeAsAppend,
	'h': moveCursorLeft,
	'j': moveCursorDown,
	'k': moveCursorUp,
	'l': moveCursorRight,
}

var cmdKeysInsertMode = map[Key]cmdFunc{
	KeyEsc:        exitProgram,
	KeyBackspace:  deleteCharBackward,
	KeyBackspace2: deleteCharBackward,
	KeyTab:        insertTab,
	KeySpace:      insertSpace,
	KeyEnter:      insertNewLine,
	KeyCtrlJ:      insertNewLine,
	KeyCtrlC:      enterNormalMode,
}

func enterNormalMode(ctx *cmdContext) {
	ctx.point.buf.mod = normalMode
	if !ctx.point.atLineStart() {
		ctx.point.moveLeft(1)
	}
}

func enterInsertMode(ctx *cmdContext) {
	ctx.point.buf.mod = insertMode
}

func enterInsertModeAsAppend(ctx *cmdContext) {
	ctx.point.buf.mod = insertMode
	// move cursor right unless empty line
	if !ctx.point.emptyLine() {
		ctx.point.moveRight(1)
	}
}

func moveCursorLeft(ctx *cmdContext) {
	ctx.point.moveLeft(1)
}

func moveCursorRight(ctx *cmdContext) {
	ctx.point.moveRight(1)
}

func moveCursorUp(ctx *cmdContext) {
	ctx.point.moveUp(1)
}

func moveCursorDown(ctx *cmdContext) {
	ctx.point.moveDown(1)
}

func exitProgram(ctx *cmdContext) {
	exitSignal <- true
}

func deleteCharBackward(ctx *cmdContext) {
	*ctx.point = eng.deleteCharBackward(*ctx.point)
}

func insertTab(ctx *cmdContext) {
	eng.insertChar(*ctx.point, '\t')
	ctx.point.moveRight(1)
}

func insertSpace(ctx *cmdContext) {
	eng.insertChar(*ctx.point, ' ')
	ctx.point.moveRight(1)
}

func insertNewLine(ctx *cmdContext) {
	eng.insertNewLineChar(*ctx.point)
	ctx.point.set(ctx.point.line+1, 0)
}

func insertChar(ctx *cmdContext) {
	eng.insertChar(*ctx.point, ctx.char)
	ctx.point.moveRight(1)
}
