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

func (m *mark) lastCharPos() int {
	return len(m.buf.text[m.line]) - 1
}

func (m *mark) maxLine() int {
	return len(m.buf.text) - 1
}

func (m *mark) checkPos() {
	if m.pos > m.lastCharPos() {
		m.pos = m.lastCharPos()
	}
	if m.pos < 0 {
		m.pos = 0
	}
}

func (m *mark) moveUp(steps int) {
	m.line -= steps
	if m.line < 0 {
		m.line = 0
	}
	m.checkPos()
}

func (m *mark) moveDown(steps int) {
	m.line += steps
	if m.line > m.maxLine() {
		m.line = m.maxLine()
	}
	m.checkPos()
}

func (m *mark) moveRight(steps int) {
	maxY := m.maxLine()
	for steps > 0 {
		maxX := m.lastCharPos()
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
