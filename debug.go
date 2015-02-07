package main

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"strings"
)

type debugLogger struct {
	*log.Logger
	logfile *os.File
}

func initDebug() *debugLogger {
	f, err := os.OpenFile("log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	check(err)
	logPrefix := "(debug) "
	logFlags := log.Ldate + log.Ltime + log.Lshortfile
	return &debugLogger{log.New(f, logPrefix, logFlags), f}
}

func (d *debugLogger) printStack() {
	b := make([]byte, 1024)
	runtime.Stack(b, false)
	d.Printf("%s", b)
}

func (d *debugLogger) stop() {
	d.logfile.Close()
}

// unsafePrintChannel reads a channel, logs the content, and put it back leaving
// the channel as it was. Not safe for concurrent access
func (d *debugLogger) unsafePrintChannel(c chan UIEvent) {
	c2 := make(chan UIEvent, 1000)
	s := "["
loop:
	for {
		select {
		case x := <-c:
			s += fmt.Sprintf(" %v ", x)
			c2 <- x
		default:
			s += "]"
			break loop
		}
	}
	debug.Println(s)
	for {
		select {
		case x := <-c2:
			c <- x
		default:
			return
		}
	}
}

func keypressesToEmitString(ks []Keypress) string {
	tokens := []string{}
	for _, k := range ks {
		tokens = append(tokens, keyToString(k))
	}
	return fmt.Sprintf("(\"%v\")", strings.Join(tokens, "\", \""))
}

func keyToString(k Keypress) string {
	if k.isSpecial {
		switch k.Special {
		case 0x03:
			return "KeyCtrlC"
		default:
			fatalError(fmt.Errorf("Unknown key %v", k.isSpecial))
			return ""
		}
	} else {
		return string(k.Char)
	}
}
