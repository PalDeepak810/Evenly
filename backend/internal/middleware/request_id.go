package middleware

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"net/http"
)

type contextKey string

const requestIDKey contextKey = "requestID"

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requestID, err := generateRequestID()

		if err != nil {
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		ctx := context.WithValue(r.Context(), requestIDKey, requestID)

		w.Header().Set("X-Request-Id", requestID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func generateRequestID() (string, error) {
	var id [16]byte

	if _, err := rand.Read(id[:]); err != nil {
		return "", fmt.Errorf("generate request ID: %w", err)
	}

	// Set UUID version 4.
	id[6] = (id[6] & 0x0f) | 0x40

	//Set UUID variant.
	id[8] = (id[8] & 0x3f) | 0x80

	return fmt.Sprintf(
		"%s-%s-%s-%s-%s",
		hex.EncodeToString(id[:4]),
		hex.EncodeToString(id[4:6]),
		hex.EncodeToString(id[6:8]),
		hex.EncodeToString(id[8:10]),
		hex.EncodeToString(id[10:16]),
	), nil
}

func RequestIDFromContext(ctx context.Context) (string, bool) {
	id, ok := ctx.Value(requestIDKey).(string)
	return id, ok
}
