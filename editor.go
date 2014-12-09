package main

import (
	"io/ioutil"
)

var (
	ui UI
	in internalEditor
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
	ui.Draw()

eventLoop:
	for {
		switch ev := ui.PollEvent(); ev.Type {
		case UIEventKey:
			switch ev.Key {
			case KeyEsc:
				break eventLoop
			case KeyBackspace, KeyBackspace2:
				in.deleteChBackward()
			case KeyTab:
				in.insertCh('\t')
			case KeySpace:
				in.insertCh(' ')
			case KeyEnter, KeyCtrlJ:
				in.insertNewLineCh()
			default:
				if ev.Ch != 0 {
					in.insertCh(ev.Ch)
				}
			}
		case UIEventError:
			panic(ev.Err)
		}
		ui.Draw()
	}
}
