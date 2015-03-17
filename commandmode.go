package main

import (
	"fmt"
	"strings"
)

var commandModePrompt = "-> "

type commandModeF func(args []string) (msg string)

var commandModeFuncs = map[string]commandModeF{
	"echo": echo,
}

func initCommandView() *view {
	buf := be.newBuffer("")
	return &view{buf, newMark(buf), 0}
}

func enterCommandMode(ctx *cmdContext) {
	be.commandMode = true
	ctx.msg = commandModePrompt
	debug.Println(be.msgLine)
}

func exitCommandMode() {
	be.commandMode = false
}

func enterCommand(s string) (msg string) {
	tokens := strings.Split(s, " ")
	cmd, args := tokens[0], tokens[1:]
	f := commandModeFuncs[cmd]
	if f == nil {
		return fmt.Sprintf("Unknown command: %v", cmd)
	}
	return f(args)
}

func echo(args []string) (msg string) {
	return strings.Join(args, " ")
}

func parseCommandMode(ev *UIEvent, ctx *cmdContext) (
	nextParser parseFunc, reprocessEvent bool) {
	switch {
	case ev.Type == UIEventTimeout:
		return parseCommandMode, false
	case ev.Key.isSpecial:
		switch ev.Key.Special {
		case KeyCtrlJ, KeyEnter:
			be.msgLine = stringToLine(enterCommand(string(
				be.msgLine[len(commandModePrompt):])))
			exitCommandMode()
			ui.Draw()
			return nil, false
		case KeyArrowRight:
		case KeyArrowLeft:
		case KeyTab:
		case KeyBackspace, KeyBackspace2:
			if len(be.msgLine) > 0 {
				be.msgLine = be.msgLine[:len(be.msgLine)-1]
			}
		case KeyDelete:
		case KeyEsc, KeyCtrlC:
			exitCommandMode()
			return nil, false
		case KeySpace:
			be.msgLine = append(be.msgLine, ' ')
		}
	default:
		be.msgLine = append(be.msgLine, ev.Key.Char)
	}
	ui.Draw()
	return parseCommandMode, false
}
