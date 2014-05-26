package util

import "io"

type RingBuffer struct {
	readCache []byte
	buf       chan []byte
	pool      Pool
	dropped   uint64
}

func NewRingBuffer(size int) *RingBuffer {
	return &RingBuffer{
		buf:  make(chan []byte, size),
		pool: NewPool(),
	}
}

func (r *RingBuffer) Write(b []byte) (n int, err error) {
	// shortcut and we don't want the reader to think
	// we've reached end of stream
	if b == nil {
		return 0, nil
	}

	c := append(r.pool.Get()[:0], b...)

loop:
	for {
		select {
		case r.buf <- c:
			break loop
		default:
			select {
			case c := <-r.buf:
				r.pool.Put(c)
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
