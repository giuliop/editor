package main

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
