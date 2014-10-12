package icecast

import (
	"io"
	"net"
)

type Client struct {
	// conn is the connection of the client
	conn net.Conn
	// meta indicates if this client wants metadata interleaved with the data
	meta bool
	// metaint is the amount of bytes between each metadata section send
	metaint int
}

func (c *Client) runLoop(r io.Reader, m Metadata) {

}
