package icecast

import (
	"bufio"
	"net"
	"net/http"
)

// ReadRequest is a light wrapper around http.ReadRequest
func ReadRequest(conn net.Conn) (*http.Request, error) {
	b := bufio.NewReader(conn)
	r, err := http.ReadRequest(b)
	if err != nil {
		return nil, err
	}

	r.RemoteAddr = conn.RemoteAddr().String()
	return r, nil
}
