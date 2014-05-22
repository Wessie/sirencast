package util

import "io"

type RingBuffer struct {
	readCache []byte
	buf       chan []byte
	dropped   uint64
}

func NewRingBuffer(size int) *RingBuffer {
	return &RingBuffer{
		buf: make(chan []byte, size),
	}
}

func (r *RingBuffer) Write(b []byte) (n int, err error) {
	// shortcut and we don't want the reader to think
	// we've reached end of stream
	if b == nil {
		return 0, nil
	}

loop:
	for {
		select {
		case r.buf <- b:
			break loop
		default:
			select {
			case <-r.buf:
			default:
			}
		}
	}

	return len(b), nil
}

func (r *RingBuffer) Read(p []byte) (n int, err error) {
	var b = r.readCache
	if len(b) == 0 {
		b = <-r.buf
	}

	// Check for a closed channel
	if b == nil {
		return 0, io.EOF
	}

	n = copy(p, b)

	r.readCache = b[n:]

	return n, nil
}
