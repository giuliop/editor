package main

func executeCommands(ui UI, cmds chan *cmdContext) {
	for c := range cmds {
		c.cmd(c)
		ui.Draw()
	}
}

func manageEventKey(ui UI, keyEvents chan UIEvent, commands chan *cmdContext) {
	var (
		ctx        *cmdContext = &cmdContext{}
		nextParser parseFunc   = parseAction
		ev         UIEvent
	)
	reprocess := make(chan UIEvent)
	for {
		select {
		case ev = <-reprocess:
		default:
			ev = <-keyEvents
		}
		// If nextPerser is nil we just completed a command and will look for
		// the next command; if a special key was pressed we always treat it
		// as a new command even if we were in the middle of another command
		if nextParser == nil || ev.Special {
			ctx = &cmdContext{}
			nextParser = parseAction
		}
		var reconsumeEvent bool
		nextParser, reconsumeEvent = nextParser(ev, ctx, commands)
		if reconsumeEvent {
			reprocess <- ev
		}
	}
}

func parseAction(ev UIEvent, ctx *cmdContext, cmds chan *cmdContext) (parseFunc, bool) {
	b := ev.Buf
	mode := b.mod
	if ev.Special {
		ctx.cmd = cmdKeys[mode][ev.Key]
		// else a char was pressed
	} else {
		switch mode {
		// we insert it in insertMode
		case insertMode:
			ctx.cmd = insertChar
			// if no cmd is waiting for input we look up the char command in
			// normalMode; otherwise the previous cmd will receive it
		case normalMode:
			ctx.cmd = cmdCharsNormalMode[ev.Char]
			if isNumber(ev.Char, ctx) && ctx.cmd == nil {
				loadNumber(ev.Char, ctx)
				return parseAction, false
			}
		}
	}
	if ctx.cmd != nil {
		ctx.point = &b.cs
		ctx.char = ev.Char
		cmds <- ctx
	}
	return nil, false
}
