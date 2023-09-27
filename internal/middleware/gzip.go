package middleware

import (
	"compress/gzip"
	"io"
	"net/http"
	"strings"
	"sync"

	"github.com/MaximPolyaev/go-metrics/internal/httpbufwritter"
)

type compressReader struct {
	r  io.ReadCloser
	zr *gzip.Reader
}

func GzipMiddleware(next http.Handler) http.Handler {
	gzipPool := sync.Pool{New: func() any {
		return gzip.NewWriter(io.Discard)
	}}

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

		bw := httpbufwritter.New(w)

		next.ServeHTTP(bw, r)

		acceptEncoding := r.Header.Get("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")
		contentType := bw.Header().Get("Content-Type")

		if supportsGzip && (contentType == "text/html" || contentType == "application/json") {
			w.Header().Set("Content-Encoding", "gzip")
			w.WriteHeader(bw.StatusCode())

			gz := gzipPool.Get().(*gzip.Writer)
			defer func() {
				if err := gz.Close(); err != nil {
					http.Error(w, err.Error(), http.StatusInternalServerError)
				}

				gzipPool.Put(gz)
			}()

			gz.Reset(w)

			if _, err := gz.Write(bw.Bytes()); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
			return
		}

		w.WriteHeader(bw.StatusCode())
		if _, err := w.Write(bw.Bytes()); err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	})
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
