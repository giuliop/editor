package main

import "unicode/utf8"

func (m mark) initLastInsert() {
	m.buf.lastInsert = insertText{
		newText: text{line{}},
		oldText: text{line{}},
		start:   &m,
	}
}

func (t *text) appendChar(ch rune) {
	if t.lastChar() == '\n' {
		*t = append(*t, line{})
	}
	(*t)[len(*t)-1] = append((*t)[len(*t)-1], ch)
}

func (t *text) appendChars(cs line) {
	if t.lastChar() == '\n' {
		*t = append(*t, line{})
	}
	(*t)[len(*t)-1] = append((*t)[len(*t)-1], cs...)
}

func (t *text) prependChar(ch rune) {
	if ch == '\n' {
		*t = append(text{line{}}, *t...)
	}
	(*t)[0] = append(line{ch}, (*t)[0]...)
}

func (t text) empty() bool {
	return len(t) == 0 ||
		(len(t) == 1 && len(t[0]) == 0)
}

// lastChar returns the last rune in text or 0 if the text is empty (that is it has
// no lines or one empty line
func (t text) lastChar() rune {
	if t.empty() {
		return 0
	}
	return t[len(t)-1][len(t[len(t)-1])-1]
}

// lastChar returns the last rune in line (typically a newline char)
// or 0 if the line is empty
func (l line) lastChar() rune {
	if len(l) == 0 {
		return 0
	}
	return l[len(l)-1]
}

// lastLIne returns the last line of text t or nil if text is empty
func (t text) lastLine() line {
	if t.empty() {
		return nil
	}
	return t[len(t)-1]
}

func (l line) toBytes() (b []byte) {
	temp := make([]byte, utf8.UTFMax)
	for _, r := range l {
		i := utf8.EncodeRune(temp, r)
		b = append(b, temp[:i]...)
	}
	return b
}

func bytestoLine(b []byte) (l line) {
	for len(b) > 0 {
		r, size := utf8.DecodeRune(b)
		l = append(l, r)
		b = b[size:]
	}
	return l
}

func stringToLine(s string) (l line) {
	return bytestoLine([]byte(s))
}
