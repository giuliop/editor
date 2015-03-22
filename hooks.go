package main

type bufferActionHooks map[filetype][]func(b *buffer)

var beforeSaveHooks = bufferActionHooks{}

func (bh bufferActionHooks) add(ft filetype, fn func(b *buffer)) {
	bh[ft] = append(bh[ft], fn)
}
