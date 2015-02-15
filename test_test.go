package main

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

const TESTFILENAME = "__testFile__"
const endOfEmission = KeyCtrlBackslash

func init() {
	cmdKeyInsertMode[endOfEmission] = command{allDoneCmd, nil}
	cmdKeyNormalMode[endOfEmission] = command{allDoneCmd, nil}
}

func allDoneCmd(ctx *cmdContext) {
	allDone <- struct{}{}
}

var (
	keys     = make(chan UIEvent, 100)
	commands = make(chan cmdContext, 100)
	allDone  = make(chan struct{}) // used this to signal all commands done
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

func TestMain(m *testing.M) {
	debug.Println("\nNew test run\n")

	stringToFile(defaultText, TESTFILENAME)
	ui = &testUI{}
	ui.Init(be.open([]string{TESTFILENAME}))

	go manageKeypress(keys, commands)
	go executeTestCommands(commands)

	defer cleanup()
	os.Exit(m.Run())
}

func cleanup() {
	debug.stop()
}

func executeTestCommands(cmds chan cmdContext) {
	defer cleanupOnError()
	for {
		ctx := <-cmds
		ctx.cmd(&ctx)
		ctx.cmdChans.done <- cmdDone
	}
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
			debug.Println(x)
			panic("Unrecognized keypress type")
		}
	}
	keyToEvents(e.b, endOfEmission)
	<-allDone
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
	b.mod = normalMode
	return b
}

func bufferToString(b *buffer) string {
	s := ""
	for _, line := range b.text {
		s += string(line)
	}
	return s
}

func TestStringToBufferToString(t *testing.T) {
	for _, s := range samples {
		if s != bufferToString(stringToBuffer(s)) {
			t.Fail()
		}
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

func recordTestMacro(ctx *cmdContext) {
	if r.macros.on {
		// save the macro keys removing the last key which is end record key
		keys := r.macros.keys[:len(r.macros.keys)-1]
		r.macros.macros[0] = keys
		r.macros.stop()
		ctx.msg = "finished recording"
		return
	}
	r.macros.start()
	ctx.msg = "started macro recording"
}
