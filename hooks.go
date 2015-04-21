package main

type actionHooks map[filetype][]func(v *view)

var beforeSaveHooks = actionHooks{}

func (bh actionHooks) add(ft filetype, fn func(v *view)) {
	bh[ft] = append(bh[ft], fn)
}
