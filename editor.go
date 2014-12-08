package main

import (
	"io/ioutil"
)

var (
	ui   Ui
	cs   *cursor
	text []line
)

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func debug(stop bool, msg string) {
	err := ioutil.WriteFile("debug", []byte(msg), 0644)
	check(err)
	if stop {
		panic(msg)
	}
}

func main() {
	var err error
	ui, err = selectUI("terminal")
	check(err)
	check(ui.Init())
	defer ui.Close()
	text = make([]line, 1, 20)
	text[0] = newLine()
	cs = &cursor{0, 0}
	draw()

mainloop:
	for {
		switch ev := ui.PollEvent(); ev.Type {
		case UiEventKey:
			switch ev.Key {
			case KeyEsc:
				break mainloop
			case KeyBackspace, KeyBackspace2:
				deleteChBackward()
			case KeyTab:
				insertChar('\t')
			case KeySpace:
				insertChar(' ')
			case KeyEnter, KeyCtrlJ:
				insertNewLineChar()
			default:
				if ev.Ch != 0 {
					insertChar(ev.Ch)
				}
			}
		case UiEventError:
			panic(ev.Err)
		}
		draw()
	}
}
