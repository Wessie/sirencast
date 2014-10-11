package icecast

import (
	"encoding/base64"
	"errors"
	"io"
	"net/http"
	"strconv"
	"strings"
)

var ErrInvalidHeader = errors.New("authorization: unable to parse header given")

// ParseDigest returns an `username` and `password` from the http.Request headers
// as specified by RFC1945 (HTTP Basic authentication). Returns a non-nil error
// if there was a problem parsing the header value.
func ParseDigest(r *http.Request) (user string, passwd string, err error) {
	authorization := strings.SplitN(r.Header.Get("Authorization"), " ", 2)

	if len(authorization) != 2 || authorization[0] != "Basic" {
		return "", "", ErrInvalidHeader
	}

	decoded, err := base64.StdEncoding.DecodeString(authorization[1])

	if err != nil {
		return "", "", ErrInvalidHeader
	}

	pair := strings.SplitN(string(decoded), ":", 2)

	if len(pair) != 2 {
		return "", "", ErrInvalidHeader
	}

	return pair[0], pair[1], nil
}

// WriteHeader writes the HTTP status line and headers to writer `w`.
func WriteHeader(w io.Writer, h http.Header, code int) error {
	text := http.StatusText(code)
	if text == "" {
		return errors.New("header: invalid status code")
	}
	statusCode := strconv.Itoa(code)
	io.WriteString(w, "HTTP/1.0 "+statusCode+" "+text+"\r\n")

	if h != nil {
		if err := h.Write(w); err != nil {
			return err
		}
	}

	_, err := io.WriteString(w, "\r\n")
	return err
}
