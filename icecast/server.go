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

const metadataSuccess = `<?xml version="1.0"?>
<iceresponse><message>Metadata update successful</message><return>1</return></iceresponse>
`

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
	b := ReadWriteCloser{
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

	if err := WriteHeader(b, nil, http.StatusOK); err != nil {
		log.Println("icecast.source: failed to write OK header:", err)
		return
	}

	if err := b.Flush(); err != nil {
		log.Println("icecast.source: failed to flush header:", err)
	}

	s.Mount(u.Path).AddSource(NewSource(b, req))
	return
}

func (s *Server) MetadataHandler(conn *sirencast.Conn) {
	r, err := ReadRequest(conn)
	if err != nil {
		log.Println("icecast.metadata: failed to construct request:", err)
		return
	}

	query := r.URL.Query()

	mount := query.Get("mount")
	if mount == "" {
		return
	}

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
		log.Println("icecast.metadata: failed to convert metadata to utf8:", err)
		return
	}

	id := NewSourceID(r)

	s.Mount(mount).SetMetadata(id, metadata)

	h := http.Header{
		"Content-Type":   {"text/xml"},
		"Content-Length": {"113"},
	}

	// now send back a xml "success" response
	if err := WriteHeader(conn, h, http.StatusOK); err != nil {
		log.Println("icecast.metadata: failed to write http header response:", err)
	}

	if _, err := io.WriteString(conn, metadataSuccess); err != nil {
		log.Println("icecast.metadata: failed to write xml success response:", err)
	}
	return
}

func (s *Server) ClientHandler(conn *sirencast.Conn) {
	r, err := ReadRequest(conn)
	if err != nil {
		log.Println("icecast.client: failed to read http request:", err)
		return
	}

	_ = r
	return
}

func (s *Server) ListClientHandler(conn *sirencast.Conn) {
	return
}

func (s *Server) MountExists(name string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
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
