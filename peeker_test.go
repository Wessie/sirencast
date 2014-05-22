package sirencast

import (
	"testing"
	"bytes"
	"reflect"
	"io"
)

var test_data = []byte("abcdefghijklmnopqrstuvwxyz")

func TestPeekerReset(t *testing.T) {
	peek := NewPeeker(bytes.NewBuffer(test_data))

	buf := make([]byte, 3)

	for i := 0; i < 50; i++ {
		_, err := peek.Read(buf)

		if err != nil {
			t.Fatalf("Error: %s", err)
		}

		if string(buf) != "abc" {
			t.Fatalf("Unequal: %v != %v", buf, "abc")
		}
		peek.Reset()
	}
}


func TestPeekerStop(t *testing.T) {
	t.Skip()

	peek := NewPeeker(bytes.NewBuffer(test_data))

	buf := make([]byte, 6)

	n, err := peek.Read(buf)

	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(buf[:n], test_data[:n]) {
		t.Errorf("Data returned is not equal to input: %v != %v", buf[:n], test_data[:n])
	}

	peek.Reset()
	peek.Stop()

	n, err = peek.Read(buf)

	if err != nil {
		t.Errorf("Peeker returned error it should not have %s", err)
	}

	n, err = peek.Read(buf)

	if err != io.EOF {
		t.Errorf("Peeker did not stop reading after Stop: %v", buf)
	}
}

func BenchmarkPeekerBuffer(b *testing.B) {
	peeker := NewPeeker(bytes.NewBuffer(test_data))
	buf := make([]byte, 6)

	peeker.Read(buf)

	for i := 0; i < b.N; i++ {
		peeker.Reset()
		_, err := peeker.Read(buf)

		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkPeekerWorst(b *testing.B) {
	buf := make([]byte, 3)

	var oldpeeker = NewPeeker(nil).(*PeekReader)

	for i := 0; i < b.N; i++ {
		peeker := NewPeeker(bytes.NewBuffer(test_data))

		// Steal the buffer from the previous peeker, otherwise
		// we allocate a lot of memory that isn't realistic in
		// the benchmark
		p := peeker.(*PeekReader)
		p.buffer = oldpeeker.buffer
		oldpeeker = p

		for n := 1; n < len(buf); n++ {
			var err error
			for tn, rn := 0, 0; tn < n; tn += rn {
				rn, err = peeker.Read(buf[:n-tn])

				if err != nil {
					b.Fatal(err)
				}
			}
		}
	}
}
