package main

import "testing"

func TestGoSupport(t *testing.T) {
	a := &asserter{}
	// test openBlock
	ss := []string{
		"func a() {",
		"func a() {  ",
		"var (",
		"var ( \t",
		"case xx == yy:",
	}
	for _, s := range ss {
		a.assert("openBlock", s, openBlock.Match([]byte(s)), true)
	}
	// test closeBlock
	ss = []string{
		"}",
		")",
		"  )",
		"\t  }",
		"  \t\t }",
	}
	for _, s := range ss {
		a.assert("closeBlock", s, closeBlock.Match([]byte(s)), true)
	}
	// test caseStatement
	ss = []string{
		"   case xxx == yyy:   ",
		"\t\tdefault:",
	}
	for _, s := range ss {
		a.assert("caseStatement", s, caseStatement.Match([]byte(s)), true)
	}
	if a.failed {
		for _, m := range a.errMsgs {
			t.Error(m)
		}
	}
}
