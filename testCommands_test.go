package main

import (
	"fmt"
	"strconv"
	"testing"
	"unicode/utf8"
)

type sample struct {
	buf  *buffer
	text string
}

var (
	defaultText = "" +
		"   xxx_yyy xxx___yyy xxx_^_ppp  \n" +
		"func (e keypressEmitter) emit(a ...interface{}) {\n" +
		"\n" +
		"   xxx***(((_ciao *** &&& ff.ff  *\n" +
		"*\n" +
		" _ \n" +
		"non c'e' male, davvero .... \n"

	emptyText = "\n"

	emptyLinesText = "" +
		"\n" +
		"\n" +
		"\n"
)

func TestLineMotions(t *testing.T) {
	var samples = []string{
		defaultText,
		emptyText,
		emptyLinesText,
	}
	a := &asserter{}
	for _, s := range samples {
		// test 'gg' and 'G'
		b := stringToBuffer(s)
		e := newKeyPressEmitter(b)
		e.emit(KeyCtrlC, "G")
		a.assert("G", "cs.pos", b.cs.pos, 0)
		a.assert("G", "cs.line", b.cs.line, len(b.text)-1)
		e.emit("gg")
		a.assert("gg", "cs.pos", b.cs.pos, 0)
		a.assert("gg", "cs.line", b.cs.line, 0)
		// test '$', '0', 'L', 'H',
		e.emit("$")
		exp := len(b.text[0]) - 2
		if exp < 0 {
			exp = 0
		}
		a.assert("$", "cs.pos", b.cs.pos, exp)
		e.emit("0")
		a.assert("0", "cs.pos", b.cs.pos, 0)
		e.emit("L")
		a.assert("L", "cs.pos", b.cs.pos, exp)
		e.emit("H")
		a.assert("H", "cs.pos", b.cs.pos, 0)
	}
	if a.failed {
		for _, m := range a.errMsgs {
			t.Error(m)
		}
	}
}

// quickmark is a mark without buffer files
type quickmark struct{ line, pos int }

// m returns a quickmark
func m(line, pos int) *quickmark {
	return &quickmark{line, pos}
}

func TestWordMotions(t *testing.T) {
	var samples = []string{
		defaultText,
		emptyText,
		emptyLinesText,
	}
	testKeys := []string{"e", "b", "w"}
	eMarks := [][]*quickmark{
		[]*quickmark{m(0, 9), m(0, 19), m(0, 24), m(0, 25), m(0, 29), m(1, 3),
			m(1, 5), m(1, 6), m(1, 22), m(1, 23), m(1, 28), m(1, 29), m(1, 30), m(1, 34),
			m(1, 43), m(1, 46), m(1, 48), m(3, 5), m(3, 11), m(3, 16), m(3, 20), m(3, 24),
			m(3, 27), m(3, 28), m(3, 30), m(3, 33), m(4, 0), m(5, 1), m(6, 2), m(6, 4),
			m(6, 5), m(6, 6), m(6, 7), m(6, 12), m(6, 13), m(6, 21), m(6, 26), m(6, 27)},
		[]*quickmark{m(0, 0)},
		[]*quickmark{m(2, 0)},
	}
	bMarks := [][]*quickmark{
		[]*quickmark{m(6, 23), m(6, 15), m(6, 13), m(6, 9), m(6, 7), m(6, 6), m(6, 5),
			m(6, 4), m(6, 0), m(5, 1), m(4, 0), m(3, 33), m(3, 29), m(3, 28),
			m(3, 26), m(3, 22), m(3, 18), m(3, 12), m(3, 6), m(3, 3), m(1, 48),
			m(1, 44), m(1, 35), m(1, 32), m(1, 30), m(1, 29), m(1, 25), m(1, 23),
			m(1, 8), m(1, 6), m(1, 5), m(1, 0), m(0, 26), m(0, 25), m(0, 21),
			m(0, 11), m(0, 3), m(0, 0)},
		[]*quickmark{m(0, 0)},
		[]*quickmark{m(0, 0)},
	}
	wMarks := [][]*quickmark{
		[]*quickmark{m(0, 3), m(0, 11), m(0, 21), m(0, 25), m(0, 26), m(1, 0), m(1, 5),
			m(1, 6), m(1, 8), m(1, 23), m(1, 25), m(1, 29), m(1, 30), m(1, 32), m(1, 35),
			m(1, 44), m(1, 48), m(3, 3), m(3, 6), m(3, 12), m(3, 18), m(3, 22), m(3, 26),
			m(3, 28), m(3, 29), m(3, 33), m(4, 0), m(5, 1), m(6, 0), m(6, 4), m(6, 5),
			m(6, 6), m(6, 7), m(6, 9), m(6, 13), m(6, 15), m(6, 23), m(6, 27)},
		[]*quickmark{m(0, 0)},
		[]*quickmark{m(2, 0)},
	}
	expected := [][][]*quickmark{eMarks, bMarks, wMarks}
	a := _testMotions(samples, testKeys, expected)
	if a.failed {
		for _, m := range a.errMsgs {
			t.Error(m)
		}
	}
}

