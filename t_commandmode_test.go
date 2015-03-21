package main

import (
	"testing"
	"time"
)

func TestPrevNextCmd(t *testing.T) {
	a := &asserter{}
	v := stringToView("")
	e := newKeyPressEmitter(v)
	waitF := func() {
		for len(be.msgLine) == 0 {
			time.Sleep(1 * time.Millisecond)
		}
	}
	e.emit(";1", KeyEnter, ";2", KeyEnter, ";3", KeyEnter)
	e.emitAsynch(waitF, ";", KeyArrowUp)
	a.assert("", "command", string(be.msgLine[len(commandModePrompt):]), "3")
	e.emit(KeyEnter)
	if a.failed {
		for _, m := range a.errMsgs {
			t.Error(m)
		}
	}
}
