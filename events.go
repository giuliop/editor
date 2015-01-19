package main

import (
	"strings"
	"time"
)

const keypressTimeout = 750 * time.Millisecond

func executeCommands(ui UI, cmds chan cmdContext) {
	defer cleanExit()
	for c := range cmds {
		c.cmd(&c)
		if !c.silent {
			ui.Draw()
		}
	}
}

func manageKeypress(ui UI, keys chan UIEvent, commands chan cmdContext) {
	defer cleanExit()
	var (
		nextParser     parseFunc = parseAction
		reconsumeEvent bool
		ev             UIEvent
	)
	reprocess := make(chan UIEvent)
	ctx := &cmdContext{}
	for {
		// first we check if we need to reprocess an old keypress, if not we either
		// wait for a new keypress or a timeout
		select {
		case ev = <-reprocess:
		default:
			select {
			case ev = <-keys:
				ctx.point = &ev.Buf.cs
			case <-time.After(keypressTimeout):
				ev.Type = UIEventTimeout
			}
		}
		//var reconsumeEvent bool
		nextParser, reconsumeEvent = nextParser(&ev, ctx, commands)
		if nextParser == nil {
			ctx = &cmdContext{}
			nextParser = parseAction
		}
		if reconsumeEvent {
			reprocess <- ev
		}
	}
}

func parseAction(ev *UIEvent, ctx *cmdContext, cmds chan cmdContext) (
	nextParser parseFunc, reprocessEvent bool) {
	switch {
	// if called by a timeout execute a matched string command if we have one
	case ev.Type == UIEventTimeout:
		if ctx.cmdString != "" {
			c := lookupStringCmd(ctx.point.buf.mod, ctx.cmdString)
			if c.cmd != nil {
				return pushCmd(c, ctx, cmds), false
			}
			if ctx.point.buf.mod == insertMode {
				ctx.cmdString = ""
			}
		}
		return parseAction, false
	case ev.Key.isSpecial:
		c := lookupStringCmd(ev.Buf.mod, ctx.cmdString)
		// if we have a valid command in the pipeline we'll execute it and reprocess
		// the special key at next iteration
		if c.cmd != nil {
			reprocessEvent = true
		} else {
			c = lookupKeyCmd(ev.Buf.mod, ev.Key.Special)
		}
		return pushCmd(c, ctx, cmds), reprocessEvent
	case isNumber(ev.Key.Char, ctx) && ctx.cmdString == "":
		loadNumber(ev.Key.Char, ctx)
		return parseAction, false
	default:
		m := ev.Buf.mod
		ctx.char = ev.Key.Char
		ctx.cmdString += string(ctx.char)
		c, submatches := matchCommand(ev.Buf.mod, ctx.cmdString, ctx.customList)
		ctx.customList = submatches
		if m == insertMode {
			// we insert the char and just delete it later if we match a command
			pushCmd(command{insertChar, nil}, ctx, cmds)
		}
		switch {
		case len(submatches) == 0:
			// if no matches, we check if we had a valid command before this char
			// and if so execute the command and reprocess the char
			if m == normalMode {
				c = lookupStringCmd(ev.Buf.mod, ctx.cmdString[:len(ctx.cmdString)-1])
				if c.cmd != nil {
					return pushCmd(c, ctx, cmds), true
				}
			}
			return nil, false
		case len(submatches) == 1 && c.cmd != nil:
			if m == insertMode {
				ctx.cmd = deleteCharBackward
				ctx.silent = true
				for i := 0; i < len(ctx.cmdString); i++ {
					cmds <- *ctx
				}
				ctx.silent = false
			}
			return pushCmd(c, ctx, cmds), false
		default:
			return parseAction, false
		}
	}
}

func pushCmd(c command, ctx *cmdContext, cmds chan cmdContext) parseFunc {
	if c.cmd != nil {
		ctx.cmd = c.cmd
		if c.parser == nil {
			cmds <- *ctx
		}
	}
	ctx.customList = nil
	return c.parser
}

func matchCommand(mod mode, s string, list []string) (
	match command, subMatches []string) {
	// if list is nil it is the first iteraction and we need to build it
	// from the appropriate command map; s will be a single char
	m := cmdStringTables[mod]
	if list == nil {
		for key := range m {
			if key[0] == s[0] {
				subMatches = append(subMatches, key)
				// if exact match
				if len(key) == 1 {
					match = m[key]
				}
			}
		}
	} else {
		for _, str := range list {
			if strings.HasPrefix(str, s) {
				subMatches = append(subMatches, str)
				// if exact match
				if len(str) == len(s) {
					match = m[s]
				}
			}
		}
	}
	return match, subMatches
}

func parseRegion(ev *UIEvent, ctx *cmdContext, cmds chan cmdContext) (
	nextParser parseFunc, reprocessEvent bool) {
	switch {
	// if called by a timeout execute a matched string command if we have one
	case ev.Type == UIEventTimeout:
		if ctx.argString != "" {
			ctx.reg = regionFuncs[ctx.argString]
			if ctx.reg != nil {
				cmds <- *ctx
				return nil, false
			}
		}
		return parseRegion, false
	case ev.Key.isSpecial:
		ctx.reg = regionFuncs[ctx.argString]
		// if we have a valid region in the pipeline we'll execute it; in any
		// case we reset parsing and reprocess the event
		if ctx.reg != nil {
			cmds <- *ctx
		}
		return nil, true
	case isNumber(ev.Key.Char, ctx) && ctx.argString == "":
		loadNumber(ev.Key.Char, ctx)
		return parseRegion, false
	default:
		ctx.argString += string(ev.Key.Char)
		match, subMatches := matchRegionFunc(ctx.argString, ctx.customList, regionFuncs)
		ctx.customList = subMatches
		switch len(subMatches) {
		case 0:
			return nil, false
		case 1:
			// if it is an exact match
			if match != "" {
				ctx.reg = regionFuncs[subMatches[0]]
				cmds <- *ctx
				return nil, false
			}
		default:
		}
		return parseRegion, false
	}
}

func matchRegionFunc(s string, list []string, m map[string]regionFunc) (
	match string, subMatches []string) {
	// if list is nil it is the first iteraction and we need to build it
	// from the map; s will be a single char
	if list == nil {
		for key := range m {
			if key[0] == s[0] {
				subMatches = append(subMatches, key)
				// if exact match
				if len(key) == 1 {
					match = key
				}
			}
		}
	} else {
		for _, str := range list {
			if strings.HasPrefix(str, s) {
				subMatches = append(subMatches, str)
				if len(str) == len(s) {
					match = s
				}
			}
		}
	}
	debug.Printf("match %v, subMatches %v\n", match, subMatches)
	return
}
