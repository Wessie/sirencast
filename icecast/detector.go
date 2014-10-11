package icecast

import (
	"bufio"
	"io"
	"log"
	"net/url"
	"strings"

	"github.com/Wessie/sirencast"
)

func init() {
	//	sirencast.RegisterDetector(Detect)
}

func (s *Server) Detect(r io.Reader) sirencast.ConnHandler {
	// TODO: Optimize this, the bufio.Reader is a bit heavy
	b := bufio.NewReader(r)
	line, err := b.ReadString('\n')
	if err != nil {
		log.Println("icecast.detector: failed to read first line:", err)
		return nil
	}

	method, uri, _, ok := parseRequestLine(line)
	if !ok {
		log.Println("icecast.detector: failed to parse request line")
		return nil
	}

	if method == "SOURCE" {
		return s.SourceHandler
	}

	// All handlers below expect a GET request, so we can return
	// early if this isn't a GET
	if method != "GET" {
		return nil
	}

	// Check for '/admin/listclients' and '/admin/metadata' request,
	// both are special for icecast. Anything else we try as a client
	// requesting a mountpoint.
	u, err := url.ParseRequestURI(uri)
	if err != nil {
		return nil
	}

	if u.Path == "/admin/listclients" {
		return s.ListClientHandler
	} else if u.Path == "/admin/metadata" {
		return s.MetadataHandler
	}

	// Check for a new client
	if s.MountExists(u.Path) {
		return s.ClientHandler
	}

	return nil
}

func parseRequestLine(line string) (method, requestURI, proto string, ok bool) {
	s1 := strings.Index(line, " ")
	s2 := strings.Index(line[s1+1:], " ")
	if s1 < 0 || s2 < 0 {
		return
	}
	s2 += s1 + 1
	return line[:s1], line[s1+1 : s2], line[s2+1:], true
}
