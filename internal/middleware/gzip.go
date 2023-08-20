package middleware

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type bufWriter struct {
	w          http.ResponseWriter
	buf        *bytes.Buffer
	statusCode int
}

type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentEncoding := r.Header.Get("Content-Encoding")
		sendsGzip := strings.Contains(contentEncoding, "gzip")
		if sendsGzip {
			cr, err := newCompressReader(r.Body)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
			r.Body = cr
			defer func() {
				if err := cr.Close(); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}()
		}

		bw := newBufWriter(w)

		next.ServeHTTP(bw, r)

		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		contentType := bw.Header().Get("Content-Type")
		if supportsGzip && (contentType == "text/html" || contentType == "application/json") {
			w.Header().Set("Content-Encoding", "gzip")
			w.WriteHeader(bw.statusCode)
			gz := gzip.NewWriter(w)

			if _, err := gz.Write(bw.buf.Bytes()); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}

			defer func() {
				if err := gz.Close(); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}()
			return
		}

		w.WriteHeader(bw.statusCode)
		if _, err := w.Write(bw.buf.Bytes()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
}

func (c *bufWriter) Header() http.Header {
	return c.w.Header()
}

func (c *bufWriter) Write(p []byte) (int, error) {
	return c.buf.Write(p)
}

func (c *bufWriter) WriteHeader(statusCode int) {
	c.statusCode = statusCode
}

func (c *compressReader) Read(p []byte) (n int, err error) {
	return c.zr.Read(p)
}

func (c *compressReader) Close() error {
	if err := c.r.Close(); err != nil {
		return err
	}
	return c.zr.Close()
}
func newBufWriter(w http.ResponseWriter) *bufWriter {
	var buf []byte

	return &bufWriter{
		w:          w,
		buf:        bytes.NewBuffer(buf),
		statusCode: http.StatusOK,
	}
}

func newCompressReader(r io.ReadCloser) (*compressReader, error) {
	zr, err := gzip.NewReader(r)
	if err != nil {
		return nil, err
	}

	return &compressReader{
		r:  r,
		zr: zr,
	}, nil
}
