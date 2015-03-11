package main

import "testing"

func TestLineIndent(t *testing.T) {
	s := "" +
		"\tHello\n" +
		"\t\t\tDude\n" +
		"   How\n" +
		"\t\t are\n" +
		"   \t\t  you?\n" +
		"ok\n" +
		" \t \n" +
		"thanks\n"

	res := [][]int{
		[]int{1 * tabStop, 1},
		[]int{3 * tabStop, 3},
		[]int{3, 3},
		[]int{2*tabStop + 1, 3},
		[]int{2*tabStop + 5, 7},
		[]int{0, 0},
		[]int{tabStop + 2, 3},
		[]int{0, 0},
	}
	a := &asserter{}
	v := stringToView(s)
	for i := range v.buf.text {
		indent, indentChars := lineIndent(v.buf, i)
		a.assert(string(i), "", indent, res[i][0])
		a.assert(string(i), "", indentChars, res[i][1])
	}
	if a.failed {
		for _, m := range a.errMsgs {
			t.Error(m)
		}
	}
}

func TestLineToBytes(t *testing.T) {
	s := "func a() {"
	b := []byte(s)
	l := bytestoLine(b)
	a := &asserter{}
	a.assert(s, string(l.toBytes()), string(l.toBytes()) == s, true)
	if a.failed {
		for _, m := range a.errMsgs {
			t.Error(m)
		}
	}
}
