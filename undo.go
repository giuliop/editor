package main

import "fmt"

const MAX_UNDO = 10000

type changeList struct {
	ops      []bufferChange
	current  int  // the position of the next op to undo (0 if none)
	redoMode bool // if true we are redoing an action, no need to record it
}

type undoContext struct {
	text  []line
	start mark
	end   mark
}

type undoAction int

const (
	undoDelete undoAction = iota
	undoWrite
	undoReplace
)

type bufferChange struct {
	redo cmdContext
	undo undoContext
}

func (c *changeList) add(redo cmdContext, undo undoContext) {
	if !c.redoMode {
		if c.current == MAX_UNDO {
			c.ops = c.ops[1:]
			c.current--
		}
		c.current++
		c.ops = append(c.ops[:c.current], bufferChange{newRedoCtx(&redo), undo})
	}
}

func (c *changeList) undo() string {
	if c.current == 0 {
		return "No more changes to undo"
	}
	ctx := c.ops[c.current].undo
	// if ctx.end is not set we don't need to delete text
	if ctx.end.buf != nil {
		region{ctx.start, ctx.end}.delete()
	}
	if !text(ctx.text).empty() {
		ctx.start.insertText(ctx.text)
	} else {
		// if we can we move left the cursor to place it before the deleted text
		if !ctx.start.atLineStart() {
			ctx.start.moveLeft(1)
		}
	}

	ctx.start.buf.cs = ctx.start
	c.current--
	return fmt.Sprintf("undid change #%v of %v", c.current+1, len(c.ops)-1)
}

func (c *changeList) redo() string {
	c.redoMode = true
	if c.current == len(c.ops)-1 {
		return "Already at latest change"
	}
	c.current++
	ctx := c.ops[c.current].redo
	p := *ctx.point
	ctx.point = &p
	pushCmd(&ctx)
	ctx.point.buf.cs = *ctx.point
	c.redoMode = false
	return fmt.Sprintf("redid change #%v of %v", c.current, len(c.ops)-1)
}

func undo(ctx *cmdContext) {
	for i := 0; i < ctx.num; i++ {
		ctx.msg = ctx.point.buf.changeList.undo()
	}
}

func redo(ctx *cmdContext) {
	for i := 0; i < ctx.num; i++ {
		ctx.msg = ctx.point.buf.changeList.redo()
	}
}

func newRedoCtx(ctx *cmdContext) cmdContext {
	p := *ctx.point
	ctx.point = &p
	ctx.silent = true
	return *ctx
}
