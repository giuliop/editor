package main

import (
	"log"
	"os"
)

var debug *log.Logger

func init() {
	f, err := os.OpenFile("log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	check(err)
	defer f.Close()
	logPrefix := " * "
	logFlags := log.Ldate + log.Ltime + log.Lshortfile
	debug = log.New(f, logPrefix, logFlags)
}
