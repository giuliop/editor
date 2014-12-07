package main

import (
	"github.com/mattn/go-runewidth"
)

func draw() {
	ui.Clear()
	// viPos tracks the visual position of chars in the line since some chars
	// might take two spaces on screen
	var viPos int
	for i, line := range text {
		viPos = 0
		for _, ch := range line {
			ui.SetCell(viPos, i, ch)
			viPos += runewidth.RuneWidth(ch)
		}
	}
	ui.SetCursor(viPos, cs.line)
	ui.Flush()
}
