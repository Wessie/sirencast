package sirencast

import (
	"errors"
	"net"
)

// NewHTTPListener returns a new HTTPListener
func NewHTTPListener() *HTTPListener {
	l := &HTTPListener{
		pipe: make(chan *SirenConn),
	}

	l.Handler = func(conn *SirenConn) {
		l.pipe <- conn
	}

	return l
}

// HTTPListener is a net.Listener that receives connections from
// the handler `HTTPListener.Handler`.
type HTTPListener struct {
	pipe     chan *SirenConn
	listener net.Listener

	Handler ConnHandler
}

// Accept waits for and returns the next connection
func (l *HTTPListener) Accept() (net.Conn, error) {
	sc, ok := <-l.pipe

	if !ok {
		return nil, errors.New("Closed listener pipe")
	}

	return sc, nil
}

// Close is a no-op and always returns nil
func (l *HTTPListener) Close() error {
	return nil
}

// Addr always returns nil
func (l *HTTPListener) Addr() net.Addr {
	return fakeAddr{network: "process", addr: "internal"}
}

type fakeAddr struct {
	network string
	addr    string
}

func (f fakeAddr) Network() string {
	return f.network
}

func (f fakeAddr) String() string {
	return f.addr
}
