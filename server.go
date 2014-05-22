package sirencast

import (
	"errors"
	"net"
	"time"
)

type Server struct {
	Environment *Environment

	Detectors *Detectors

	listener net.Listener
}

func SetupServer(e *Environment) (*Server, error) {
	s := &Server{
		Environment: e,
		Detectors:   DefaultDetectors,
	}

	return s, nil
}

func (server *Server) Serve() (err error) {
	addr, err := net.ResolveTCPAddr("tcp", server.Environment.Server.BindAddress)

	if err != nil {
		return err
	}

	l, err := net.ListenTCP("tcp", addr)

	if err != nil {
		return err
	}

	server.listener = l

	var tempDelay time.Duration
	for {
		conn, err := l.Accept()

		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}

				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}

				time.Sleep(tempDelay)
			} else {
				// Unrecoverable
				return err
			}
		}

		c, err := server.newConn(conn)

		if err != nil {
			conn.Close()
			continue
		}

		go c.serve()
	}

	return nil
}

func (server *Server) newConn(c net.Conn) (sc *SirenConn, err error) {
	sc = &SirenConn{}

	sc.conn = c
	sc.peeker = NewPeeker(c)
	sc.reader = sc.peeker

	sc.handler = server.Detectors.Detect(sc.peeker)

	if sc.handler == nil {
		return nil, errors.New("Unsupported stream")
	}

	return sc, nil
}

func Run(e *Environment) error {
	server, err := SetupServer(e)

	if err != nil {
		return err
	}

	return server.Serve()
}
