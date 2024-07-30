package middleware

import (
	"bytes"
	"io"
	"net/http"
)

type Decoder interface {
	Decode(data []byte) ([]byte, error)
}

// WithDecrypt middleware для расшифровки сообщений
func WithDecrypt(decoder Decoder, decryptURI string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if decoder == nil || r.RequestURI != decryptURI {
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

			reqBody, err = decoder.Decode(reqBody)
			if err != nil {
				http.Error(w, err.Error(), http.StatusBadRequest)
				return
			}

			r.Body = io.NopCloser(bytes.NewBuffer(reqBody))

			next.ServeHTTP(w, r)
		})
	}
}
