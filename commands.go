package main

import "unicode"

type direction int

const (
	right direction = iota
	left
	up
	down
)

type cmdContext struct {
	num        int        // times to execute the command
	char       rune       // optional char object
	cmd        cmdFunc    // the commnad to execute
	reg        regionFunc // optional region object
	point      *mark      // the cursor position
	custom     string     // optional string object
	customList []string   // optional string slice object
}

type command struct {
	cmd    cmdFunc
	parser parseFunc
}

type cmdFunc func(ctx *cmdContext)
type parseFunc func(ev *UIEvent, ctx *cmdContext, cmds chan *cmdContext) (parseFunc, bool)

var cmdKeys [2]map[Key]command

func initCmdTables() {
	cmdKeys[insertMode] = cmdKeysInsertMode
	cmdKeys[normalMode] = cmdKeysNormalMode
}

var cmdKeysNormalMode = map[Key]command{
	KeyEsc: command{exitProgram, nil},
}

var cmdCharsNormalMode = map[rune]command{
	'i': command{insertAtCs, nil},
	'a': command{appendAtCs, nil},
	'h': command{moveCursorLeft, nil},
	'j': command{moveCursorDown, nil},
	'k': command{moveCursorUp, nil},
	'l': command{moveCursorRight, nil},
	'd': command{delete_, parseRegion},
}

var cmdKeysInsertMode = map[Key]command{
	KeyEsc:        command{exitProgram, nil},
	KeyBackspace:  command{deleteCharBackward, nil},
	KeyBackspace2: command{deleteCharBackward, nil},
	KeyTab:        command{insertTab, nil},
	KeySpace:      command{insertSpace, nil},
	KeyEnter:      command{insertNewLine, nil},
	KeyCtrlJ:      command{insertNewLine, nil},
	KeyCtrlC:      command{toNormalMode, nil},
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

func delete_(ctx *cmdContext) {
	if ctx.num == 0 {
		ctx.num = 1
	}
	for i := 0; i < ctx.num; i++ {
		r := ctx.reg(*ctx.point)
		*ctx.point = eng.deleteRegion(r)
	}
	ctx.point.buf.cs = *ctx.point
}

type region struct {
	start mark
	end   mark
}
type regionFunc func(m mark) region

var regionFuncs = map[string]regionFunc{
	"e": toWordEnd,
	//"E":  toWORDEnd,
	//"w":  toNextWordStart,
	//"W":  toNextWORDStart,
	//"b":  toWordStart,
	//"B":  toWORDStart,
	//"0":  toFirstCharInLine,
	//"gh": toFirstCharInLine,
	//"$":  toLastCharInLine,
	//"gl": toLastCharInLine,
	//"iw": innerword,
	//"aw": aword,
}

func toWordEnd(m mark) region {
	m2 := m
	for {
		m2.moveRight(1)
		if m2.atLastLine() && m2.pos == m2.lastCharPos() {
			return region{m, m2}
		}
		c := m2.char()
		if !(unicode.IsLetter(c) || unicode.IsNumber(c)) {
			m2.moveLeft(1)
			return region{m, m2}
		}
	}
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
