package icecast

import "sync"

type Metadata interface {
	Set(SourceID, string)
	Get(SourceID) string
}

func NewMetadataContainer() *MetadataContainer {
	return &MetadataContainer{
		m: make(map[SourceID]string, 1),
	}
}

type MetadataContainer struct {
	mu sync.Mutex
	m  map[SourceID]string
}

func (m *MetadataContainer) Set(id SourceID, meta string) {
	m.mu.Lock()
	m.m[id] = meta
	m.mu.Unlock()
	return
}

func (m *MetadataContainer) Get(id SourceID) (meta string) {
	m.mu.Lock()
	meta = m.m[id]
	m.mu.Unlock()
	return meta
}
