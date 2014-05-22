package sirencast

import (
	"io"
	"log"
	"net"
	"runtime"
	"time"
)

// SirenConn wraps a net.Conn to support peeking at the front of the stream.
// This is done so that we can detect what kind of content is arriving before
// giving it off to the correct handler.
type SirenConn struct {
	conn   net.Conn
	peeker Peeker

	// reader is the current reader to read from, either `conn` or `peeker`
	reader io.Reader
	// handler is the handler that is called when `serve` is called.
	handler ConnHandler
}

func (sc *SirenConn) serve() {
	defer func() {
		if err := recover(); err != nil {
			buf := make([]byte, 4096)
			buf = buf[:runtime.Stack(buf, false)]
			log.Printf("sirencast: panic serving %s: %v\n%s", sc.RemoteAddr(), err, buf)
		}
	}()

	sc.handler(sc)
}

func (sc *SirenConn) Read(b []byte) (n int, err error) {
	n, err = sc.reader.Read(b)

	if err != nil {
		if sc.reader == sc.peeker {
			sc.reader = sc.conn

			// If we read some bytes before, return those
			if n > 0 {
				return
			}

			// otherwise proxy ourself to the new reader
			return sc.reader.Read(b)
		}

		return
	}

	return
}

func (sc *SirenConn) Write(b []byte) (n int, err error) {
	return sc.conn.Write(b)
}

func (sc *SirenConn) Close() error {
	return sc.conn.Close()
}

func (sc *SirenConn) LocalAddr() net.Addr {
	return sc.conn.LocalAddr()
}

func (sc *SirenConn) RemoteAddr() net.Addr {
	return sc.conn.RemoteAddr()
}

func (sc *SirenConn) SetDeadline(t time.Time) error {
	return sc.conn.SetDeadline(t)
}

func (sc *SirenConn) SetReadDeadline(t time.Time) error {
	return sc.conn.SetReadDeadline(t)
}

func (sc *SirenConn) SetWriteDeadline(t time.Time) error {
	return sc.conn.SetWriteDeadline(t)
}
