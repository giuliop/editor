package main

type indentFunc func(m *mark)

var indentFuncs = map[filetype]indentFunc{
	_go: goindent,
}

func goindent(m *mark) {

}

func (m *mark) indentLine() {
	indentFuncs[m.buf.filetype](m)
}
