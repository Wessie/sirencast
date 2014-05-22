package web

import "net/http"

var Root = http.NewServeMux()

func init() {
	http.Handle("/", Root)
}
