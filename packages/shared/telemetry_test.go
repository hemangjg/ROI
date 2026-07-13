package shared

import (
	"context"
	"testing"
)

func TestInitTelemetryDisabled(t *testing.T) {
	shutdown, err := InitTelemetry(context.Background(), "test-service", "local", "", nil)
	if err != nil {
		t.Fatalf("InitTelemetry: %v", err)
	}
	if err := shutdown(context.Background()); err != nil {
		t.Fatalf("shutdown: %v", err)
	}
}

func TestNormalizeGRPCEndpoint(t *testing.T) {
	tests := []struct {
		in   string
		want string
	}{
		{"http://localhost:4317", "localhost:4317"},
		{"https://collector:4317", "collector:4317"},
		{"otel-collector:4317", "otel-collector:4317"},
		{"  http://otel-collector:4317  ", "otel-collector:4317"},
	}

	for _, tc := range tests {
		if got := normalizeGRPCEndpoint(tc.in); got != tc.want {
			t.Errorf("normalizeGRPCEndpoint(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}