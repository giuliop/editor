package main

type mark struct {
	line int
	pos  int
	buf  *buffer
}

func newMark(b *buffer) mark {
	return mark{0, 0, b}
}

func (m *mark) atFirstLine() bool {
	return m.line == 0
}

func (m *mark) atLastLine() bool {
	return m.line == m.maxLine()
}

func (m *mark) atLineStart() bool {
	return m.pos == 0
}

// lastCharPos return the position of the last char in the line (before the newline
// char if present. If the line is empty it returns -1
func (m *mark) lastCharPos() int {
	last := len(m.buf.text[m.line]) - 1
	if last >= 0 && m.buf.text[m.line][last] == '\n' {
		last -= 1
	}
	return last
}

func (m *mark) emptyLine() bool {
	return m.lastCharPos() == -1
}

func (m *mark) maxLine() int {
	return len(m.buf.text) - 1
}

// fixPos checks that the mark is within the line, if it is over the end of the line
// puts the cursor back to the end, if before start of line, back to start of line
func (m *mark) fixPos() {
	// if in insert mode max position is after the last char
	max := m.lastCharPos() + 1
	// in normal mode we want the cursor one position to the left, unless it is an empty line
	if m.buf.mod == normalMode && max > 0 {
		max -= 1
	}
	if m.pos > max {
		m.pos = max
	}
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
	if m.line > m.maxLine() {
		m.line = m.maxLine()
	}
	m.fixPos()
}

func (m *mark) moveRight(steps int) {
	maxY := m.maxLine()
	for steps > 0 {
		maxX := m.lastCharPos() + 1
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
	m.fixPos()
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
