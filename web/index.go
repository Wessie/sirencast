package web

import "net/http"

var Root = http.NewServeMux()

func init() {
	http.Handle("/", Root)
	Root.HandleFunc("/", func(rw http.ResponseWriter, r *http.Request) {
		rw.Write([]byte("hello world"))
	})
}
