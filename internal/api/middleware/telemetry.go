package middleware

import (
	"net/http"

	"go.opentelemetry.io/otel/propagation"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

func Telemetry(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		propagator := otel.GetTextMapPropagator()
		extractedCtx := propagator.Extract(r.Context(), propagation.HeaderCarrier(r.Header))

		tracer := otel.Tracer("TemplateApp")
		ctx, span := tracer.Start(
			extractedCtx,
			r.URL.Path,
			trace.WithAttributes(
				attribute.String(string(RequestIDKeyValue), GetRequestID(r)),
			),
		)
		defer span.End()

		r = r.Clone(ctx)

		next.ServeHTTP(w, r)
	})
}
