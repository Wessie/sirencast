package sirencast

import (
	"io"
)

const PeekBufferSize = 8192
const PeekReadSize = 4096

type Peeker interface {
	io.Reader
	// Reset peeking position to the start
	Reset()
	// Command the peeker to stop peeking, this makes the peeker read until
	// the buffer end is reached and return io.EOF instead of reading more
	// from the input.
	Stop()
}

type PeekReader struct {
	input      io.Reader
	buffer     []byte
	wpos, rpos int
	stopped    bool
}

func NewPeeker(input io.Reader) Peeker {
	return &PeekReader{
		input:   input,
		buffer:  nil,
		wpos:    0,
		rpos:    0,
		stopped: false,
	}
}

func (pk *PeekReader) Read(p []byte) (n int, err error) {
	if pk.buffer == nil {
		pk.buffer = make([]byte, PeekBufferSize)
	}

	if (pk.wpos - pk.rpos) > 0 {
		// Copy over as much as we can from the buffer
		n = copy(p, pk.buffer[pk.rpos:pk.wpos])

		pk.rpos += n

		if n > 0 {
			return n, nil
		}
	}

	if pk.stopped {
		return 0, io.EOF
	}

	// We ran out of bytes in the buffer, so instead get ready to
	// read from the input reader.
	var buffer []byte
	if pk.wpos > (PeekBufferSize - PeekReadSize) {
		buffer = make([]byte, PeekReadSize)
		pk.buffer = append(pk.buffer, buffer...)
	}

	buffer = pk.buffer[pk.wpos : pk.wpos+PeekReadSize]

	// Otherwise read from the original source
	n, err = pk.input.Read(buffer)

	pk.wpos += n

	if err != nil {
		return
	}

	// Lower our n if we've read more than caller asked for
	if n > len(p) {
		n = len(p)
	}

	copy(p, buffer[:n])

	return
}

func (pk *PeekReader) Reset() {
	pk.rpos = 0
}

func (pk *PeekReader) Stop() {
	pk.stopped = true
}
