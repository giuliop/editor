package main

import (
	"io/ioutil"
)

var (
	ui UI
	in internal
)

// check panics if passed an error
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// debug writes msg to a file called debug and optionally panics based on the value of stop
func debug(stop bool, msg string) {
	err := ioutil.WriteFile("debug", []byte(msg), 0644)
	check(err)
	if stop {
		panic(msg)
	}
}

func main() {
	// initialize internal editor
	// for now we create one empty buffer
	in = initInternal()

	// initialize ui frontend
	var err error
	ui, err = selectUI("terminal")
	check(err)
	check(ui.Init())
	defer ui.Close()
	draw()

eventLoop:
	for {
		switch ev := ui.PollEvent(); ev.Type {
		case UIEventKey:
			switch ev.Key {
			case KeyEsc:
				break eventLoop
			case KeyBackspace, KeyBackspace2:
				in.cb.deleteChBackward()
			case KeyTab:
				in.cb.insertChar('\t')
			case KeySpace:
				in.cb.insertChar(' ')
			case KeyEnter, KeyCtrlJ:
				in.cb.insertNewLineChar()
			default:
				if ev.Ch != 0 {
					in.cb.insertChar(ev.Ch)
				}
			}
		case UIEventError:
			panic(ev.Err)
		}
		draw()
	}
}
