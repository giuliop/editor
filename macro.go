package main

type macroRegister struct {
	*keyLogger
	macros [10][]Keypress
}

type keyLogger struct {
	keys []Keypress
	on   bool
}

func (k *keyLogger) start() {
	k.on = true
}

func (k *keyLogger) stop() {
	k.on = false
	k.keys = nil
}

func (k *keyLogger) record(key Keypress) {
	k.keys = append(k.keys, key)
}

func recordMacro(ctx *cmdContext) {
	if r.macros.on {
		// save the macro keys removing the last key which is end record key
		keys := r.macros.keys[:len(r.macros.keys)-1]
		debug.Printf("macro:\n%v\n", keypressesToEmitString(keys))
		debug.Printf("buffer:\n%v\n", _bufferToString(ctx.point.buf))
		r.macros.macros[0] = keys
		r.macros.stop()
		ctx.msg = "finished recording"
		return
	}
	r.macros.start()
	ctx.msg = "started macro recording"
}

func _bufferToString(b *buffer) string {
	s := ""
	for _, line := range b.text {
		s += string(line)
	}
	return s
}
