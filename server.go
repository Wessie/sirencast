package sirencast

import (
	"errors"
	"net"
	"time"

	"github.com/Wessie/sirencast/config"
)

type Server struct {
	Config    *config.Config
	Detectors *Detectors
	listener  net.Listener
}

func SetupServer(e *config.Config) (*Server, error) {
	s := &Server{
		Config:    e,
		Detectors: DefaultDetectors,
	}

	return s, nil
}

func (server *Server) Serve() (err error) {
	addr, err := net.ResolveTCPAddr("tcp", server.Config.Addr)
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

// newConn wraps the given connection into a Conn and tries
// to find a handler suitable for the connection.
//
// newConn will return an error if it is unable to find a handler.
func (server *Server) newConn(c net.Conn) (*Conn, error) {
	var (
		p = NewPeeker(c)
		h ConnHandler
	)

	if h = server.Detectors.Detect(p); h == nil {
		return nil, errors.New("Unsupported stream")
	}

	return &Conn{
		conn:    c,
		start:   p,
		reader:  p,
		handler: h,
	}, nil
}

// Run runs a sirencast server with the configuration given. This
// is a blocking function and will return the error returned by
// SetupServer if any, otherwise it will return the error from Server.Serve
func Run(e *config.Config) error {
	server, err := SetupServer(e)

	if err != nil {
		return err
	}

	return server.Serve()
}
