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
			switch ev.Type {
			case UIEventKey:
				ctx := new(cmdContext)
				ctx.point = &(ui.CurrentBuffer().cs)
				cmd, ok := cmdKeys[ev.Key]
				if !ok && ev.Char != 0 {
					cmd = insertChar
					ctx.char = ev.Char
				}
				cmd(ctx)
			case UIEventError:
				panic(ev.Err)
			}
			ui.Draw()
		}
	}()

	// wait for exit signal
	<-exitSignal
}
