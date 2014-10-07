package icecast

import (
	"bufio"
	"io"
	"log"
	"net/http"
	"net/textproto"
	"net/url"
	"sync"

	"github.com/Wessie/sirencast"
	"github.com/Wessie/sirencast/util/taxtic"
)

func NewServer() *Server {
	return &Server{
		mu:     new(sync.RWMutex),
		mounts: make(map[string]*Mount),
	}
}

type Server struct {
	mu     *sync.RWMutex
	mounts map[string]*Mount
}

type ReadWriteCloser struct {
	*bufio.Reader
	*bufio.Writer
	io.Closer
}

// SourceHandler parses an incoming request from an icecast
// source client.
func (s *Server) SourceHandler(conn *sirencast.Conn) {
	// A request from an icecast source looks similar to plain HTTP.
	//
	// The initial line is of the form:
	// 	SOURCE /mountpoint ICE/1.0
	//
	// The HTTP version part of this line (ICE/1.0) is the only thing
	// holding us back from using the builtin net/http request parsers.
	//
	// Instead we parse the first line ourself, and construct a http.Request
	// from this.
	b := &ReadWriteCloser{
		Reader: bufio.NewReader(conn),
		Writer: bufio.NewWriter(conn),
		Closer: conn,
	}

	line, err := b.ReadString('\n')
	if err != nil {
		log.Println("source: invalid http request")
		// TODO: Log errors
		return
	}

	method, uri, proto, ok := parseRequestLine(line)
	if !ok {
		log.Println("source: invalid http request line")
		// TODO: Log errors
		return
	}

	if method != "SOURCE" {
		log.Println("source: received non-source method request.")
		return
	}

	u, err := url.ParseRequestURI(uri)
	if err != nil {
		log.Println("source: invalid http request uri: ", err)
		// TODO: Log errors
		return
	}

	tp := textproto.NewReader(b.Reader)
	mimeHeader, err := tp.ReadMIMEHeader()
	if err != nil {
		log.Println("source: invalid http headers: ", err)
		// TODO: Log errors
		return
	}

	req := &http.Request{
		Body:       b,
		Method:     method,
		Proto:      proto,
		ProtoMajor: 1,
		ProtoMinor: 0,
		Header:     http.Header(mimeHeader),
		RequestURI: uri,
		URL:        u,
		Host:       u.Host,
		RemoteAddr: conn.RemoteAddr().String(),
	}

	// Adjust the Host field if it wasn't included in the URI but we
	// do have a Host header.
	if req.Host == "" {
		req.Host = req.Header.Get("Host")
	}

	if err := WriteHeader(b, req.Header, http.StatusOK); err != nil {
		// TODO: Log errors
		return
	}

	s.Mount(u.Path).AddSource(NewSource(b, req))
	return
}

func (s *Server) MetadataHandler(conn *sirencast.Conn) {
	b := bufio.NewReader(conn)
	r, err := http.ReadRequest(b)
	if err != nil {
		// TODO: Log errors
		return
	}

	query := r.URL.Query()

	metadata := query.Get("song")
	if metadata == "" {
		// TODO: Decide on ignoring empty values or updating to empty
		return
	}

	charset := query.Get("charset")
	if charset == "" {
		// TODO: Check what we want to do for encoding, defaulting to utf8 is
		// pretty sane, but might not be the correct approach for icecast compatibility.
		charset = "utf8"
	}

	metadata, err = taxtic.Convert(charset, metadata)
	if err != nil {
		// TODO: Log errors
		return
	}

	id := NewSourceID(r)

	s.Mount(r.URL.Path).SetMetadata(id, metadata)
	return
}

func (s *Server) ClientHandler(conn *sirencast.Conn) {
	b := bufio.NewReader(conn)
	r, err := http.ReadRequest(b)
	if err != nil {
		log.Println("icecast: client init failure:", err)
		return
	}

	_ = r
	return
}

func (s *Server) ListClientHandler(conn *sirencast.Conn) {
	return
}

func (s *Server) MountExists(name string) bool {
	_, ok := s.mounts[name]
	return ok
}

func (s *Server) Mount(name string) *Mount {
	s.mu.Lock()
	m, ok := s.mounts[name]
	if !ok {
		m = NewMount(name)
		s.mounts[name] = m
	}
	s.mu.Unlock()

	return m
}
