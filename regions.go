package main

import "unicode"

type region struct {
	start mark
	end   mark
}

// a regionFunc returns a region, that is a start mark and an end mark
// while a text object is obviously a region, motions are as well (with the
// current cursor position as start mark)
type regionFunc func(m mark) region

var regionFuncs = map[string]regionFunc{
	"e": toWordEnd,
	//"E":  toWORDEnd,
	//"w":  toNextWordStart,
	//"W":  toNextWORDStart,
	//"b":  toWordStart,
	//"B":  toWORDStart,
	//"0":  toFirstCharInLine,
	//"gh": toFirstCharInLine,
	//"$":  toLastCharInLine,
	//"gl": toLastCharInLine,
	//"iw": innerword,
	//"aw": aword,
}

func toWordEnd(m mark) region {
	m2 := m
	for {
		m2.moveRight(1)
		if m2.atLastLine() && m2.pos == m2.lastCharPos() {
			return region{m, m2}
		}
		c := m2.char()
		if !(unicode.IsLetter(c) || unicode.IsNumber(c)) {
			m2.moveLeft(1)
			return region{m, m2}
		}
	}
}
