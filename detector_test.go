package sirencast

import (
	"bytes"
	"io"
	"io/ioutil"
	"testing"
)

type TestSourceClient struct {
	conn io.Reader
}

func (c *TestSourceClient) Read(p []byte) (n int, err error) {
	return c.conn.Read(p)
}

func EchoDetector(input io.Reader) ConnHandler {
	return func(c *Conn) {}
}

func TestDetectorDetect(t *testing.T) {
	d := NewDetectors()
	d.Register(EchoDetector)

	input := ioutil.NopCloser(bytes.NewBuffer(nil))
	pk := NewPeeker(input)

	if source := d.Detect(pk); source == nil {
		t.Error("Expected input return, got nil")
	}
}
