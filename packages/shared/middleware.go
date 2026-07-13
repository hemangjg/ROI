package shared

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"time"

	"github.com/ai-finops/ai-finops/packages/logger"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/google/uuid"
	"go.opentelemetry.io/otel/trace"
)

type ctxKey string

const requestIDKey ctxKey = "request_id"

// RequestIDKey is the context key for request IDs.
const RequestIDHeader = "X-Request-ID"

// RequestID middleware assigns a unique request ID to each request.
func RequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := r.Header.Get(RequestIDHeader)
		if id == "" {
			id = uuid.NewString()
		}
		ctx := context.WithValue(r.Context(), requestIDKey, id)
		w.Header().Set(RequestIDHeader, id)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// RequestIDFromContext extracts the request ID from context.
func RequestIDFromContext(ctx context.Context) string {
	if v, ok := ctx.Value(requestIDKey).(string); ok {
		return v
	}
	return ""
}

// ErrorHandler converts AppError and unknown errors into JSON responses.
func ErrorHandler(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			rw := &errorResponseWriter{ResponseWriter: w, status: http.StatusOK}
			next.ServeHTTP(rw, r)

			if rw.err != nil {
				status := http.StatusInternalServerError
				code := "internal_error"
				message := "internal server error"

				if appErr, ok := rw.err.(*AppError); ok {
					status = appErr.Status
					code = appErr.Code
					message = appErr.Message
				}

				reqLog := logger.WithRequestID(log, RequestIDFromContext(r.Context()))
				reqLog.Error("request failed", slog.Int("status", status), slog.String("error", rw.err.Error()))

				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(status)
				_ = json.NewEncoder(w).Encode(map[string]string{
					"code":    code,
					"message": message,
				})
			}
		})
	}
}

// Recoverer logs panics and returns 500.
func Recoverer(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if rec := recover(); rec != nil {
					log.Error("panic recovered", slog.Any("panic", rec), slog.String("request_id", RequestIDFromContext(r.Context())))
					w.WriteHeader(http.StatusInternalServerError)
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}

// RequestLogger logs each HTTP request with structured slog fields.
func RequestLogger(log *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			next.ServeHTTP(ww, r)

			reqLog := logger.WithRequestID(log, RequestIDFromContext(r.Context()))
			if span := trace.SpanFromContext(r.Context()); span.SpanContext().IsValid() {
				reqLog = reqLog.With(slog.String("trace_id", span.SpanContext().TraceID().String()))
			}
			reqLog.Info("request completed",
				slog.String("method", r.Method),
				slog.String("path", r.URL.Path),
				slog.Int("status", ww.Status()),
				slog.Int("bytes", ww.BytesWritten()),
				slog.Duration("duration", time.Since(start)),
			)
		})
	}
}

type errorResponseWriter struct {
	http.ResponseWriter
	status int
	err    error
}

func (w *errorResponseWriter) WriteHeader(status int) {
	w.status = status
	w.ResponseWriter.WriteHeader(status)
}