package sirencast


import (
	"io"
	"fmt"
	"time"
	"testing"
	"bytes"
	"net"
	"net/http"
	"net/http/httptest"
)


func TestListenerAcceptance(t *testing.T) {
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello world")
	})
	server := httptest.NewUnstartedServer(testHandler)
	defer server.Close()

	listener := NewHTTPListener()
	server.Listener = listener

	server.Start()

	// Now forge a fake request with a SirenConn

	r, err := http.NewRequest("GET", server.URL, nil)

	if err != nil {
		t.Fatal(err)
	}

	rb := bytes.NewBuffer(nil)

	if err = r.Write(rb); err != nil {
		t.Fatal(err)
	}

	sc := &SirenConn{}

	conn := newFakeConn(t, rb)

	sc.conn = conn
	sc.reader = conn

	listener.Handler(sc)

	// Wait till fake connection is closed
	<-conn.finished

	buf := make([]byte, 8096)

	n, err := conn.out.Read(buf)

	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(buf[:n]))
}


func newFakeConn(t *testing.T, in io.Reader) *fakeConn {
	return &fakeConn{
		t: t,
		in: in,
		out: bytes.NewBuffer(nil),
		finished: make(chan bool, 15),
	}
}


type fakeConn struct {
	t 		*testing.T
	in 		io.Reader
	out     *bytes.Buffer
	finished chan bool
}


func (fc *fakeConn) Read(b []byte) (n int, err error) {
	return fc.in.Read(b)
}

func (fc *fakeConn) Write(b []byte) (n int, err error) {
	return fc.out.Write(b)
}

func (fc *fakeConn) Close() error {
	fc.finished <- true

	fc.t.Log("Closing fake connection")
	return nil
}

func (fc *fakeConn) LocalAddr() net.Addr {
	fc.t.Log("Local address lookup on fake connection")
	return &net.TCPAddr{
		IP: net.IPv4(255, 255, 255, 255),
		Port: 9999,
		Zone: "",
	}
}

func (fc *fakeConn) RemoteAddr() net.Addr {
	fc.t.Log("Remote address lookup on fake connection")
	return &net.TCPAddr{
		IP: net.IPv4(255, 255, 255, 255),
		Port: 9999,
		Zone: "",
	}
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
