package util

import "sync"

func NewByteSlicePool(size int) *ByteSlicePool {
	p := new(sync.Pool)
	p.New = func() interface{} {
		return make([]byte, size)
	}

	return &ByteSlicePool{p}
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