func TestWORDMotions(t *testing.T) {
	var samples = []string{
		defaultText,
		emptyText,
		emptyLinesText,
	}
	testKeys := []string{"E", "B", "W"}
	eMarks := [][]*quickmark{
		[]*quickmark{m(0, 9), m(0, 19), m(0, 29), m(1, 3), m(1, 6), m(1, 23), m(1, 30),
			m(1, 46), m(1, 48), m(3, 16), m(3, 20), m(3, 24), m(3, 30), m(3, 33), m(4, 0),
			m(5, 1), m(6, 2), m(6, 7), m(6, 13), m(6, 21), m(6, 26), m(6, 27)},
		[]*quickmark{m(0, 0)},
		[]*quickmark{m(2, 0)},
	}
	bMarks := [][]*quickmark{
		[]*quickmark{m(6, 23), m(6, 15), m(6, 9), m(6, 4), m(6, 0), m(5, 1), m(4, 0),
			m(3, 33), m(3, 26), m(3, 22), m(3, 18), m(3, 3), m(1, 48), m(1, 32), m(1, 25),
			m(1, 8), m(1, 5), m(1, 0), m(0, 21), m(0, 11), m(0, 3), m(0, 0)},
		[]*quickmark{m(0, 0)},
		[]*quickmark{m(0, 0)},
	}
	wMarks := [][]*quickmark{
		[]*quickmark{m(0, 3), m(0, 11), m(0, 21), m(1, 0), m(1, 5), m(1, 8), m(1, 25),
			m(1, 32), m(1, 48), m(3, 3), m(3, 18), m(3, 22), m(3, 26), m(3, 33), m(4, 0),
			m(5, 1), m(6, 0), m(6, 4), m(6, 9), m(6, 15), m(6, 23), m(6, 27)},
		[]*quickmark{m(0, 0)},
		[]*quickmark{m(2, 0)},
	}
	expected := [][][]*quickmark{eMarks, bMarks, wMarks}
	a := _testMotions(samples, testKeys, expected)
	if a.failed {
		for _, m := range a.errMsgs {
			t.Error(m)
		}
	}
}

func _testMotions(samples []string, testKeys []string, expected [][][]*quickmark) *asserter {
	a := &asserter{}
	for _s, s := range samples {
		b := stringToBuffer(s)
		e := newKeyPressEmitter(b)
		e.emit(KeyCtrlC)
		for _k, k := range testKeys {
			for i, m := range expected[_k][_s] {
				e.emit(k)
				head := k + " - " + strconv.Itoa(i)
				a.assert(head, "cs.pos", b.cs.pos, m.pos)
				a.assert(head, "cs.line", b.cs.line, m.line)
				if b.cs.atLastTextChar() || b.cs.atStartOfText() {
					a.assert(k, "all marks consummed", i == len(expected[_k][_s])-1, true)
				}
			}
		}
	}
	return a
}

func TestInsert(t *testing.T) {
	b := stringToBuffer("123")
	e := newKeyPressEmitter(b)
	e.emit(KeyCtrlC, "li0")
	err := equalStrings(bufferToString(b), "1023\n")
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestAppend(t *testing.T) {
	b := stringToBuffer("123")
	e := newKeyPressEmitter(b)
	e.emit(KeyCtrlC, "la0")
	err := equalStrings(bufferToString(b), "1203\n")
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestAppendEndOfLine(t *testing.T) {
	b := stringToBuffer("123")
	e := newKeyPressEmitter(b)
	e.emit(KeyCtrlC, "lA0")
	err := equalStrings(bufferToString(b), "1230\n")
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestAppendEndOfLineInsertMode(t *testing.T) {
	b := stringToBuffer("123")
	e := newKeyPressEmitter(b)
	e.emit(KeyCtrlC, "liAA0")
	err := equalStrings(bufferToString(b), "1230\n")
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func TestMultiInsertAppend(t *testing.T) {
	b := stringToBuffer(defaultText)
	e := newKeyPressEmitter(b)
	e.emit(KeyCtrlC, "l", "l", "A", "g", "g", "g", KeyCtrlC, "j", "h", "h", "h", "h",
		"h", "l", "i", "c", "A", "A", "v", KeyCtrlC, "j", "i", "c", KeyCtrlC, "j",
		"j", "a", "d")
	expected := "" +
		"   xxx_yyy xxx___yyy xxx_^_ppp  ggg\n" +
		"func (e keypressEmitter) emit(ca ...interface{}) {v\n" +
		"c\n" +
		"   xxx***(((_ciao *** &&& ff.ff  *\n" +
		"*d\n" +
		" _ \n" +
		"non c'e' male, davvero .... \n"
	err := equalStrings(bufferToString(b), expected)
	if err != nil {
		t.Fatalf(err.Error())
	}
}

func equalStrings(actual, expected string) error {
	count, line, offset := 0, 0, 0
	for i, c := range actual {
		if count >= len(expected) {
			return fmt.Errorf("expected shorter than actual!")
		}
		r, size := utf8.DecodeRuneInString(expected[count:])
		count += size
		if c != r {
			debug.Printf("actual\n%v\nexpected\n%v\n", actual, expected)
			return fmt.Errorf("no match at line %v pos %v, expected %q, found %q",
				line+1, i-offset, r, c)
		}
		if r == '\n' {
			line++
			offset = i
		}
	}
	if expected[count:] != "" {
		return fmt.Errorf("actual shorter than expected!")
	}
	return nil
}

func XXXTestDeleteToEndAndStartOfLine(t *testing.T) {
	b := stringToBuffer(defaultText)
	e := newKeyPressEmitter(b)
	e.emit()
	expected := "" +
		"   xxx_yyy xxx___yyy xxx_^_ppp  ggg\n" +
		"func (e keypressEmitter) emit(ca ...interface{}) {v\n" +
		"c\n" +
		"   xxx***(((_ciao *** &&& ff.ff  *\n" +
		"*d\n" +
		" _ \n" +
		"non c'e' male, davvero .... \n"
	err := equalStrings(bufferToString(b), expected)
	if err != nil {
		t.Fatalf(err.Error())
	}
}
