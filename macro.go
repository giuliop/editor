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
	if r.macro.on {
		// save the macro keys removing the last key which is end record key
		keys := r.macro.keys[:len(r.macro.keys)-1]
		r.macro.macros[0] = keys
		debug.Printf("macro:\n%v\n", keypressesToEmitString(keys))
		r.macro.stop()
		ctx.msg = "finished recording"
		return
	}
	r.macro.start()
	ctx.msg = "started macro recording"
}
