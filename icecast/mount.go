package icecast

import (
	"io"
	"net/http"
)

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
	return &Source{rwc, r}
}

type Source struct {
	io.ReadWriteCloser
	req *http.Request
}

func (s *Source) ID() SourceID {
	return NewSourceID(s.req)
}

type Mount struct {
	Name string
}

func NewMount(name string) *Mount {
	return &Mount{
		Name: name,
	}
}

func (m *Mount) AddSource(s *Source) {
	return
}

// SetMetadata sets the metadata of the source bound to the given
// SourceID. Setting an empty string means deleting the current
// metadata.
func (m *Mount) SetMetadata(id SourceID, metadata string) {
	return
}
