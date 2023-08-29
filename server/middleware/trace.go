package middleware

import (
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/goodleby/golang-server/tracing"
	"go.opentelemetry.io/otel/attribute"
)

func Trace(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ctx, span := tracing.Span(r.Context(), "Root")

		next.ServeHTTP(w, r.WithContext(ctx))

		routeID := fmt.Sprintf("%s %s", r.Method, chi.RouteContext(ctx).RoutePattern())
		span.SetName(routeID)
		span.SetAttributes(attribute.String("RequestURI", r.RequestURI))

		span.End()
	})
}