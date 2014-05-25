package sirencast

import (
	"io"
	"log"
	"net"
	"runtime"
	"time"
)

// Conn wraps a net.Conn to support peeking at the front of the stream.
// This is done so that we can detect what kind of content is arriving before
// giving it off to the correct handler.
type Conn struct {
	conn   net.Conn
	peeker Peeker

	// reader is the current reader to read from, either `conn` or `peeker`
	reader io.Reader
	// handler is the handler that is called when `serve` is called.
	handler ConnHandler
}

// serve calls the appointed handler with sc as argument, it recovers from any panics
// that occur inside the handler to avoid the whole server going down.
func (sc *Conn) serve() {
	defer func() {
		if err := recover(); err != nil {
			buf := make([]byte, 4096)
			buf = buf[:runtime.Stack(buf, false)]
			log.Printf("sirencast: panic serving %s: %v\n%s", sc.RemoteAddr(), err, buf)
		}
	}()

	sc.handler(sc)
}

func (sc *Conn) Read(b []byte) (n int, err error) {
	n, err = sc.reader.Read(b)

	if err != nil {
		// We're already reading from the net.Conn, so
		// we can return whatever the net.Conn returned.
		if sc.reader != sc.peeker {
			return
		}

		// Otherwise we're going to swap to the net.Conn now
		// and continue reading.
		sc.reader = sc.conn

		// If we read some bytes before, return those first
		// before continueing.
		if n > 0 {
			return n, nil
		}

		// otherwise proxy ourself to the new reader
		return sc.reader.Read(b)
	}

	return
}

func (sc *Conn) Write(b []byte) (n int, err error) {
	return sc.conn.Write(b)
}

func (sc *Conn) Close() error {
	return sc.conn.Close()
}

func (sc *Conn) LocalAddr() net.Addr {
	return sc.conn.LocalAddr()
}

func (sc *Conn) RemoteAddr() net.Addr {
	return sc.conn.RemoteAddr()
}

func (sc *Conn) SetDeadline(t time.Time) error {
	return sc.conn.SetDeadline(t)
}

func (sc *Conn) SetReadDeadline(t time.Time) error {
	return sc.conn.SetReadDeadline(t)
}

func (sc *Conn) SetWriteDeadline(t time.Time) error {
	return sc.conn.SetWriteDeadline(t)
}
