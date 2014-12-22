package main

type direction int

const (
	right direction = iota
	left
	up
	down
)

type cmdContext struct {
	times  int
	action string
	object string
	dir    direction
	point  mark
}

var ctx = new(cmdContext)

func insertChar() {
	for i := 0; i < ctx.times; i++ {
		//eng.insertChar(ctx.mark, ctx.object)
	}
}
