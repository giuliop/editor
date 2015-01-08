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
		nextParser, reconsumeEvent = nextParser(&ev, ctx, commands)
		if reconsumeEvent {
			reprocess <- ev
		}
	}
}

func parseAction(ev *UIEvent, ctx *cmdContext, cmds chan *cmdContext) (parseFunc, bool) {
	// ev == nil means we were called for a timeout
	if ev == nil {
		return parseAction, false
	}
	b := ev.Buf
	mode := b.mod
	var parser parseFunc
	if ev.Special {
		ctx.cmd = cmdKeys[mode][ev.Key].cmd
		parser = cmdKeys[mode][ev.Key].parser
		// else a char was pressed
	} else {
		switch mode {
		case insertMode:
			// we insert the char in insertMode
			ctx.cmd = insertChar
			// we look up the char command in normalMode
		case normalMode:
			if isNumber(ev.Char, ctx) && ctx.cmd == nil {
				loadNumber(ev.Char, ctx)
				return parseAction, false
			}
			ctx.cmd = cmdCharsNormalMode[ev.Char].cmd
			parser = cmdCharsNormalMode[ev.Char].parser
		}
	}
	if ctx.cmd != nil {
		ctx.point = &b.cs
		ctx.char = ev.Char
		if parser == nil {
			cmds <- ctx
		}
	}
	return parser, false
}

func parseRegion(ev *UIEvent, ctx *cmdContext, cmds chan *cmdContext) (parseFunc, bool) {
	// ev == nil means we were called for a timeout
	// we'll use ctx.custom to build list of candiates
	b := ev.Buf
	if ev == nil {
		if ctx.custom != "" {
			ctx.point = &b.cs
			ctx.reg = regionFuncs[ctx.custom]
			cmds <- ctx
			return nil, false
		} else {
			return parseRegion, false
		}
	}
	if isNumber(ev.Char, ctx) && ctx.custom == "" {
		loadNumber(ev.Char, ctx)
		return parseRegion, false
	}
	candidates := make([]string, 0, 20)
	if ctx.custom == "" {
		for key := range regionFuncs {
			if rune(key[0]) == ev.Char {
				candidates = append(candidates, key)
			}
		}
	} else {
		sofar := ctx.custom + string(ev.Char)
		for _, c := range ctx.customList {
			if c[:len(sofar)] == sofar {
				candidates = append(candidates, c)
			}
		}
	}
	switch len(candidates) {
	case 0:
		return nil, false
	case 1:
		ctx.point = &b.cs
		ctx.reg = regionFuncs[candidates[0]]
		cmds <- ctx
		return nil, false
	default:
		ctx.customList = candidates
		return parseRegion, false
	}
}
