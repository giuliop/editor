package main

import (
	"strconv"
	"unicode"
)

// isNumber takes a key and a context and returns whether the key should be
// treated as a number; returns true if key is a digit but false if key is '0'
// and there is no number in ctx.num (that is the 0 is not there to complete
// a number like 10 or 02)

func isNumber(ch rune, ctx *cmdContext) bool {
	if !unicode.IsDigit(ch) || ctx.point.buf.mod != normalMode {
		return false
	}
	if ch == '0' && ctx.num == 0 {
		return false
	}
	return true
}

func loadNumber(key rune, ctx *cmdContext) error {
	num, err := strconv.Atoi(strconv.Itoa(ctx.num) + string(key))
	if err == nil {
		ctx.num = num
	}
	return err
}
