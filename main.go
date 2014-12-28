package main

import "os"

var (
	eng        textEngine        // the buffer collection backend
	exitSignal = make(chan bool) // a channel to signal quitting the program
)

// check panics if passed an error
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// debug writes msg to a file called debug and optionally panics based on the value of stop
func debug(stop bool, msg string) {
	f, err := os.OpenFile("debug", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	defer f.Close()
	check(err)
	_, err = f.WriteString(msg + "\n")
	check(err)
	if stop {
		panic(msg)
	}
}
func init() {
	initCmdTables()
}

func main() {
	// initialize internal engine
	eng = initEngine()

	// create an empty buffer as current buffer
	b := eng.newBuffer("")

	// initialize ui frontend
	ui, err := selectUI("terminal")
	check(err)
	ui.Init(b)
	check(err)
	defer ui.Close()
	ui.Draw()

	// activate channel for IO events
	uiEvents := make(chan UIEvent, 100)
	go func() {
		for {
			uiEvents <- ui.PollEvent()
		}
	}()

	//activate command manager
	go func() {
		for ev := range uiEvents {
			b := ui.CurrentBuffer()
			switch ev.Type {
			case UIEventKey:
				var cmd cmdFunc
				ctx := new(cmdContext)
				ctx.point = &(b.cs)
				if ev.Char != 0 {
					ctx.char = ev.Char
					if b.mod == insertMode {
						cmd = insertChar
					}
					if b.mod == normalMode {
						cmd = cmdCharsNormalMode[ev.Char]
					}
				} else {
					cmd = cmdKeys[b.mod][ev.Key]
				}
				if cmd != nil {
					cmd(ctx)
				}
			case UIEventError:
				check(ev.Err)
			}
			ui.Draw()
		}
	}()

	// wait for exit signal
	<-exitSignal
}
