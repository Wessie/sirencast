// +build !go1.3

package util

import "sync"

type Pool interface {
	Get() []byte
	Put([]byte)
}

func NewPool() Pool {
	return NewMutexPool()
}

type buff struct {
	next *buff
	buf  []byte
}

func NewMutexPool() Pool {
	return &MutexPool{
		head: nil,
		New:  nil,
	}
}

type MutexPool struct {
	sync.Mutex

	head *buff
	New  func() []byte
}

func (p *MutexPool) Get() (b []byte) {
	p.Lock()
	defer p.Unlock()

	// Return nil if we have nothing
	b = nil

	if p.head == nil && p.New != nil {
		b = p.New()
	} else if p.head != nil {
		b = p.head.buf
		p.head = p.head.next
	}

	return b
}

func (p *MutexPool) Put(b []byte) {
	p.Lock()
	defer p.Unlock()

	buf := &buff{
		next: p.head,
		buf:  b[:0],
	}

	p.head = buf
}
