package main

func manageEventKey(ui UI, keyEvents chan UIEvent) {
	ctx := &cmdContext{}
	for ev := range keyEvents {
		b := ev.Buf
		mod := b.mod
		if ctx.cmd == nil && !ctx.save {
			ctx = &cmdContext{point: &(b.cs)}
		}

		if ev.Char != 0 {
			ctx.key = ev.Char
			if ctx.cmd == nil {
				switch mod {
				case insertMode:
					ctx.cmd = insertChar
				case normalMode:
					ctx.cmd = cmdCharsNormalMode[ev.Char]
				}
			}
		} else {
			ctx.cmd = cmdKeys[mod][ev.Key]
		}

		if ctx.cmd != nil {
			cmd := ctx.cmd
			ctx.cmd = nil
			ctx.save = false
			cmd(ctx)
		}
		ui.Draw()
	}
}
