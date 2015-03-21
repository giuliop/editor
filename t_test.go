package main

import (
	"fmt"
	"os"
	"strings"
	"testing"
)

const testFileName = "__testFile__"

func init() {
	cmdKeyInsertMode[testEndOfEmission] = command{allDoneCmd, nil}
	cmdKeyNormalMode[testEndOfEmission] = command{allDoneCmd, nil}
}

func allDoneCmd(ctx *cmdContext) {
	testChan <- struct{}{}
}

var (
	keys = make(chan UIEvent, 100)
)

type testUI struct {
	curBuf *buffer
}

func (u *testUI) Init(b *buffer) error   { return nil }
func (u *testUI) Close()                 {}
func (u *testUI) CurrentView() *view     { return nil }
func (u *testUI) Draw()                  {}
func (u *testUI) PollEvent() UIEvent     { return UIEvent{} }
func (u *testUI) CurrentBuffer() *buffer { return nil }
func (u *testUI) UserMessage(s string)   {}
func (u *testUI) SplitHorizontal()       {}
func (u *testUI) SplitVertical()         {}
func (u *testUI) ToPane(dir direction)   {}

func TestMain(m *testing.M) {
	debug.Println("\nNew test run\n")

	stringToFile(defaultText, testFileName)
	ui = &testUI{}
	ui.Init(be.open([]string{testFileName}))

	go manageKeypress(keys, commands)
	go executeCommands(commands)

	code := func() (code int) {
		defer cleanup()
		code = m.Run()
		os.Remove(testFileName)
		return code
	}
	os.Exit(code())
}

func cleanup() {
	debug.stop()
}

type keypressEmitter struct {
	c chan UIEvent
	v *view
}

func newKeyPressEmitter(v *view) *keypressEmitter {
	return &keypressEmitter{c: keys, v: v}
}

func (e keypressEmitter) emit(a ...interface{}) {
	for _, x := range a {
		switch x.(type) {
		case string:
			stringToEvents(e.v, x.(string))
		case Key:
			keyToEvents(e.v, x.(Key))
		default:
			debug.Println(x)
			panic("Unrecognized keypress type")
		}
	}
	keyToEvents(e.v, testEndOfEmission)
	<-testChan
}

func stringToEvents(v *view, s string) {
	for _, c := range s {
		ev := UIEvent{
			View: v,
			Type: UIEventKey,
			Key:  Keypress{Char: c},
		}
		keys <- ev
	}
}

func keyToEvents(v *view, k Key) {
	ev := UIEvent{
		View: v,
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
	v := stringToView(text)
	v.buf.filename = filename
	v.buf.save()
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

func stringToView(s string) *view {
	b := be.newBuffer("")
	b.text = stringToLines(s)
	b.mod = normalMode
	return &view{b, &mark{0, 0, b}, 0}
}

func viewToString(v *view) string {
	s := ""
	for _, line := range v.buf.text {
		s += string(line)
	}
	return s
}

func TestStringToBufferToString(t *testing.T) {
	for _, s := range samples {
		if s != viewToString(stringToView(s)) {
			t.Fail()
		}
	}
}

func TestStringToFileToBufferToString(t *testing.T) {
	s := defaultText
	stringToFile(s, testFileName)
	b := be.newBuffer("")
	err := be.openFile(b, testFileName)
	if err != nil {
		t.Fatal(err)
	}
	if s != viewToString(&view{b, &mark{0, 0, b}, 0}) {
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
