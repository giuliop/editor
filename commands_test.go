package main

import (
	"strings"
	"testing"
)

var _be = initBackend()
var text1 = "" +
	"ciao bello, come va?\n" +
	"tutto bene grazie e tu?\n" +
	"non c'e' male, davvero\n"

func newTestBuffer(name, text string) *buffer {
	b := _be.newBuffer(name)
	lines := strings.Split(text, "\n")
	t := make([]line, len(lines))
	for i, l := range lines {
		t[i] = line(l + "\n")
	}
	b.text = t
	return b
}

//func TestAppendAtEndOfLine(*testing.T) {
//b := newBuffer("1", text1)
//m := b.cs
//ctx := &cmdContext{point: &m}
//appendAtEndOfLine(ctx)
//}
