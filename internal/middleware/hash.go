package middleware

import (
	"bytes"
	"io"
	"net/http"

	"github.com/MaximPolyaev/go-metrics/internal/hash"
	"github.com/MaximPolyaev/go-metrics/internal/httpbufwritter"
)

// WithHashing - check hash key with request body for access process
func WithHashing(hashKey string) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			reqHash := r.Header.Get("HashSHA256")
			if reqHash == "" {
				next.ServeHTTP(w, r)
				return
			}

			reqBody, err := io.ReadAll(r.Body)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if err := r.Body.Close(); err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			r.Body = io.NopCloser(bytes.NewBuffer(reqBody))

			if reqHash != hash.Encode(reqBody, hashKey) {
				http.Error(w, "incorrect hash sign", http.StatusBadRequest)
				return
			}

			bw := httpbufwritter.New(w)

			next.ServeHTTP(bw, r)

			respBytes := bw.Bytes()
			w.Header().Set("HashSHA256", hash.Encode(respBytes, hashKey))
			w.WriteHeader(bw.StatusCode())

			if _, err := w.Write(respBytes); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})
	}
}
