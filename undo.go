package main

import "fmt"

const MAX_UNDO = 10000

type changeList struct {
	ops      []bufferChange
	current  int  // the position of the next op to undo (0 if none)
	total    int  // the total number of actions in the register
	redoMode bool // if true we are redoing an action, no need to record it
}

type undoContext struct {
	action undoAction
	text   []line
	start  mark
	end    mark
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
		c.ops = append(c.ops[:c.current+1], bufferChange{*redoCtx(&redo), undo})
		c.current++
	}
}

func (c *changeList) undo() string {
	if c.current == 0 {
		return "No more changes to undo"
	}
	ctx := c.ops[c.current].undo
	switch ctx.action {
	case undoDelete:
		ctx.start.insertText(ctx.text)
	case undoWrite:
		region{ctx.start, ctx.end}.delete(right)
	case undoReplace:
		region{ctx.start, ctx.end}.delete(right)
		ctx.start.insertText(ctx.text)
	}
	ctx.start.buf.cs = ctx.start
	c.current--
	return fmt.Sprintf("undid change #%v of %v", c.current, c.total)
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
	return fmt.Sprintf("redid change #%v of %v", c.current, c.total)
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

func redoCtx(ctx *cmdContext) *cmdContext {
	p := *ctx.point
	ctx.point = &p
	ctx.silent = true
	return ctx
}
