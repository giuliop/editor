package main

import (
	"regexp"
	"unicode"
)

type region struct {
	start mark
	end   mark
}

// a regionFunc returns a region, that is a start mark and an end mark
// while a text object is obviously a region, motions are as well (with the
// current cursor position as start mark)
type regionFunc func(m mark) (region, direction)

var motions = map[string]regionFunc{
	"e":  toWordEnd,
	"E":  toWORDEnd,
	"w":  toNextWordStart,
	"W":  toNextWORDStart,
	"b":  toWordStart,
	"B":  toWORDStart,
	"$":  toLineEnd,
	"L":  toLineEnd,
	"0":  toLineStart,
	"H":  toLineStart,
	"gg": toFirstLine,
	"G":  toLastLine,
}

var regionFuncs = map[string]regionFunc{
//"iw": innerword,
//"aw": aword,
}

// we add all motions to RegionFuncs since all motions are regionFuncs but not
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

func isSymbolWordChar(c rune) bool {
	return !(unicode.IsSpace(c) || c == '\n' || isWordChar(c))
}

func doMotion(m *mark, atMotion func() bool, moveMark func(n int)) region {
	m2 := *m
	moveMark(1)
	for ; !(atMotion() || m.atLastTextChar() || m.atStartOfText()); moveMark(1) {
	}
	return region{m2, *m}
}

func toLineEnd(m mark) (region, direction) {
	m2 := mark{m.line, m.maxCursPos(), m.buf}
	return region{m, m2}, right
}

func toLineStart(m mark) (region, direction) {
	m2 := mark{m.line, 0, m.buf}
	return region{m, m2}, left
}

func toLastLine(m mark) (region, direction) {
	m2 := mark{m.lastLine(), 0, m.buf}
	return region{m, m2}, right
}

func toFirstLine(m mark) (region, direction) {
	m2 := mark{0, 0, m.buf}
	return region{m, m2}, left
}

func toWORDEnd(m mark) (region, direction) {
	return doMotion(&m, m.atEndOfWORD, m.moveRight), right
}

func toWordEnd(m mark) (region, direction) {
	return doMotion(&m, m.atEndOfWord, m.moveRight), right
}

func toNextWORDStart(m mark) (region, direction) {
	return doMotion(&m, m.atStartOfWORD, m.moveRight), right
}

func toNextWordStart(m mark) (region, direction) {
	return doMotion(&m, m.atStartOfWord, m.moveRight), right
}
func toWORDStart(m mark) (region, direction) {
	return doMotion(&m, m.atStartOfWORD, m.moveLeft), left
}

func toWordStart(m mark) (region, direction) {
	return doMotion(&m, m.atStartOfWord, m.moveLeft), left
}

// a non-space followed by either a space or newline
func (m *mark) atEndOfWORD() bool {
	c, n := m.char(), m.nextChar()
	return m.atLastTextChar() ||
		(!unicode.IsSpace(c) && (unicode.IsSpace(n) || n == '\n'))
}

// a non-space preceded by either a space or newline
func (m *mark) atStartOfWORD() bool {
	c, p := m.char(), m.prevChar()
	return m.atStartOfText() ||
		(!unicode.IsSpace(c) && (unicode.IsSpace(p) || p == 0))
}

// (symbol) word char followed by a non (symbol) word char
func (m *mark) atEndOfWord() bool {
	c, n := m.char(), m.nextChar()
	return m.atLastTextChar() ||
		(isWordChar(c) && !isWordChar(n)) ||
		(isSymbolWordChar(c) && !isSymbolWordChar(n))
}

// (symbol) word char precede by a non (symbol) word char
func (m *mark) atStartOfWord() bool {
	c, p := m.char(), m.prevChar()
	return m.atStartOfText() ||
		(isWordChar(c) && !isWordChar(p)) ||
		(isSymbolWordChar(c) && (!isSymbolWordChar(p) || unicode.IsSpace(p))) ||
		(!unicode.IsSpace(c) && (p == 0))
}

func findRight(m mark, r *regexp.Regexp) mark {
	text := m.buf.text
	offset := m.pos + 1
	for ln := m.line; ln <= m.lastLine(); ln++ {
		s := string(text[ln][offset:])
		pos := r.FindStringIndex(s)
		if pos != nil {
			return mark{ln, pos[0] + offset, m.buf}
		}
		offset = 0
	}
	return m.lastTextCharPos()
}

func findLeft(m mark, r *regexp.Regexp) mark {
	text := m.buf.text
	for ln := m.line; ln >= 0; ln-- {
		s := string(text[ln])
		if ln == m.line {
			s = s[:m.pos]
		}
		matches := r.FindAllStringSubmatchIndex(s, -1)
		if matches != nil {
			pos := matches[len(matches)-1]
			return mark{ln, pos[2], m.buf}
		}
	}
	return m.firstTextCharPos()
}
