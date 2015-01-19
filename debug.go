package main

import (
	"os"
)

// log writes msg to file log
func log(msg string) {
	f, err := os.OpenFile("log", os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0666)
	defer f.Close()
	check(err)
	_, err = f.WriteString(msg + "\n")
	check(err)
}
