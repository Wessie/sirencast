package util

import "io"

// RingBufferSize is the byte slice allocation size
var RingBufferSize = 16384

// RingBuffer is an io.ReadWriter implementor that drops the
// oldest write when a reader is too slow.
type RingBuffer struct {
	readCache []byte
	buf       chan []byte
	pool      *ByteSlicePool
	dropped   uint64
}

// NewRingBuffer allocates a new RingBuffer with the given amount
// of spots. Each spot is equal to one Write, therefore it is
// suggested to wrap the RingBuffer with the bufio package.
func NewRingBuffer(spots int) *RingBuffer {
	return &RingBuffer{
		buf:  make(chan []byte, spots),
		pool: NewByteSlicePool(RingBufferSize),
	}
}

// Write writes to the buffer and drops older writes if
// they have not been read yet. Write is non-blocking and
// will always complete.
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

// Read reads from buffer, if no data is available Read waits
// until there is some available.
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
