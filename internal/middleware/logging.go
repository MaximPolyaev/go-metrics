package middleware

import (
	"net/http"
	"time"

	"github.com/MaximPolyaev/go-metrics/internal/logger"
	"github.com/sirupsen/logrus"
)

type (
	Middleware func(next http.Handler) http.Handler

	responseData struct {
		status int
		size   int
	}

	loggingResponseWriter struct {
		http.ResponseWriter
		responseData *responseData
	}
)

func WithLogging(log *logger.Logger) Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			responseData := &responseData{
				status: http.StatusOK,
				size:   0,
			}

			lw := loggingResponseWriter{
				ResponseWriter: w,
				responseData:   responseData,
			}

			next.ServeHTTP(&lw, r)
			duration := time.Since(start)

			log.WithFields(logrus.Fields{
				"URI":      r.RequestURI,
				"method":   r.Method,
				"duration": duration,
				"status":   responseData.status,
				"size":     responseData.size,
			}).Info("")
		})
	}
}

func (r *loggingResponseWriter) Write(b []byte) (int, error) {
	size, err := r.ResponseWriter.Write(b)
	r.responseData.size += size

	return size, err
}

func (r *loggingResponseWriter) WriteHeader(statusCode int) {
	r.ResponseWriter.WriteHeader(statusCode)
	r.responseData.status = statusCode
}
