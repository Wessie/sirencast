package icecast

import (
	"log"

	"github.com/Wessie/sirencast/util"
)

// Mount depicts a singular icecast mountpoint. A mountpoint can have
// many clients (same as plain icecast) and have many sources (not the
// same as icecast).
type Mount struct {
	Name    string
	sources *Container
	meta    Metadata

	mw *MultiWriter
}

func NewMount(name string) *Mount {
	return &Mount{
		Name:    name,
		sources: NewContainer(),
		meta:    NewMetadataContainer(),
		mw:      NewMultiWriter(),
	}
}

func (m *Mount) AddClient(c *Client) {
	m.log("adding client: %v", c)

	r := util.NewRingBuffer(5)
	m.mw.Add(r)
	go c.runLoop(r, m.meta)
}

// AddSource adds a new source to the mountpoint, the mountpoint will
// be responsible for sources output and removal after disconnection
func (m *Mount) AddSource(s *Source) {
	m.log("adding source: %v", s)
	m.sources.Add(s)

	go func() {
		// read from the source and remove when it returns
		s.readLoop()
		m.sources.Remove(s)
		m.log("removing source: %v", s)
	}()
}

// SetMetadata sets the metadata of the source bound to the given
// SourceID. Setting an empty string means deleting the current
// metadata.
func (m *Mount) SetMetadata(id SourceID, metadata string) {
	m.log("setting metadata: id: %s meta: %s", id, metadata)
	m.meta.Set(id, metadata)
	return
}

func (m *Mount) log(f string, args ...interface{}) {
	f = "icecast.mount." + m.Name + ": " + f
	log.Printf(f, args...)
}
