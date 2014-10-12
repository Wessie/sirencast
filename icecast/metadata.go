package icecast

import "sync"

type ReadOnlyMetadata interface {
	Get() string
}

func NewMetadata() *Metadata {
	return new(Metadata)
}

type Metadata struct {
	meta string
	mu   sync.Mutex
}

func (m *Metadata) Get() string {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.meta
}

func (m *Metadata) Set(s string) {
	m.mu.Lock()
	m.meta = s
	m.mu.Unlock()
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
