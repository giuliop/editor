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
	if !(m.pos > m.lastCharPos()) {
		m.pos += 1
	} else {
		if !m.atLastLine() {
			m.set(m.line+1, 0)
		}
	}
}

func (m *mark) moveLeft(steps int) {
	if !m.atLineStart() {
		m.pos -= 1
	} else {
		if !m.atFirstLine() {
			m.set(m.line-1, len(m.buf.text[m.line-1])-1)
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
