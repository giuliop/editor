package main

import (
	"strings"
	"time"
)

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
			select {
			case ev = <-keyEvents:
			case <-time.After(5 * time.Minute):
				ev.Type = UIEventTimeout
			}
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

func parseAction(ev *UIEvent, ctx *cmdContext, cmds chan *cmdContext) (
	parser parseFunc, reprocess bool) {
	// if called by a timeout execute a matched command or wait for more
	switch {
	case ev.Type == UIEventTimeout:
		ctx.cmd = cmdCharsNormalMode[ctx.cmdString].cmd
		parser = cmdCharsNormalMode[ctx.cmdString].parser
	case ev.Special:
		// special keys immediately match a command, no ambiguity possible
		ctx.cmd = cmdKeys[ev.Buf.mod][ev.Key].cmd
		parser = cmdKeys[ev.Buf.mod][ev.Key].parser
	case ev.Buf.mod == insertMode:
		ctx.cmd = insertChar
	case ev.Buf.mod == normalMode:
		if isNumber(ev.Char, ctx) && ctx.cmdString == "" {
			loadNumber(ev.Char, ctx)
			return parseAction, false
		}
		ctx.cmdString += string(ev.Char)
		match, submatches := matchCommand(ctx.cmdString, ctx.customList,
			cmdCharsNormalMode)
		switch len(submatches) {
		case 0:
			return nil, false
		case 1:
			// if it is an exact match
			if match != "" {
				ctx.cmd = cmdCharsNormalMode[match].cmd
				parser = cmdCharsNormalMode[match].parser
			}
		default:
			ctx.customList = submatches
			return parseAction, false
		}
	}
	if ctx.cmd != nil {
		if parser == nil {
			cmds <- ctx
		}
	}
	// we will pass the baton to a new parser, let's set up ctx
	ctx.point = &ev.Buf.cs
	ctx.char = ev.Char
	ctx.customList = nil
	return parser, false
}

func matchCommand(s string, list []string, m map[string]command) (
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
		return
	}
	for _, str := range list {
		if strings.HasPrefix(str, s) {
			subMatches = append(subMatches, str)
			if len(str) == len(s) {
				match = s
			}
		}
	}
	return
}

func parseRegion(ev *UIEvent, ctx *cmdContext, cmds chan *cmdContext) (parseFunc, bool) {
	// we'll use ctx.customList to save the submatches so far
	switch {
	// if we were called by a timeout...
	case ev.Type == UIEventTimeout:
		ctx.reg = regionFuncs[ctx.argString]
		if ctx.reg != nil {
			cmds <- ctx
			return nil, false
		}
	// if we need to load a number...
	case isNumber(ev.Char, ctx) && ctx.argString == "":
		loadNumber(ev.Char, ctx)
	// otherwise we match the char
	default:
		ctx.argString += string(ev.Char)
		match, subMatches := matchRegionFunc(ctx.argString, ctx.customList, regionFuncs)
		switch len(subMatches) {
		case 0:
			return nil, false
		case 1:
			// if it is an exact match
			if match != "" {
				ctx.reg = regionFuncs[subMatches[0]]
				cmds <- ctx
				return nil, false
			}
		}
		ctx.customList = subMatches
	}
	return parseRegion, false
}

// stringMatch takes a string s and list of strings and returns two values:
// s (or "") if s has a perfect match in list (or not) and the sublist of
// strings in list that have s as a prefix
func matchRegionFunc(s string, list []string, m map[string]regionFunc) (
	match string, subMatches []string) {
	// if list is nil it is the first iteraction and we need to build it
	// from the map; s will be a single char
	if list == nil {
		for key := range m {
			if key[0] == s[0] {
				list = append(list, key)
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
	return
}
