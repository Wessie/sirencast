package util

import (
	"testing"
	"time"
)

func TestRingReadBlocking(t *testing.T) {
	var (
		testValue = "Hello World"
		r         = NewRingBuffer(16)
		b         = make([]byte, 32)
	)

	go func() {
		time.Sleep(time.Millisecond * 200)

		n, err := r.Write([]byte(testValue))
		if n < len(testValue) || err != nil {
			t.Error("Failed writing:", err)
		}
	}()

	n, err := r.Read(b)
	if err != nil {
		t.Error("Failed reading:", err)
	}

	b = b[:n]

	if string(b) != testValue {
		t.Errorf("Received different value than expected: %s != %s", string(b), testValue)
	}
}

func TestRingWriteNonBlocking(t *testing.T) {
	var (
		testValue  = "Hello World"
		testValueB = []byte(testValue)
		r          = NewRingBuffer(2)
	)

	for i := 0; i < 6; i++ {
		n, err := r.Write(testValueB)
		if n < len(testValue) {
			t.Errorf("Write was too short: %d != %d", n, len(testValue))
		} else if err != nil {
			t.Error("Write returned an error:", err)
		}
	}
}

// TestWriteNonBlocking tests that the buffer drops from the head and
// in a consistent manner
func TestRingWriteDropBehaviour(t *testing.T) {
	var (
		testValueGo   = []byte("Hello World")
		testValueDrop = []byte("Dropping")
		buf           = make([]byte, 128)
		r             = NewRingBuffer(1)
	)

	for i := 0; i < 512; i++ {
		n, err := r.Write(testValueDrop)
		if err != nil {
			t.Fatal("Failed writing (drop):", err)
		} else if n != len(testValueDrop) {
			t.Fatal("Write (drop) was too short")
		}

		n, err = r.Write(testValueGo)
		if err != nil {
			t.Fatal("Failed writing (go):", err)
		} else if n != len(testValueGo) {
			t.Fatal("Write (go) was too short")
		}

		n, err = r.Read(buf)
		if err != nil {
			t.Fatal("Failed reading:", err)
		}

		if string(buf[:n]) != string(testValueGo) {
			t.Fatal("Received different value than expected")
		}
	}
}

func BenchmarkRingWriteDrop(b *testing.B) {
	var (
		r    = NewRingBuffer(16)
		from = make([]byte, 32)
	)
	for i := 0; i < b.N; i++ {
		r.Write(from)
	}
}

func BenchmarkRingFastWriteRead(b *testing.B) {
	var (
		r    = NewRingBuffer(16)
		from = make([]byte, 32)
		to   = make([]byte, 32)
	)

	for i := 0; i < b.N; i++ {
		r.Write(from)
		r.Read(to)
	}
}

func BenchmarkRingWrite2ChunkedRead(b *testing.B) {
	var (
		r    = NewRingBuffer(16)
		from = make([]byte, 32)
		to   = make([]byte, 16)
	)

	for i := 0; i < b.N; i++ {
		r.Write(from)
		r.Read(to)
		r.Read(to)
	}
}

func BenchmarkRingWrite16ChunkedRead(b *testing.B) {
	var (
		r    = NewRingBuffer(16)
		from = make([]byte, 32)
		to   = make([]byte, 2)
	)

	for i := 0; i < b.N; i++ {
		r.Write(from)
		for k := 0; k < len(from); k += len(to) {
			r.Read(to)
		}
	}
}
