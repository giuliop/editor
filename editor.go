package main

import (
	"io/ioutil"
)

var (
	ui  UI
	eng textEngine
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
	eng = initEngine()

	// initialize ui frontend
	var err error
	ui, err = selectUI("terminal")
	check(err)
	check(ui.Init())
	defer ui.Close()
	ui.Draw()

eventLoop:
	for {
		switch ev := ui.PollEvent(); ev.Type {
		case UIEventKey:
			switch ev.Key {
			case KeyEsc:
				break eventLoop
			case KeyBackspace, KeyBackspace2:
				eng.deleteChBackward()
			case KeyTab:
				eng.insertCh('\t')
			case KeySpace:
				eng.insertCh(' ')
			case KeyEnter, KeyCtrlJ:
				eng.insertNewLineCh()
			default:
				if ev.Ch != 0 {
					eng.insertCh(ev.Ch)
				}
			}
		case UIEventError:
			panic(ev.Err)
		}
		ui.Draw()
	}
}
