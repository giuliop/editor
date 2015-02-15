package main

type direction int

const (
	right direction = iota
	left
	up
	down
)

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
	silent     bool       // if true does not redraw the screen after execution
	msg        string     // to comunicate back to user
	cmdChans   cmdStack   // channels to push the command and wait for done signal
}

type command struct {
	cmd    cmdFunc   // the command function
	parser parseFunc // a function to parse command arguments (if needed)
}
type cmdFunc func(ctx *cmdContext)
type parseFunc func(ev *UIEvent, ctx *cmdContext) (parseFunc, bool)

var cmdStringTables = [2]map[string]command{cmdStringInsertMode, cmdStringNormalMode}
var cmdKeyTables = [2]map[Key]command{cmdKeyInsertMode, cmdKeyNormalMode}

func lookupStringCmd(m mode, s string) command {
	return cmdStringTables[m][s]
}

func lookupKeyCmd(m mode, key Key) command {
	return cmdKeyTables[m][key]
}

var cmdKeyNormalMode = map[Key]command{
	KeyCtrlS: command{saveToFile, nil},
	KeyCtrlX: command{exitProgram, nil},
	KeyCtrlR: command{redo, nil},
}

// commands should be at most two chars to avoid risk of over-shadowing one char
// command (e.g., 'dgg' could overshadow command 'd')
var cmdStringNormalMode = map[string]command{
	",q": command{exitProgram, nil},
	"i":  command{insertAtCs, nil},
	"a":  command{appendAtCs, nil},
	"A":  command{appendAtEndOfLine, nil},
	"h":  command{moveCursorLeft, nil},
	"j":  command{moveCursorDown, nil},
	"k":  command{moveCursorUp, nil},
	"l":  command{moveCursorRight, nil},
	"d":  command{delete_, parseRegion},
	"dd": command{deleteLine, nil},
	"x":  command{deleteCharForward, nil},
	"e":  command{moveCursorTo, nil},
	"E":  command{moveCursorTo, nil},
	"B":  command{moveCursorTo, nil},
	"b":  command{moveCursorTo, nil},
	"w":  command{moveCursorTo, nil},
	"W":  command{moveCursorTo, nil},
	"0":  command{moveCursorTo, nil},
	"$":  command{moveCursorTo, nil},
	"H":  command{moveCursorTo, nil},
	"L":  command{moveCursorTo, nil},
	"gg": command{moveCursorTo, nil},
	"G":  command{moveCursorTo, nil},
	"m":  command{recordMacro, nil},
	"u":  command{undo, nil},
}

var cmdKeyInsertMode = map[Key]command{
	KeyEsc:        command{toNormalMode, nil},
	KeyBackspace:  command{deleteCharBackward, nil},
	KeyBackspace2: command{deleteCharBackward, nil},
	KeyTab:        command{insertTab, nil},
	KeySpace:      command{insertSpace, nil},
	KeyEnter:      command{insertNewLine, nil},
	KeyCtrlJ:      command{insertNewLine, nil},
	KeyCtrlC:      command{toNormalMode, nil},
	KeyDelete:     command{deleteCharForward, nil},
	KeyCtrlS:      command{saveToFile, nil},
}

var cmdStringInsertMode = map[string]command{
	"AA": command{appendAtEndOfLine, nil},
}

func toNormalMode(ctx *cmdContext) {
	ctx.point.buf.mod = normalMode
	if !ctx.point.atLineStart() {
		ctx.point.moveLeft(1)
	}
	ctx.msg = "Normal mode"
}

func insertAtCs(ctx *cmdContext) {
	ctx.point.buf.mod = insertMode
}

func appendAtCs(ctx *cmdContext) {
	// move cursor right unless empty line
	ctx.point.buf.mod = insertMode
	if !ctx.point.atEmptyLine() {
		ctx.point.moveRight(1)
	}
}

func appendAtEndOfLine(ctx *cmdContext) {
	// move cursor right unless empty line
	ctx.point.buf.mod = insertMode
	ctx.point.pos = ctx.point.maxCursPos()
}

func moveCursorLeft(ctx *cmdContext) {
	ctx.point.moveLeft(ctx.num)
}

func moveCursorRight(ctx *cmdContext) {
	ctx.point.moveRight(ctx.num)
}

func moveCursorUp(ctx *cmdContext) {
	ctx.point.moveUp(ctx.num)
}

func moveCursorDown(ctx *cmdContext) {
	ctx.point.moveDown(ctx.num)
}

func moveCursorTo(ctx *cmdContext) {
	ctx.reg = motions[ctx.cmdString]
	for i := 0; i < ctx.num; i++ {
		r, _ := ctx.reg(*ctx.point)
		*ctx.point = r.end
	}
	ctx.point.buf.cs = *ctx.point
}

func delete_(ctx *cmdContext) {
	switch ctx.argString {
	case "gg":
		deleteToStart(ctx)
	case "G":
		deleteToEnd(ctx)
	default:
		for i := 0; i < ctx.num; i++ {
			r, dir := ctx.reg(*ctx.point)
			if (ctx.argString == "W" || ctx.argString == "w") &&
				!r.end.atLastTextChar() {
				r.end.pos--
			}
			*ctx.point = r.delete(dir)
		}
		ctx.point.buf.cs = *ctx.point
	}
}

func deleteLine(ctx *cmdContext) {
	p := ctx.point
	toline := p.line + ctx.num - 1
	if toline > p.maxLine() {
		toline = p.maxLine()
	}

	// add undo info
	start, end := mark{p.line, 0, p.buf}, mark{toline, 0, p.buf}
	end.pos = end.lineEndPos()
	p.buf.changeList.add(*ctx, undoContext{undoDelete, start.copy(end), start, mark{}})

	p.buf.deleteLines(*p, mark{toline, 0, p.buf})
	if p.line > p.maxLine() {
		p.line--
	}
	p.fixPos()
}

func deleteToStart(ctx *cmdContext) {
	b := ctx.point.buf
	b.deleteLines(mark{0, 0, b}, *ctx.point)
	b.cs = mark{0, 0, b}
}

func deleteToEnd(ctx *cmdContext) {
	b := ctx.point.buf
	b.deleteLines(*ctx.point, mark{ctx.point.lastLine(), 0, b})
	b.cs = mark{ctx.point.line - 1, 0, b}
	if b.cs.line < 0 {
		b.cs.line = 0
	}
}

func exitProgram(ctx *cmdContext) {
	exit <- true
}

func deleteCharForward(ctx *cmdContext) {
	for i := 0; i < ctx.num; i++ {
		ctx.point.deleteCharForward()
		ctx.point.fixPos()
	}
}

func deleteCharBackward(ctx *cmdContext) {
	*ctx.point = ctx.point.deleteCharBackward()
}

func insertTab(ctx *cmdContext) {
	ctx.point.insertChar('\t')
	ctx.point.moveRight(1)
}

func insertSpace(ctx *cmdContext) {
	ctx.point.insertChar(' ')
	ctx.point.moveRight(1)
}

func insertNewLine(ctx *cmdContext) {
	ctx.point.insertNewLineChar()
	ctx.point.set(ctx.point.line+1, 0)
}

func insertChar(ctx *cmdContext) {
	ctx.point.insertChar(ctx.char)
	ctx.point.moveRight(1)
}

func saveToFile(ctx *cmdContext) {
	err := ctx.point.buf.save()
	if err != nil {
		ctx.msg = err.Error()
	} else {
		ctx.msg = "file saved"
	}
}
