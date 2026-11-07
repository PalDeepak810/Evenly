package middleware

import (
	"log/slog"
	"net/http"
	"time"
)

type responseWriter struct {
	http.ResponseWriter
	status      int
	wroteHeader bool
}

func (rw *responseWriter) WriteHeader(statusCode int) {
	if rw.wroteHeader {
		return
	}

	rw.status = statusCode
	rw.wroteHeader = true
	rw.ResponseWriter.WriteHeader(statusCode)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if !rw.wroteHeader {
		rw.WriteHeader(http.StatusOK)
	}

	return rw.ResponseWriter.Write(b)
}

func RequestLogger(logger *slog.Logger) func(handler http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			rw := &responseWriter{
				ResponseWriter: w,
				status:         http.StatusOK,
			}

			next.ServeHTTP(rw, r)

			duration := time.Since(start)

			requestID, _ := RequestIDFromContext(r.Context())

			logger.InfoContext(
				r.Context(),
				"HTTP request completed",
				"request_id", requestID,
				"method", r.Method,
				"path", r.URL.Path,
				"status", rw.status,
				"duration", duration.String(),
			)
		})
	}
}
