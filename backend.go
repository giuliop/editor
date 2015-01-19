package main

// backend holds the buffers open in the editor
type backend struct {
	bufs []buffer // the open buffers
}

// initBackend returns the backend after having initialized it
func initBackend() backend {
	be := backend{}
	return be
}

// newBuffer adds a new empty buffer to the backend and returns a pointer to it
// Note that the last line of a buffer ends with a newline which is removed before
// saving to file
func (be *backend) newBuffer(name string) *buffer {
	b := &buffer{
		text: make([]line, 1, 20),
		name: name,
	}
	b.cs = newMark(b)
	b.text[0] = newLine()
	be.bufs = append(be.bufs, *b)
	return b
}
