package sirencast

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestListenerAcceptance tests the HTTPListener accepting a non-special net.Conn
// that it was passed by calling (*HTTPListener).Handler
func TestListenerAcceptance(t *testing.T) {
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello world")
	})
	server := httptest.NewUnstartedServer(testHandler)
	defer server.Close()

	listener := NewHTTPListener("internal")
	server.Listener = listener

	server.Start()

	// Now forge a fake request with a Conn

	r, err := http.NewRequest("GET", server.URL, nil)
	if err != nil {
		t.Fatal(err)
	}

	in := bytes.NewBuffer(nil)
	if err = r.Write(in); err != nil {
		t.Fatal(err)
	}

	sc := &Conn{}
	out := new(bytes.Buffer)

	conn := &fakeConn{
		Reader: in,
		Writer: out,
		Closer: ioutil.NopCloser(nil),
		closer: make(chan struct{}),
	}

	sc.conn = conn
	sc.reader = conn

	listener.Handler(sc)

	// Wait till fake connection is closed
	conn.Wait()

	buf := make([]byte, 8096)

	n, err := out.Read(buf)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(buf[:n]))
}
