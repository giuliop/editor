package main

import "testing"

type sample struct {
	buf  *buffer
	text string
}

var samples = []string{
	defaultText,
	emptyText,
	emptyLinesText,
}

func TestMotions(t *testing.T) {
	a := &asserter{}
	for _, s := range samples {
		b := stringToBuffer(s)
		e := newKeyPressEmitter(b)
		e.emit(KeyCtrlC, "G")
		a.assert("G", "cs.pos", b.cs.pos, 0)
		a.assert("G", "cs.line", b.cs.line, len(b.text)-1)
		e.emit("gg")
		a.assert("gg", "cs.pos", b.cs.pos, 0)
		a.assert("gg", "cs.line", b.cs.line, 0)
	}
	if a.failed {
		for _, m := range a.errMsgs {
			t.Error(m)
		}
	}
}
