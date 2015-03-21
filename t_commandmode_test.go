package main

import "testing"

func _TestPrevNextCmd(t *testing.T) {
	a := &asserter{}
	v := stringToView("\n")
	e := newKeyPressEmitter(v)
	e.emit(";1", KeyEnter, ";2", KeyEnter, ";3", KeyEnter)
	e.emit(";", KeyArrowUp)
	a.assert("1", "command", string(be.msgLine[len(commandModePrompt):]), "3")
	e.emit(KeyArrowUp)
	a.assert("2", "command", string(be.msgLine[len(commandModePrompt):]), "2")
	e.emit(KeyArrowUp)
	a.assert("3", "command", string(be.msgLine[len(commandModePrompt):]), "1")
	e.emit(KeyArrowUp)
	a.assert("4", "command", string(be.msgLine[len(commandModePrompt):]), "1")
	e.emit(KeyArrowUp)
	a.assert("5", "command", string(be.msgLine[len(commandModePrompt):]), "1")
	e.emit(KeyArrowDown)
	a.assert("6", "command", string(be.msgLine[len(commandModePrompt):]), "2")
	e.emit(KeyArrowDown)
	a.assert("7", "command", string(be.msgLine[len(commandModePrompt):]), "3")
	e.emit(KeyArrowDown)
	a.assert("8", "command", string(be.msgLine[len(commandModePrompt):]), "")
	e.emit(KeyArrowDown)
	a.assert("", "command", string(be.msgLine[len(commandModePrompt):]), "")
	e.emit(KeyEsc)
	if a.failed {
		for _, m := range a.errMsgs {
			t.Error(m)
		}
	}
}

func _TestPrevNextCmdSubMatches(t *testing.T) {
	r.commands.reset()
	a := &asserter{}
	v := stringToView("\n")
	e := newKeyPressEmitter(v)
	e.emit(";aabbcc", KeyEnter, ";aacc", KeyEnter, ";aabbdd", KeyEnter,
		";bbaa", KeyEnter)
	e.emit(";aa", KeyArrowUp)
	a.assert("1", "command", string(be.msgLine[len(commandModePrompt):]), "aabbdd")
	e.emit(KeyArrowUp)
	a.assert("2", "command", string(be.msgLine[len(commandModePrompt):]), "aacc")
	e.emit(KeyArrowUp)
	a.assert("3", "command", string(be.msgLine[len(commandModePrompt):]), "aabbcc")
	e.emit(KeyArrowUp)
	a.assert("4", "command", string(be.msgLine[len(commandModePrompt):]), "aabbcc")
	e.emit(KeyBackspace, KeyBackspace)
	e.emit(KeyArrowDown)
	a.assert("5", "command", string(be.msgLine[len(commandModePrompt):]), "aabbdd")
	e.emit(KeyArrowDown)
	a.assert("6", "command", string(be.msgLine[len(commandModePrompt):]), "aabb")
	e.emit(KeyEsc)
	if a.failed {
		for _, m := range a.errMsgs {
			t.Error(m)
		}
	}
}
