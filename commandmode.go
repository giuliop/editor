package main

import (
	"fmt"
	"strings"
)

type commandModeF func(args []string) (msg string)

var commandModeFuncs = map[string]commandModeF{
	"echo": echo,
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
			return nil, false
		case KeyArrowRight:
		case KeyArrowLeft:
		case KeyTab:
		case KeyBackspace, KeyBackspace2:
		case KeyDelete:
		case KeyEsc, KeyCtrlC:
			return nil, false
		case KeySpace:
		}
	default:
		ctx.argString += string(ev.Key.Char)
	}
	ui.Draw()
	return parseCommandMode, false
}
