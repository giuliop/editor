package main

// UI is the controller for a terminal or gui frontend
type UI interface {
	Init(b *buffer) error
	Close()
	Draw()
	//Clear() error
	//Flush() error
	//SetCursor(x, y int)
	//HideCursor()
	//SetCell(x, y int, ch rune)
	PollEvent() UIEvent
	CurrentBuffer() *buffer
	userMessage(s string)
}

type UIEvent struct {
	Buf    *buffer
	Type   UIEventType // one of Event* constants
	Mod    UIModifier  // one of Mod* constants or 0
	Key    Keypress    // a keypress
	Width  int         // width of the screen
	Height int         // height of the screen
	Err    error       // error in case if input failed
	MouseX int         // x coord of mouse
	MouseY int         // y coord of mouse
}

type UIEventType uint8

const (
	UIEventKey UIEventType = iota
	UIEventResize
	UIEventMouse
	UIEventError
	UIEventTimeout
)

type UIModifier uint8

const (
	ModAlt UIModifier = 0x01
)

type Keypress struct {
	Special   Key  // a special key or CTRL combination
	Char      rune // a normal char key
	isSpecial bool // true is a special key was pressed
}

type Key uint16

const (
	KeyF1 Key = 0xFFFF - iota
	KeyF2
	KeyF3
	KeyF4
	KeyF5
	KeyF6
	KeyF7
	KeyF8
	KeyF9
	KeyF10
	KeyF11
	KeyF12
	KeyInsert
	KeyDelete
	KeyHome
	KeyEnd
	KeyPgup
	KeyPgdn
	KeyArrowUp
	KeyArrowDown
	KeyArrowLeft
	KeyArrowRight

	MouseLeft
	MouseMiddle
	MouseRight
)

const (
	KeyCtrlTilde      Key = 0x00
	KeyCtrl2          Key = 0x00
	KeyCtrlSpace      Key = 0x00
	KeyCtrlA          Key = 0x01
	KeyCtrlB          Key = 0x02
	KeyCtrlC          Key = 0x03
	KeyCtrlD          Key = 0x04
	KeyCtrlE          Key = 0x05
	KeyCtrlF          Key = 0x06
	KeyCtrlG          Key = 0x07
	KeyBackspace      Key = 0x08
	KeyCtrlH          Key = 0x08
	KeyTab            Key = 0x09
	KeyCtrlI          Key = 0x09
	KeyCtrlJ          Key = 0x0A
	KeyCtrlK          Key = 0x0B
	KeyCtrlL          Key = 0x0C
	KeyEnter          Key = 0x0D
	KeyCtrlM          Key = 0x0D
	KeyCtrlN          Key = 0x0E
	KeyCtrlO          Key = 0x0F
	KeyCtrlP          Key = 0x10
	KeyCtrlQ          Key = 0x11
	KeyCtrlR          Key = 0x12
	KeyCtrlS          Key = 0x13
	KeyCtrlT          Key = 0x14
	KeyCtrlU          Key = 0x15
	KeyCtrlV          Key = 0x16
	KeyCtrlW          Key = 0x17
	KeyCtrlX          Key = 0x18
	KeyCtrlY          Key = 0x19
	KeyCtrlZ          Key = 0x1A
	KeyEsc            Key = 0x1B
	KeyCtrlLsqBracket Key = 0x1B
	KeyCtrl3          Key = 0x1B
	KeyCtrl4          Key = 0x1C
	KeyCtrlBackslash  Key = 0x1C
	KeyCtrl5          Key = 0x1D
	KeyCtrlRsqBracket Key = 0x1D
	KeyCtrl6          Key = 0x1E
	KeyCtrl7          Key = 0x1F
	KeyCtrlSlash      Key = 0x1F
	KeyCtrlUnderscore Key = 0x1F
	KeySpace          Key = 0x20
	KeyBackspace2     Key = 0x7F
	KeyCtrl8          Key = 0x7F
)
