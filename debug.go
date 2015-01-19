package main

import (
	"log"
	"os"
	"runtime"
)

type debugLogger struct {
	*log.Logger
	logfile *os.File
}

var debug *debugLogger

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
	d.Printf(" * Fatal error * \n\n%s\n", b)
}

func (d *debugLogger) stop() {
	d.logfile.Close()
}
