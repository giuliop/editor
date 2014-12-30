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
	key    rune
	dir    direction
	point  *mark
}

type cmdFunc func(ctx *cmdContext) (done bool)

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

func enterNormalMode(ctx *cmdContext) bool {
	ctx.point.buf.mod = normalMode
	if !ctx.point.atLineStart() {
		ctx.point.moveLeft(1)
	}
	return true
}

func enterInsertMode(ctx *cmdContext) bool {
	ctx.point.buf.mod = insertMode
	return true
}

func enterInsertModeAsAppend(ctx *cmdContext) bool {
	ctx.point.buf.mod = insertMode
	// move cursor right unless empty line
	if !ctx.point.emptyLine() {
		ctx.point.moveRight(1)
	}
	return true
}

func moveCursorLeft(ctx *cmdContext) bool {
	ctx.point.moveLeft(1)
	return true
}

func moveCursorRight(ctx *cmdContext) bool {
	ctx.point.moveRight(1)
	return true
}

func moveCursorUp(ctx *cmdContext) bool {
	ctx.point.moveUp(1)
	return true
}

func moveCursorDown(ctx *cmdContext) bool {
	ctx.point.moveDown(1)
	return true
}

func exitProgram(ctx *cmdContext) bool {
	exitSignal <- true
	return true
}

func deleteCharBackward(ctx *cmdContext) bool {
	*ctx.point = eng.deleteCharBackward(*ctx.point)
	return true
}

func insertTab(ctx *cmdContext) bool {
	eng.insertChar(*ctx.point, '\t')
	ctx.point.moveRight(1)
	return true
}

func insertSpace(ctx *cmdContext) bool {
	eng.insertChar(*ctx.point, ' ')
	ctx.point.moveRight(1)
	return true
}

func insertNewLine(ctx *cmdContext) bool {
	eng.insertNewLineChar(*ctx.point)
	ctx.point.set(ctx.point.line+1, 0)
	return true
}

func insertChar(ctx *cmdContext) bool {
	eng.insertChar(*ctx.point, ctx.key)
	ctx.point.moveRight(1)
	return true
}
