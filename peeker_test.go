package sirencast

import (
	"bytes"
	"io"
	"reflect"
	"testing"
)

var testData = []byte("abcdefghijklmnopqrstuvwxyz")

// TestPeekerReset tests the Reset method on the Peeker type.
//
// This does 50 loops and reads 3 bytes each loop. These 3 bytes
// should equal testData[0:3] on every loop.
func TestPeekerReset(t *testing.T) {
	peek := NewPeeker(bytes.NewBuffer(testData))
	buf := make([]byte, 3)
	expect := string(testData[0:3])

	for i := 0; i < 50; i++ {
		_, err := peek.Read(buf)

		if err != nil {
			t.Fatalf("Error: %s", err)
		}

		if string(buf) != expect {
			t.Fatalf("Unequal: %v != %v", buf, expect)
		}
		peek.Reset()
	}
}

// TestPeekerStop tests that a Peeker actually stops reading after
// (Peeker).Stop was called.
func TestPeekerStop(t *testing.T) {
	peek := NewPeeker(bytes.NewBuffer(testData[:6]))

	buf := make([]byte, 6)

	n, err := peek.Read(buf)
	if err != nil {
		t.Error(err)
	}

	if !reflect.DeepEqual(buf[:n], testData[:n]) {
		t.Errorf("peeker returned invalid data (1): %v != %v", buf[:n], testData[:n])
	}

	peek.Reset()
	peek.Stop()

	n, err = peek.Read(buf)
	if err != nil {
		t.Errorf("peeker returned unexpected error: %v", err)
	} else if string(buf) != string(testData[:n]) {
		t.Errorf("peeker returned invalid data (2): %v != %v", string(buf), string(testData[:n]))
	}

	n, err = peek.Read(buf)
	if err != io.EOF {
		t.Logf("peeker dump: %+v", peek.(*PeekReader))
		t.Errorf("peeker did not stop reading after Stop: %v (%d %v)", buf, n, err)
	}
}

func BenchmarkPeekerBuffer(b *testing.B) {
	peeker := NewPeeker(bytes.NewBuffer(testData))
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
		peeker := NewPeeker(bytes.NewBuffer(testData))

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
