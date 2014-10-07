package icecast

import (
	"io"
	"log"
	"net/http"
	"sync"
)

type nullWriter struct{}

func (nw nullWriter) Write(p []byte) (int, error) {
	return len(p), nil
}

var discardWriter = nullWriter{}
var ReadBufferSize = 4096

func NewSourceID(r *http.Request) SourceID {
	id := SourceID{}

	if r.URL.Path == "/admin/metadata" || r.URL.Path == "/admin/listclients" {
		id.Mount = r.URL.Query().Get("mount")
	} else {
		id.Mount = r.URL.Path
	}

	return id
}

type SourceID struct {
	Mount string
}

func NewSource(rwc io.ReadWriteCloser, r *http.Request) *Source {
	s := &Source{
		ReadWriteCloser: rwc,
		req:             r,
		out:             discardWriter,
	}
	go s.readLoop()
	return s
}

// Source is an icecast source client, a source sends audio data and
// metadata of this audio to be send to listening clients.
type Source struct {
	// source input/output
	io.ReadWriteCloser
	// source request
	req *http.Request

	// protects 'out' below
	mu sync.Mutex
	// mount output
	out io.Writer
	// source name
	Name string
}

// ID returns the SourceID generated by the sources initial request,
// this identifier can be thought of as unique to a source.
func (s *Source) ID() SourceID {
	return NewSourceID(s.req)
}

func (s *Source) readLoop() {
	b := make([]byte, ReadBufferSize)
	for {
		n, err := s.Read(b)
		if err != nil {
			log.Println("reading error:", err)
			return
		}

		s.mu.Lock()
		_, err = s.out.Write(b[:n])
		s.mu.Unlock()
		if err != nil {
			log.Println("writing error:", err)
			return
		}

	}
}

// SwapOut swaps the source output with the new writer passed in.
func (s *Source) SwapOut(new io.Writer) {
	s.mu.Lock()
	s.out = new
	s.mu.Unlock()
}