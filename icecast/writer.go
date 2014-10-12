package icecast

import "io"

func NewMultiWriter() *MultiWriter {
	return &MultiWriter{
		pending: make(chan io.Writer, 5),
		w:       make([]io.Writer, 0),
	}
}

type MultiWriter struct {
	pending chan io.Writer
	// w are the writers we will be writing to
	w []io.Writer
}

func (mw *MultiWriter) Add(c io.Writer) {
	mw.pending <- c
}

func (mw *MultiWriter) Write(p []byte) (n int, err error) {
pending:
	for {
		select {
		case w := <-mw.pending:
			mw.w = append(mw.w, w)
		default:
			break pending
		}
	}

	for i := 0; i < len(mw.w); {
		n, err = mw.w[i].Write(p)
		if n != len(p) || err != nil {
			mw.w[i], mw.w[len(mw.w)-1] = mw.w[len(mw.w)-1], nil
			mw.w = mw.w[:len(mw.w)-1]
			continue
		}
		i++
	}

	return len(p), nil
}
