// Package httpbufwritter используется для буферизации HTTP запросов
package httpbufwritter

import (
	"bytes"
	"net/http"
)

type BufWriter struct {
	w          http.ResponseWriter
	buf        *bytes.Buffer
	statusCode int
}

func New(w http.ResponseWriter) *BufWriter {
	var buf []byte

	return &BufWriter{
		w:          w,
		buf:        bytes.NewBuffer(buf),
		statusCode: http.StatusOK,
	}
}

func (b *BufWriter) StatusCode() int {
	return b.statusCode
}

func (b *BufWriter) Bytes() []byte {
	return b.buf.Bytes()
}

func (b *BufWriter) Header() http.Header {
	return b.w.Header()
}

func (b *BufWriter) Write(p []byte) (int, error) {
	return b.buf.Write(p)
}

func (b *BufWriter) WriteHeader(statusCode int) {
	b.statusCode = statusCode
}
