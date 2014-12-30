package main

func manageEventKey(ui UI, keyEvents chan UIEvent) {
	var (
		cmd     cmdFunc
		ctx     *cmdContext
		cmdDone bool = true
	)
	for ev := range keyEvents {
		b := ev.Buf
		mod := b.mod
		if cmdDone {
			ctx = &cmdContext{point: &(b.cs)}
		}

		if ev.Char != 0 {
			ctx.key = ev.Char
			if cmdDone {
				switch mod {
				case insertMode:
					cmd = insertChar
				case normalMode:
					cmd = cmdCharsNormalMode[ev.Char]
				}
			}
		} else {
			cmd = cmdKeys[mod][ev.Key]
		}

		if cmd != nil {
			cmdDone = cmd(ctx)
		}
		ui.Draw()
	}
}
