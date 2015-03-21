package main

import "strings"

const (
	commandModePrompt  = "-> "
	commandModeMaxCmds = 100
)

type commandModeF func(args []string) (msg string)

type commandRegister struct {
	list    []line // the list of commands
	last    int    // to retrieve past commands
	current line   // the current command being typed
}

func (c *commandRegister) reset() {
	c.list = c.list[:0]
	c.last = -1
	c.current = c.current[:0]
}

func (c *commandRegister) add(cmd line) {
	c.list = append(c.list, cmd)
	c.last++
	if len(c.list) > commandModeMaxCmds {
		c.list = c.list[1:]
		c.last--
	}
}

func (c *commandRegister) previous() {
	for found := false; found == false; {
		if c.list[c.last].hasPrefix(c.current) {
			be.msgLine = append(be.msgLine[:len(commandModePrompt)],
				c.list[c.last]...)
			found = true
		}
		if c.last > 0 {
			c.last--
		} else {
			return
		}
	}
}

func (c *commandRegister) next() {
	for {
		if c.last == len(c.list)-1 {
			be.msgLine = append(be.msgLine[:len(commandModePrompt)], c.current...)
			return
		}
		c.last++
		if c.list[c.last].hasPrefix(c.current) {
			be.msgLine = append(be.msgLine[:len(commandModePrompt)],
				c.list[c.last]...)
			return
		}
	}
}

var commandModeFuncs = map[string]commandModeF{
	"q":    quit,
	"echo": echo,
}

func echo(args []string) (msg string) {
	return strings.Join(args, " ")
}

func quit(args []string) (msg string) {
	exitProgram(nil)
	return "Bye-bye"
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
	r.commands.current = r.commands.current[:0]
	be.commandMode = false
	ui.Draw()
}

func enterCommand(cmd line) (msg string) {
	if len(cmd) == 0 {
		return ""
	}
	r.commands.add(cmd)
	tokens := strings.Split(string(cmd), " ")
	c, args := tokens[0], tokens[1:]
	f := commandModeFuncs[c]
	if f == nil {
		return "Unknown command: " + string(cmd)
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

func parseCommandMode(ev *UIEvent, ctx *cmdContext) (
	nextParser parseFunc, reprocessEvent bool) {
	switch {
	case ev.Type == UIEventTimeout:
		return parseCommandMode, false
	case ev.Key.isSpecial:
		switch ev.Key.Special {
		case KeyCtrlJ, KeyEnter:
			be.msgLine = stringToLine(
				enterCommand(be.msgLine[len(commandModePrompt):]))
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
				r.commands.current = append(line{},
					be.msgLine[len(commandModePrompt):]...)
			}
		case KeyDelete:
		case KeyEsc, KeyCtrlC:
			be.msgLine = be.msgLine[:0]
			exitCommandMode()
			return nil, false
		case KeySpace:
			be.msgLine = append(be.msgLine, ' ')
		case testEndOfEmission: // to support automated testing
			testChan <- struct{}{}
		}

	default:
		be.msgLine = append(be.msgLine, ev.Key.Char)
		r.commands.current = append(line{}, be.msgLine[len(commandModePrompt):]...)
	}
	ui.Draw()
	return parseCommandMode, false
}
