package main

import (
	"strconv"
	"unicode"
)

type direction int

const (
	right direction = iota
	left
	up
	down
)

type cmdContext struct {
	num   int
	key   rune
	cmd   cmdFunc
	point *mark
	save  bool
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
	'i': insertAtCs,
	'a': appendAtCs,
	'h': moveCursorLeft,
	'j': moveCursorDown,
	'k': moveCursorUp,
	'l': moveCursorRight,
	'0': number,
	'1': number,
	'2': number,
	'3': number,
	'4': number,
	'5': number,
	'6': number,
	'7': number,
	'8': number,
	'9': number,
}

var cmdKeysInsertMode = map[Key]cmdFunc{
	KeyEsc:        exitProgram,
	KeyBackspace:  deleteCharBackward,
	KeyBackspace2: deleteCharBackward,
	KeyTab:        insertTab,
	KeySpace:      insertSpace,
	KeyEnter:      insertNewLine,
	KeyCtrlJ:      insertNewLine,
	KeyCtrlC:      toNormalMode,
}

func toNormalMode(ctx *cmdContext) {
	ctx.point.buf.mod = normalMode
	if !ctx.point.atLineStart() {
		ctx.point.moveLeft(1)
	}
}

func insertAtCs(ctx *cmdContext) {
	ctx.point.buf.mod = insertMode
}

func appendAtCs(ctx *cmdContext) {
	// move cursor right unless empty line
	ctx.point.buf.mod = insertMode
	if !ctx.point.emptyLine() {
		ctx.point.moveRight(1)
	}
}

func number(ctx *cmdContext) {
	if unicode.IsDigit(ctx.key) {
		num, err := strconv.Atoi(strconv.Itoa(ctx.num) + string(ctx.key))
		if err == nil {
			ctx.num = num
		}
		ctx.save = true
	}
}

func moveCursorLeft(ctx *cmdContext) {
	if ctx.num == 0 {
		ctx.num = 1
	}
	ctx.point.moveLeft(ctx.num)
}

func moveCursorRight(ctx *cmdContext) {
	if ctx.num == 0 {
		ctx.num = 1
	}
	ctx.point.moveRight(ctx.num)
}

func moveCursorUp(ctx *cmdContext) {
	if ctx.num == 0 {
		ctx.num = 1
	}
	ctx.point.moveUp(ctx.num)
}

func moveCursorDown(ctx *cmdContext) {
	if ctx.num == 0 {
		ctx.num = 1
	}
	ctx.point.moveDown(ctx.num)
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
	eng.insertChar(*ctx.point, ctx.key)
	ctx.point.moveRight(1)
}
