package sirencast

import (
	"errors"
	"net"
)

// NewHTTPListener returns a new HTTPListener
func NewHTTPListener(addr string) *HTTPListener {
	l := &HTTPListener{
		pipe: make(chan *Conn),
		addr: addr,
	}

	l.Handler = func(conn *Conn) {
		l.pipe <- conn
	}

	return l
}

// HTTPListener is a net.Listener that receives connections from
// a ConnHandler rather than from a network listener. The handler
// should be passed to the appropiate register function by the
// user.
type HTTPListener struct {
	addr string
	pipe chan *Conn

	Handler ConnHandler
}

// Accept waits for and returns the next connection
func (l *HTTPListener) Accept() (net.Conn, error) {
	sc, ok := <-l.pipe

	if !ok {
		return nil, errors.New("closed listener pipe")
	}

	return sc, nil
}

// Close is a no-op and always returns nil
func (l *HTTPListener) Close() error {
	return nil
}

func (l *HTTPListener) Addr() net.Addr {
	return fakeAddr{network: "internal", addr: l.addr}
}

type fakeAddr struct {
	network string
	addr    string
}

func (f fakeAddr) Network() string {
	return f.network
}

func (f fakeAddr) String() string {
	return "shared(" + f.addr + ")"
}
