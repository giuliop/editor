package main

import (
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

const TESTFILENAME = "__testFile__"

var (
	ui       = &testUI{}
	keys     = make(chan UIEvent, 100)
	commands = make(chan cmdContext, 100)
)

type testUI struct {
	curBuf *buffer
}

func (u *testUI) Init(b *buffer) error   { return nil }
func (u *testUI) Close()                 {}
func (u *testUI) Draw()                  {}
func (u *testUI) PollEvent() UIEvent     { return UIEvent{} }
func (u *testUI) CurrentBuffer() *buffer { return nil }
func (u *testUI) userMessage(s string)   {}

var (
	defaultText = "" +
		"ciao bello, come va?\n" +
		"tutto bene grazie e tu?\n" +
		"non c'e' male, davvero\n"

	emptyText = "\n"

	emptyLinesText = "" +
		"\n" +
		"\n" +
		"\n"
)

func TestMain(m *testing.M) {
	be = initBackend()
	debug = initDebug()
	defer debug.stop()
	debug.Println("New test run\n")

	stringToFile(defaultText, TESTFILENAME)
	ui.curBuf = be.open([]string{TESTFILENAME})

	go manageKeypress(ui, keys, commands)
	go executeCommands(ui, commands)

	defer func() {}()
	os.Exit(m.Run())
}

type keypressEmitter struct {
	c chan UIEvent
	b *buffer
}

func newKeyPressEmitter(b *buffer) *keypressEmitter {
	return &keypressEmitter{c: keys, b: b}
}

func (e keypressEmitter) emit(a ...interface{}) {
	for _, x := range a {
		switch x.(type) {
		case string:
			stringToEvents(e.b, x.(string))
		case Key:
			keyToEvents(e.b, x.(Key))
		default:
			panic("Unrecognized keypress type")
		}
		// yield to let commands run before new events
		time.Sleep(1 * time.Millisecond)
	}
}

func stringToEvents(b *buffer, s string) {
	for _, c := range s {
		ev := UIEvent{
			Buf:  b,
			Type: UIEventKey,
			Key:  Keypress{Char: c},
		}
		keys <- ev
	}
}

func keyToEvents(b *buffer, k Key) {
	ev := UIEvent{
		Buf:  b,
		Type: UIEventKey,
		Key:  Keypress{Special: k, isSpecial: true},
	}
	keys <- ev
}

type asserter struct {
	failed  bool
	errMsgs []string
}

func (a *asserter) assert(title, name string, actual, expected interface{}) {
	if actual != expected {
		a.failed = true
		a.errMsgs = append(a.errMsgs, fmt.Sprintf(
			"%v - expected %v = %v, got %v", title, name, expected, actual))
	}
	return
}

func stringToFile(text, filename string) {
	b := stringToBuffer(text)
	b.filename = filename
	b.save()
}

func stringToLines(s string) []line {
	lines := strings.Split(s, "\n")
	// last item is empty string if s ends with newline
	if lines[len(lines)-1] == "" {
		lines = lines[:len(lines)-1]
	}
	t := make([]line, len(lines))
	for i, l := range lines {
		t[i] = line(l + "\n")
	}
	return t
}

func stringToBuffer(s string) *buffer {
	b := be.newBuffer("")
	b.text = stringToLines(s)
	return b
}

func bufferToString(b *buffer) string {
	s := ""
	for _, line := range be.text(b) {
		s += string(line)
	}
	return s
}

func TestStringToBufferToString(t *testing.T) {
	if defaultText != bufferToString(stringToBuffer(defaultText)) {
		t.Fail()
	}
}

func TestStringToFileToBufferToString(t *testing.T) {
	s := defaultText
	stringToFile(s, TESTFILENAME)
	b, err := be.openFile(TESTFILENAME)
	if err != nil {
		t.Fatal(err)
	}
	if s != bufferToString(b) {
		t.Fail()
	}
}
