// +build go1.3

package util

import "sync"

type Pool interface {
	Get() []byte
	Put([]byte)
}

func NewPool() Pool {
	return NewByteSlicePool()
}

func NewByteSlicePool() Pool {
	return &ByteSlicePool{new(sync.Pool)}
}

type ByteSlicePool struct {
	pool *sync.Pool
}

func (p *ByteSlicePool) Get() []byte {
	return p.pool.Get().([]byte)
}

func (p *ByteSlicePool) Put(b []byte) {
	p.pool.Put(b)
}
