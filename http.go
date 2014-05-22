package sirencast


import (
	"net"
	"errors"
)

func NewHTTPListener() *HTTPListener {
	l := &HTTPListener{
		pipe: make(chan *SirenConn),
	}

	l.Handler = func(conn *SirenConn) {
		l.pipe <- conn
	}

	l.Detector = func(input Peeker) ConnHandler {
		return l.Handler
	}

	return l
}

type HTTPListener struct {
	pipe		chan *SirenConn
	listener 	net.Listener

	Handler 	ConnHandler
	Detector 	Detector
}

func (l *HTTPListener) Accept() (net.Conn, error) {
	sc, ok := <-l.pipe

	if !ok {
		return nil, errors.New("Closed listener pipe")
	}

	return sc, nil
}

func (l *HTTPListener) Close() error {
	return nil
}

func (l *HTTPListener) Addr() net.Addr {
	if l.listener == nil {
		return &net.TCPAddr{
			IP: net.IPv4(255, 255, 255, 255),
			Port: 9999,
			Zone: "",
		}
	}

	return l.listener.Addr()
}
