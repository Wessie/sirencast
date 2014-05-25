package sirencast

import (
	"io"
	"net"
	"strings"
	"testing"
	"time"
)

// TestConnReaderSwap tests the read-swap of Conn when the first reader
// is exhausted.
func TestConnReaderSwap(t *testing.T) {
	startContent, connContent := "hello", "world"

	start := strings.NewReader(startContent)
	end := strings.NewReader(connContent)
	c := &fakeConn{Reader: end}

	conn := Conn{
		conn:    c,
		start:   start,
		reader:  start,
		handler: nil,
	}

	b := make([]byte, 5)

	_, err := conn.Read(b)
	if err != nil {
		t.Fatal("reading start failed:", err)
	}

	if string(b) != startContent {
		t.Fatal("reading start returned different content:", string(b), "!=", startContent)
	}

	_, err = conn.Read(b)
	if err != nil {
		t.Fatal("reading conn failed:", err)
	}

	if string(b) != connContent {
		t.Fatal("reading conn returned different content:", string(b), "!=", connContent)
	}

	return
}

type fakeConn struct {
	io.Reader
	io.Writer
	io.Closer

	RemoteAddr_ net.Addr
	LocalAddr_  net.Addr

	closer chan struct{}
	closed bool
}

func (fc *fakeConn) Wait() {
	<-fc.closer
}

func (fc *fakeConn) Close() error {
	if !fc.closed {
		close(fc.closer)
		fc.closed = true
	}
	return fc.Closer.Close()
}

func (fc *fakeConn) RemoteAddr() net.Addr {
	if fc.RemoteAddr_ == nil {
		return fakeAddr{}
	}
	return fc.RemoteAddr_
}

func (fc *fakeConn) LocalAddr() net.Addr {
	if fc.LocalAddr_ == nil {
		return fakeAddr{}
	}
	return fc.LocalAddr_
}

func (fc *fakeConn) SetDeadline(t time.Time) error {
	return nil
}

func (fc *fakeConn) SetReadDeadline(t time.Time) error {
	return nil
}

func (fc *fakeConn) SetWriteDeadline(t time.Time) error {
	return nil
}
