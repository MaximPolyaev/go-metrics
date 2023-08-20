package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
)

type compressWriter struct {
	w  http.ResponseWriter
	zw *gzip.Writer
}

type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func GzipMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ow := w

		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		if supportsGzip {
			cw := newCompressWriter(w)
			ow = cw
			defer func() {
				if err := cw.Close(); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}
			}()
		}

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

		next.ServeHTTP(ow, r)
	})
}

func (c *compressWriter) Header() http.Header {
	return c.w.Header()
}

func (c *compressWriter) Write(p []byte) (int, error) {
	if c.isAllowGzip() {
		return c.zw.Write(p)
	}

	return c.w.Write(p)
}

func (c *compressWriter) WriteHeader(statusCode int) {
	if c.isAllowGzip() && statusCode < 300 {
		c.w.Header().Set("Content-Encoding", "gzip")
	}
	c.w.WriteHeader(statusCode)
}

func (c *compressWriter) Close() error {
	return c.zw.Close()
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

func (c *compressWriter) isAllowGzip() bool {
	ct := c.Header().Get("Content-Type")
	return ct == "text/html" || ct == "application/json"
}

func newCompressWriter(w http.ResponseWriter) *compressWriter {
	w.Header().Set("Content-Encoding", "gzip")

	return &compressWriter{
		w:  w,
		zw: gzip.NewWriter(w),
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
