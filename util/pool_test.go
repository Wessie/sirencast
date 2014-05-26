// +build !go1.3

package util

import (
	"runtime"
	"sync"
	"sync/atomic"
	"testing"
)

func TestPoolNew(t *testing.T) {
	p := NewMutexPool()
	p.(*MutexPool).New = func() []byte { return make([]byte, 4096) }

	if p.Get() == nil {
		t.Fatal("Expected new slice")
	}
}

func TestPoolNoNew(t *testing.T) {
	var p Pool = &MutexPool{}

	if p.Get() != nil {
		t.Fatal("expected empty")
	}

	p.Put([]byte{20})
	p.Put([]byte{10})

	if g := p.Get(); g == nil {
		t.Fatalf("got nil, wanted slice", g)
	}

	if g := p.Get(); g == nil {
		t.Fatalf("got nil, wanted slice", g)
	}

	if g := p.Get(); g != nil {
		t.Fatal("expected empty")
	}
}

func BenchmarkMutexPool(b *testing.B) {
	var p MutexPool
	var wg sync.WaitGroup

	n0 := uintptr(b.N)
	n := n0
	for i := 0; i < runtime.GOMAXPROCS(0); i++ {
		wg.Add(1)

		go func() {
			defer wg.Done()
			in := []byte{1}

			for atomic.AddUintptr(&n, ^uintptr(0)) < n0 {
				for b := 0; b < 100; b++ {
					p.Put(in)
					p.Get()
				}
			}
		}()
	}

	wg.Wait()
}

func BenchmarkMutexPoolOverflow(b *testing.B) {
	var p MutexPool
	var wg sync.WaitGroup

	in := []byte{1}
	n0 := uintptr(b.N)
	n := n0

	for i := 0; i < runtime.GOMAXPROCS(0); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			for atomic.AddUintptr(&n, ^uintptr(0)) < n0 {
				for b := 0; b < 100; b++ {
					p.Put(in)
				}

				for b := 0; b < 100; b++ {
					p.Get()
				}
			}
		}()
	}
	wg.Wait()
}
