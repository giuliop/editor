package main

type mark struct {
	line int     // the line the mark is at (line starts from 0)
	pos  int     // the position in the line (0 is left of first character)
	buf  *buffer // the buffer the mark is in
}

func newMark(b *buffer) *mark {
	return &mark{0, 0, b}
}

func (m *mark) atFirstLine() bool {
	return m.line == 0
}

func (m *mark) atLastLine() bool {
	return m.line == m.lastLine()
}

func (m *mark) lastLine() int {
	return len(m.buf.text) - 1
}

func (m *mark) atLineStart() bool {
	return m.pos == 0
}

// atLineEnd returns whether the mark is at line end, that is left of newline char
func (m *mark) atLineEnd() bool {
	return m.pos == len(m.buf.text[m.line])-1
}

func (m *mark) atStartOfText() bool {
	return m.atFirstLine() && m.atLineStart()
}

func (m *mark) atEndOfText() bool {
	return m.atLastLine() && m.atLineEnd()
}

// atLastTextChar return whether the mark is left of the last character of the last
// line (or before the newline if the last line is empty)
func (m *mark) atLastTextChar() bool {
	return m.atLastLine() && (m.pos == m.lastCharPos() || m.atLineEnd())
}

// lastCharPos return the position of the last char in the line before the newline
// If the line is empty it returns -1
func (m *mark) lastCharPos() int {
	return len(m.buf.text[m.line]) - 2
}

// maxCursPos returns the maximum position the cursor might be on, which is left of
//the newline char in insert mode or for empty lines and left of the last char for
// non empty lines in normal mode
func (m *mark) maxCursPos() int {
	max := m.lastCharPos()
	if m.buf.mod == insertMode || max < 0 {
		max++
	}
	return max
}

// lineEndPos return the position of the newline char at the end of the line
func (m mark) lineEndPos() int {
	return len(m.buf.text[m.line]) - 1
}

func (m *mark) atEmptyLine() bool {
	return m.lastCharPos() == -1
}

// fixPos checks that the mark is within the line, if it is over the end of the line
// puts the cursor back to the end, if before start of line, back to start of line
func (m *mark) fixPos() {
	max := m.maxCursPos()
	if m.pos > max {
		m.pos = max
		return
	}
	if m.pos < 0 {
		m.pos = 0
	}
}

//fixLine checks that the cursor is on a valid line and puts the cursor
// on either the first or last line otherwise
func (m *mark) fixLine() {
	max := len(m.buf.text) - 1
	if m.line > max {
		m.line = max
		return
	}
	if m.line < 0 {
		m.line = 0
	}
}

func (m *mark) fixLineAndPos() {
	m.fixLine()
	m.fixPos()
}

func (m *mark) moveUp(steps int) {
	m.line -= steps
	if m.line < 0 {
		m.line = 0
	}
	m.fixPos()
}

func (m *mark) moveDown(steps int) {
	m.line += steps
	if m.line > m.lastLine() {
		m.line = m.lastLine()
	}
	m.fixPos()
}

func (m *mark) moveRight(steps int) {
	maxY := m.lastLine()
	for steps > 0 {
		maxX := m.maxCursPos()
		if maxX >= m.pos+steps {
			m.pos += steps
			break
		}
		if m.line < maxY {
			steps -= (maxX - m.pos + 1)
			m.set(m.line+1, 0)
		} else {
			m.pos = maxX
			break
		}
	}
}

func (m *mark) moveLeft(steps int) {
	for steps > 0 {
		if m.pos-steps >= 0 {
			m.pos -= steps
			break
		}
		if m.line > 0 {
			steps -= (m.pos + 1)
			m.line -= 1
			m.pos = m.lastCharPos()
			m.fixPos()
		} else {
			m.pos = 0
			break
		}
	}
}

func (m *mark) set(line, pos int) {
	m.line = line
	m.pos = pos
}

func (m *mark) hide() {
	m.set(-1, -1)
}

func (m *mark) char() rune {
	return m.buf.text[m.line][m.pos]
}

func (m *mark) prevChar() rune {
	if m.atLineStart() {
		return 0
	}
	return m.buf.text[m.line][m.pos-1]
}

func (m *mark) nextChar() rune {
	if m.atLineEnd() {
		return 0
	}
	return m.buf.text[m.line][m.pos+1]
}

// orderMarks takes two marks (assumed in same buffer) and returns the two marks
// in order, that is the mark earlier in the text is returned first
func orderMarks(m1, m2 mark) (mark, mark) {
	switch {
	case m1.line < m2.line:
		return m1, m2
	case m1.line > m2.line:
		return m2, m1
	case m1.pos < m2.pos:
		return m1, m2
	default:
		return m2, m1
	}
}

// deltaChars returns the distance in number of chars (including newlines)
// between m and m2; the value is negative if m2 is before m and 0 if they overlap
// the function assume the marks are on the same buffer
func (m *mark) deltaChars(m2 mark) (delta int) {
	t := m.buf.text
	fr, to := orderMarks(*m, m2)
	if fr.line == to.line {
		delta = len(t[fr.line][fr.pos:to.pos])
	} else {
		delta = len(t[fr.line][fr.pos:]) + len(t[to.line][:to.pos])
		for fr.line++; fr.line < to.line; fr.line++ {
			delta += len(t[fr.line])
		}
	}
	if fr.line == m2.line && fr.pos == m2.pos {
		delta *= -1
	}
	return delta
}

func (m *mark) maxLine() int {
	return len(m.buf.text) - 1
}

func (m mark) lastTextCharPos() mark {
	m2 := mark{m.lastLine(), 0, m.buf}
	m2.pos = m2.lastCharPos()
	m2.fixPos()
	return m2
}

func (m mark) firstTextCharPos() mark {
	m2 := mark{0, 0, m.buf}
	return m2
}

func (m mark) isBefore(m2 mark) bool {
	return m.line < m2.line ||
		(m.line == m2.line && m.pos < m2.pos)
}

// toEndofText returns a mark at then of the text t assuming that it starts at
// the mark m
func (m mark) toEndofText(t text) mark {
	if t.empty() {
		return m
	}

	end := mark{m.line, 0, m.buf}
	if len(t) == 1 && t[0].lastChar() != '\n' {
		end.pos = m.pos + len(t[0])
		return end
	}

	end.line += len(t) - 1
	if t.lastChar() == '\n' {
		end.line++
	} else {
		end.pos = len(t.lastLine())
	}
	return end
}

func (m mark) isSame(m2 mark) bool {
	return m.buf == m2.buf &&
		m.pos == m2.pos &&
		m.line == m2.line
}
