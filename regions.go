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
type regionFunc func(m mark) region

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

/*
 *var (
 *    WORDEnd = regexp.MustCompile(`[\P{Zs}][\p{Zs}\n]`)
 *    // a non-space, non-newline preceded by either a space or start of text
 *    WORDStart       = regexp.MustCompile(`(?:\A|\p{Zs})([^\p{Zs}\n])`)
 *    wordEnd         *regexp.Regexp
 *    symbolWordEnd   *regexp.Regexp
 *    wordStart       *regexp.Regexp
 *    symbolWordStart *regexp.Regexp
 *)
 *s := ``
 *for r := range specialWordChars {
 *    s += string(r)
 *}
 *wordEnd = regexp.MustCompile(`[\pL\d` + s + `][^\pL\d` + s + `]`)
 *symbolWordEnd = regexp.MustCompile(`[^\pL\d\n\p{Zs}` + s + `][\pL\d\n\p{Zs}` + s + `]`)
 *wordStart = regexp.MustCompile(`(?:\A|[^\pL\d` + s + `])([\pL\d` + s + `])`)
 *symbolWordStart = regexp.MustCompile(`(?:\A|[\pL\d\p{Zs}` + s + `])([^\pL\d\p{Zs}` + s + `])`)
 */

// we add all motiond to RegionFuncs since all motions are regionFuncs but not

// vicecersa; we also compile motion regexp
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

func toLineEnd(m mark) region {
	m2 := mark{m.line, m.maxCursPos(), m.buf}
	return region{m, m2}
}

func toLineStart(m mark) region {
	m2 := mark{m.line, 0, m.buf}
	return region{m, m2}
}

func toLastLine(m mark) region {
	m2 := mark{m.lastLine(), 0, m.buf}
	return region{m, m2}
}

func toFirstLine(m mark) region {
	m2 := mark{0, 0, m.buf}
	return region{m, m2}
}

func toWORDEnd(m mark) region {
	return doMotion(&m, m.atEndOfWORD, m.moveRight)
}

func toWordEnd(m mark) region {
	return doMotion(&m, m.atEndOfWord, m.moveRight)
}

func toNextWORDStart(m mark) region {
	return doMotion(&m, m.atStartOfWORD, m.moveRight)
}

func toNextWordStart(m mark) region {
	return doMotion(&m, m.atStartOfWord, m.moveRight)
}
func toWORDStart(m mark) region {
	return doMotion(&m, m.atStartOfWORD, m.moveLeft)
}

func toWordStart(m mark) region {
	return doMotion(&m, m.atStartOfWord, m.moveLeft)
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
	return be.lastTextCharPos(m)
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
	return be.firstTextCharPos(m)
}
