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

			if closeErr := r.Body.Close(); closeErr != nil {
				http.Error(w, closeErr.Error(), http.StatusBadRequest)
				return
			}

			r.Body = io.NopCloser(bytes.NewBuffer(reqBody))

			encodedHash, err := hash.Encode(reqBody, hashKey)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			if reqHash != encodedHash {
				http.Error(w, "incorrect hash sign", http.StatusBadRequest)
				return
			}

			bw := httpbufwritter.New(w)

			next.ServeHTTP(bw, r)

			respBytes := bw.Bytes()

			encodedHash, err = hash.Encode(respBytes, hashKey)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			w.Header().Set("HashSHA256", encodedHash)
			w.WriteHeader(bw.StatusCode())

			if _, err := w.Write(respBytes); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
			}
		})
	}
}
