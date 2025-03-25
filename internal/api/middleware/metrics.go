package middleware

import (
	"net/http"
	"strings"
	"time"

	"ugc/internal/metrics"
)

// Wrapper for http.ResponseWriter that preserves HTTP status code.
type statusRecorder struct {
	http.ResponseWriter
	code int
}

func (rec *statusRecorder) WriteHeader(code int) {
	rec.code = code
	rec.ResponseWriter.WriteHeader(code)
}

func Metrics(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := statusRecorder{w, 200}

		t := time.Now()

		next.ServeHTTP(&rec, r)

		var endpoint string
		parts := strings.Split(r.URL.Path, "/")
		if len(parts) > 0 {
			endpoint = parts[len(parts)-1]
		}

		if strings.Contains(r.URL.Path, "/api/v") {
			metrics.RequestsByEndpoint.WithLabelValues(endpoint).
				Observe(float64(time.Since(t).Milliseconds()))
			metrics.StatusCodes.WithLabelValues(http.StatusText(rec.code)).Inc()
		}
	})
}
