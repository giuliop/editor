// Editor is a great editor, or at least it will be one day!
package main

import "os"

var (
	be   backend           // the open buffers collection backend
	ui   UI                // the user interface
	r    = register{}      // holds all global lists: macros...
	exit = make(chan bool) // a channel to signal quitting the program
)

var debug *debugLogger

type register struct {
	macro *macroRegister
}

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

func init() {
	// initialize debug logging, available through a logger called debug
	debug = initDebug()
	// initialize internal engine
	be = initBackend()
	// initialize the global registers
	r = initRegisters()
}

func initRegisters() register {
	r := register{}
	r.macro = &macroRegister{&keyLogger{}, [10][]Keypress{}}
	return r
}

func main() {
	defer debug.stop()
	// initialize the user interface
	curBuf := be.open(os.Args[1:])
	var err error
	ui, err = initFrontEnd(curBuf)
	check(err)
	ui.Draw()
	defer ui.Close()

	//activate channels for keypresses and recognized commands
	keys := make(chan UIEvent, 100)
	commands := make(chan cmdContext)
	go manageKeypress(keys, commands)
	go executeCommands(commands)

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
