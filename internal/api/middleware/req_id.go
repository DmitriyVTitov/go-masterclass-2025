package middleware

import (
	"context"
	"net/http"

	"github.com/google/uuid"
)

const RequestIDKeyValue RequestIDKeyType = "X-Request-ID"

type RequestIDKeyType string

func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var requestID string
		if r.Header.Get(string(RequestIDKeyValue)) != "" {
			requestID = r.Header.Get(string(RequestIDKeyValue))
		} else {
			requestID = uuid.New().String()
		}
		ctx := context.WithValue(r.Context(), RequestIDKeyValue, requestID)
		r = r.Clone(ctx)

		next.ServeHTTP(w, r)
	})
}

func GetRequestID(r *http.Request) string {
	reqID := r.Context().Value(RequestIDKeyValue)
	if s, ok := reqID.(string); ok {
		return s
	}

	return ""
}
