// Editor is a great editor, or at least it will be one day!
package main

import "os"

var (
	be   backend           // the open buffers collection backend
	exit = make(chan bool) // a channel to signal quitting the program
)

// check panics if passed an error
func check(e error) {
	if e != nil {
		panic(e)
	}
}

func forceExit() {
	debug.Print(" * Force Exit * \n\n")
	//debug.Printf("%+v", *ctx)
	exit <- true
}

func fatalError(e error) {
	debug.Print(" * Fatal error * \n\n")
	debug.Println(e)
	exit <- true
}

func cleanupOnError() {
	if r := recover(); r != nil {
		debug.Print(" * Fatal error * \n\n")
		debug.Println(r)
		debug.printStack()
		exit <- true
	}
}

func initFrontEnd(activeBuf *buffer) (UI, error) {
	ui, err := selectUI("terminal")
	if err == nil {
		err = ui.Init(activeBuf)
	}
	return ui, err
}

func main() {
	// initialize debug logging, available through a logger called debug
	debug = initDebug()
	defer debug.stop()

	// initialize internal engine and create an empty buffer as current buffer
	be := initBackend()
	curBuf := be.open(os.Args[1:])

	// initialize ui frontend with the new empty buffer as active buffer
	ui, err := initFrontEnd(curBuf)
	defer ui.Close()
	check(err)
	ui.Draw()

	//activate channels for keypresses and recognized commands
	keys := make(chan UIEvent, 100)
	commands := make(chan cmdContext, 100)
	go manageKeypress(ui, keys, commands)
	go executeCommands(ui, commands)

	// listen for events and route them to appropriate channel
	uiEvents := make(chan UIEvent, 100)
	go func() {
		for {
			uiEvents <- ui.PollEvent()
		}
	}()
	go func() {
		for ev := range uiEvents {
			switch ev.Type {
			case UIEventKey:
				if ev.Key.Special == KeyF1 {
					forceExit()
				}
				keys <- ev
			case UIEventError:
				check(ev.Err)
			}
		}
	}()

	// wait for exit signal
	<-exit
}
