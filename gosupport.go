package main

import (
	"os/exec"
	"regexp"
)

var (
	openBlock     = regexp.MustCompile(`[{(:][[:space:]]*$`)
	closeBlock    = regexp.MustCompile(`^[[:space:]]*[})]`)
	caseStatement = regexp.MustCompile(`^[[:space:]]*(?:case .*)|(?:default):$`)
)

func init() {
	indentFuncs[_go] = goindent
	indentKeys[_go] = []rune{')', '}', ':'}
	commandModeFuncs["gofmt"] = gofmt
	beforeSaveHooks.add(_go, func(v *view) { gofmt(v, nil) })
}

// goIndent returns the indentation needed for the line under the mark
func goindent(m *mark) (indent int) {
	// indent first line in text with no indentation
	if m.line == 0 {
		return 0
	}

	// previous and current line without final newline char
	prev := stripCommentsAndNewline(m.buf.text[m.line-1])
	curr := stripCommentsAndNewline(m.buf.text[m.line])

	// we start from previous line indent
	indent, _ = lineIndent(m.buf, m.line-1)

	// indent line after new block
	if openBlock.Match(prev.toBytes()) {
		indent += tabStop
	}
	// indent closing block line
	if closeBlock.Match(curr.toBytes()) {
		indent -= tabStop
	}
	//outindent back case: or default: in switch statement
	if caseStatement.Match(curr.toBytes()) {
		indent -= tabStop
	}
	if indent < 0 {
		indent = 0
	}
	return indent
}

//  stripCommentsAndNewlinereturns the line without end comments and newline
func stripCommentsAndNewline(ln line) line {
	for i, r := range ln {
		switch r {
		case '/':
			if ln[i+1] == '/' {
				return ln[:i]
			}
		case '\n':
			return ln[:i]
		}
	}
	return ln
}

// gofmt runs goimports (or gofmt if goimports is not installed) on view v
func gofmt(v *view, args []string) string {
	path, err := exec.LookPath("goimports")
	program := "goimports"
	if err != nil {
		path, err = exec.LookPath("gofmt")
		program = "gofmt"
		if err != nil {
			return "gofmt is not installed, sorry!"
		}
	}
	cmd := exec.Command(path, v.buf.filename)
	out, err := cmd.Output()
	if err != nil {
		debug.Printf("%v error: %v\n", program, err)
		return "gofmt error, sorry!"
	}
	v.buf.text = bytesToText(out)
	// make sure cursor is OK
	v.cs.fixLineAndPos()
	return program + " run"
}
