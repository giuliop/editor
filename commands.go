package main

//type direction int
//const (
//right direction = iota
//left
//up
//down
//)

// cmdContext is used to store all the info we need to process commands
type cmdContext struct {
	num        int        // times to execute the command
	cmd        cmdFunc    // the commnad to execute
	point      *mark      // the cursor position
	char       rune       // the last input char
	cmdString  string     // the input string defining the command
	argString  string     // optional input string defining the command arg
	reg        regionFunc // optional region object
	customList []string   // optional string slice object
}

type command struct {
	cmd    cmdFunc   // the command function
	parser parseFunc // a function to parse command arguments (if needed)
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

var cmdCharsNormalMode = map[string]command{
	"i": command{insertAtCs, nil},
	"a": command{appendAtCs, nil},
	"h": command{moveCursorLeft, nil},
	"j": command{moveCursorDown, nil},
	"k": command{moveCursorUp, nil},
	"l": command{moveCursorRight, nil},
	"d": command{delete_, parseRegion},
	"x": command{deleteCharForward, nil},
	"e": command{moveCursorTo, nil},
	"E": command{moveCursorTo, nil},
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
	KeyDelete:     command{deleteCharForward, nil},
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

func moveCursorTo(ctx *cmdContext) {
	if ctx.num == 0 {
		ctx.num = 1
	}
	ctx.reg = motions[ctx.cmdString]
	for i := 0; i < ctx.num; i++ {
		r := ctx.reg(*ctx.point)
		*ctx.point = r.end
	}
	ctx.point.buf.cs = *ctx.point
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

func deleteCharForward(ctx *cmdContext) {
	eng.deleteCharForward(*ctx.point)
	ctx.point.fixPos()
}
