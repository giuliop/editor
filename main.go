package main

import (
	"os"
)

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

// log writes msg to file log
func log(msg string) {
	f, err := os.OpenFile("log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	defer f.Close()
	check(err)
	_, err = f.WriteString(msg + "\n")
	check(err)
}

func main() {
	// initialize internal engine and create an empty buffer as current buffer
	eng := *initTextEngine()
	b := eng.newBuffer("")

	// initialize ui frontend
	ui, err := selectUI("terminal")
	check(err)
	err = ui.Init(b)
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
		// activate key command manager
		keyEvents := make(chan UIEvent, 100)
		cmdToExecute := make(chan cmdContext, 10)
		go manageEventKey(ui, keyEvents, cmdToExecute)
		go executeCommands(ui, cmdToExecute)
		// listen for events and route them to appropriate channel
		for ev := range uiEvents {
			switch ev.Type {
			case UIEventKey:
				keyEvents <- ev
			case UIEventError:
				check(ev.Err)
			}
		}
	}()

	// wait for exit signal
	<-exitSignal
}
