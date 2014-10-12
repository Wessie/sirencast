package icecast

import (
	"log"

	"github.com/Wessie/sirencast/util"
)

type mountEvent int

const (
	EventPanic mountEvent = iota
	EventNewSource
	EventRemoveSource
	EventNewMetadata
	EventDestroyMount
)

// Mount depicts a singular icecast mountpoint. A mountpoint can have
// many clients (same as plain icecast) and have many sources (not the
// same as icecast).
type Mount struct {
	Name       string
	sources    *Container
	events     chan mountEvent
	sourceMeta *MetadataContainer

	meta *Metadata
	mw   *MultiWriter
}

func NewMount(name string) *Mount {
	m := Mount{
		Name:       name,
		sources:    NewContainer(),
		sourceMeta: NewMetadataContainer(),
		meta:       NewMetadata(),
		mw:         NewMultiWriter(),
		events:     make(chan mountEvent),
	}
	go m.runLoop()
	return &m
}

func (m *Mount) runLoop() {
	var current *Source
	for {
		switch <-m.events {
		case EventNewSource, EventRemoveSource:
			next := m.sources.Top()
			if next == nil && current != nil {
				current.SwapOutput(discardWriter)
				continue
			}

			if next != current && current != nil {
				current.SwapOutput(discardWriter)
			}

			next.SwapOutput(m.mw)
			current = next
		case EventNewMetadata:
			m.meta.Set(
				m.sourceMeta.Get(current.ID()),
			)
		case EventDestroyMount:
			return
		default:
			panic("icecast.mount: invalid mount event issued")
		}
	}
}

// Close disconnects all sources and clients
func (m *Mount) Close() {

	return
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
	m.events <- EventNewSource

	go func() {
		// read from the source and remove when it returns
		s.readLoop()
		m.sources.Remove(s)
		m.events <- EventRemoveSource
		m.log("removing source: %v", s)
	}()
}

// SetMetadata sets the metadata of the source bound to the given
// SourceID. Setting an empty string means deleting the current
// metadata.
func (m *Mount) SetMetadata(id SourceID, metadata string) {
	m.log("setting metadata: id: %s meta: %s", id, metadata)
	m.sourceMeta.Set(id, metadata)
	m.events <- EventNewMetadata
	return
}

func (m *Mount) log(f string, args ...interface{}) {
	f = "icecast.mount." + m.Name + ": " + f
	log.Printf(f, args...)
}
