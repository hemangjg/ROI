package shared

import (
	"fmt"
	"net/http"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
)

// OTelHTTP instruments HTTP handlers with OpenTelemetry spans.
// Health and readiness probes are excluded to avoid trace noise.
func OTelHTTP(serviceName string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return otelhttp.NewHandler(next, serviceName,
			otelhttp.WithFilter(func(r *http.Request) bool {
				switch r.URL.Path {
				case "/healthz", "/readyz":
					return false
				default:
					return true
				}
			}),
			otelhttp.WithSpanNameFormatter(func(_ string, r *http.Request) string {
				return fmt.Sprintf("%s %s", r.Method, r.URL.Path)
			}),
		)
	}
}