package main

import (
	"fmt"
	"strings"
)

const (
	commandModePrompt  = "-> "
	commandModeMaxCmds = 100
)

type commandModeF func(args []string) (msg string)

type commandRegister struct {
	list []string // the list of commands
	last int      // to retrieve past commands
}

func (c *commandRegister) add(cmd string) {
	c.list = append(c.list, cmd)
	c.last++
	debug.Println(c.last)
	if len(c.list) > commandModeMaxCmds {
		c.list = c.list[1:]
		c.last--
	}
}

func (c *commandRegister) previous() {
	be.msgLine = append(be.msgLine[:len(commandModePrompt)],
		stringToLine(c.list[c.last])...)
	if c.last > 0 {
		c.last--
	}
}

func (c *commandRegister) next() {
	be.msgLine = be.msgLine[:len(commandModePrompt)]
	if c.last == len(c.list)-1 {
		return
	}
	c.last++
	be.msgLine = append(be.msgLine, stringToLine(c.list[c.last])...)
}

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
}

func exitCommandMode() {
	r.commands.last = len(r.commands.list) - 1
	be.commandMode = false
	ui.Draw()
}

func enterCommand(s string) (msg string) {
	if s == "" {
		return ""
	}
	r.commands.last = len(r.commands.list) - 1
	r.commands.add(s)
	tokens := strings.Split(s, " ")
	cmd, args := tokens[0], tokens[1:]
	f := commandModeFuncs[cmd]
	if f == nil {
		return fmt.Sprintf("Unknown command: %v", cmd)
	}
	msg = f(args)

	// make sure the cursor is valid in case the command changed the buffer
	cs := ui.CurrentView().cs
	if cs.line > len(ui.CurrentView().buf.text)-1 {
		cs.line = len(ui.CurrentView().buf.text) - 1
	}
	cs.fixPos()

	return msg
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
			return nil, false
		case KeyArrowRight:
		case KeyArrowLeft:
		case KeyArrowDown:
			r.commands.next()
		case KeyArrowUp:
			r.commands.previous()
		case KeyTab:
		case KeyBackspace, KeyBackspace2:
			if len(be.msgLine) > len(commandModePrompt) {
				be.msgLine = be.msgLine[:len(be.msgLine)-1]
			}
		case KeyDelete:
		case KeyEsc, KeyCtrlC:
			be.msgLine = line{}
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
