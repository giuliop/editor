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

var motions = map[string]regionFunc{
	"e": toWordEnd,
	"E": toWORDEnd,
	//"w":  toNextWordStart,
	//"W":  toNextWORDStart,
	//"b":  toWordStart,
	//"B":  toWORDStart,
	//"0":  toFirstCharInLine,
	//"gh": toFirstCharInLine,
	//"$":  toLastCharInLine,
	//"gl": toLastCharInLine,
}

var regionFuncs = map[string]regionFunc{
//"iw": innerword,
//"aw": aword,
}

// we add all motiond to RegionFuncs since all motions are regionFuncs but not
// vicecersa
func init() {
	for k, f := range motions {
		regionFuncs[k] = f
	}
}

var specialWordChars = map[rune]bool{
	'_': true,
}

func isWordChar(c rune) bool {
	return unicode.IsLetter(c) || unicode.IsNumber(c) || specialWordChars[c]
}

func toWordEnd(m mark) region {
	m2 := m
	start := isWordChar(m2.char())
	for {
		m2.moveRight(1)
		c := m2.char()
		end := isWordChar(c)
		switch {
		case m2.atLastTextChar():
			return region{m, m2}
		case m2.pos == m2.lastCharPos():
			return region{m, m2}
		case start != end || unicode.IsSpace(c):
			// we go back to the end of the word, unless we would go back to the
			// initial position
			if m2.line > m.line || m2.pos > m.pos+1 {
				m2.moveLeft(1)
				if !(m2.pos == m.pos && m2.line == m.line) {
					return region{m, m2}
				}
				m2.moveRight(1)
			}
			if unicode.IsSpace(c) {
				m2.moveRight(1)
			}
			start = isWordChar(m2.char())
		}
	}
}

func toWORDEnd(m mark) region {
	m2 := m
	for {
		m2.moveRight(1)
		switch {
		case m2.atLastTextChar():
			return region{m, m2}
		case m2.pos == m2.lastCharPos():
			return region{m, m2}
		case unicode.IsSpace(m2.char()):
			// we go back to the end of the word, unless we would go back to the
			// initial position
			if m2.line > m.line || m2.pos > m.pos+1 {
				m2.moveLeft(1)
				return region{m, m2}
			}
		}
	}
}
