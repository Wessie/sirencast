package icecast

import (
	"bufio"
	"io"
	"log"
	"net"
)

type Client struct {
	// conn is the connection of the client
	conn    net.Conn
	bufconn *bufio.Writer

	// meta indicates if this client wants metadata interleaved with the data
	meta bool
	// metaint is the amount of bytes between each metadata section send
	metaint int
}

func (c *Client) runLoop(r io.ReadCloser, m ReadOnlyMetadata) {
	defer c.conn.Close()
	defer r.Close()

	var n, wn int
	var err error
	var p = make([]byte, 16384)

	// use a simpler loop when no meta is requested
	if !c.meta {
		log.Println("icecast.client: using non-meta loop")
		for {
			n, err = r.Read(p)
			if err != nil {
				return
			}

			wn, err = c.bufconn.Write(p[:n])
			if wn != n || err != nil {
				return
			}
		}
	}

	log.Println("icecast.client: using meta loop")
	// we switch to the metaint size for the buffer
	// this allows us to do a read/write/meta cycle in the loop
	// with little extra tracking.
	p = make([]byte, c.metaint)
	var (
		metabuf = make([]byte, 255*16+1)
		// zero is the no-metadata buffer
		zero = []byte{0}
		// variable used for writing
		metadata []byte
		// meta is the temporary metadata variable
		meta string
		// the latest metadata we saw
		curMeta string
	)

	for {
		n, err = io.ReadFull(r, p)
		if err != nil {
			return
		}

		if n != len(p) {
			log.Println("icecast.client: failed to read full buffer")
			return
		}

		wn, err = c.bufconn.Write(p)
		if wn != n || err != nil {
			return
		}

		// handle metadata, we first need to check if we have new metadata
		// at all, we can send a 0 length meta block if we have nothing.
		if meta = m.Get(); meta == curMeta {
			metadata = zero
		} else {
			curMeta = meta
			metadata = fillMetaBuffer(metabuf, curMeta)
		}

		wn, err = c.bufconn.Write(metadata)
		if wn != len(metadata) || err != nil {
			return
		}
	}
}

var (
	metaFront   = []byte("StreamTitle='")
	metaBack    = []byte("';")
	metaPadding = make([]byte, 16)
)

func fillMetaBuffer(m []byte, meta string) []byte {
	p := m[1:1]
	p = append(p, metaFront...)
	p = append(p, meta...)
	p = append(p, metaBack...)
	p = append(p, getPadding(len(p))...)
	m[0] = byte(len(p) / 16)

	return m[:len(p)+1]
}

func getPadding(length int) []byte {
	return metaPadding[:calculatePadding(length)]
}

func calculatePadding(length int) int {
	return 16 - (length % 16)
}
