package main

type indentFunc func(m *mark) (indent int)

var (
	tab     = []rune{'\t'}
	tabStop = 4
)

var indentFuncs = map[filetype]indentFunc{}

var indentKeys = map[filetype][]rune{}

func isIndentKey(r rune, b *buffer) bool {
	for _, k := range indentKeys[b.filetype] {
		if k == r {
			return true
		}
	}
	return false
}

// indentLine applies identation by getting the indentation from the appropriate
// file type indent func if present or simply the prevoius line indent; it then
// returns the change in indentation chars to move the cursor if needed
func (m *mark) indentLine() (indentChars int) {

	_, currIndentChars := lineIndent(m.buf, m.line)
	// determine indentation
	f := indentFuncs[m.buf.filetype]
	var indent int
	switch {
	case f != nil:
		indent = f(m)
	case m.line == 0:
		indent = 0
	default:
		indent, _ = lineIndent(m.buf, m.line-1)
	}

	tabs, spaces := indent/tabStop, indent%tabStop
	indentRunes := line{}
	for i := 0; i < tabs; i++ {
		indentRunes = append(indentRunes, tab...)
	}
	for i := 0; i < spaces; i++ {
		indentRunes = append(indentRunes, ' ')
	}
	_, oldIndent := lineIndent(m.buf, m.line)
	m.buf.text[m.line] = append(indentRunes, m.buf.text[m.line][oldIndent:]...)
	return tabs + spaces - currIndentChars
}

// indent returns the indentation of the line and the numbers of indent chars
func lineIndent(b *buffer, ln int) (indent, indentChars int) {
	for _, r := range b.text[ln] {
		switch r {
		case '\t':
			indent += tabStop
			indentChars++
		case ' ':
			indent += 1
			indentChars++
		default:
			return indent, indentChars
		}
	}
	return indent, indentChars
}
