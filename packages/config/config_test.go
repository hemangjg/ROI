package config_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/ai-finops/ai-finops/packages/config"
)

func TestLoadFromEnv(t *testing.T) {
	t.Setenv("INGESTION_SERVICE_ENVIRONMENT", "local")
	t.Setenv("INGESTION_SERVICE_LOG_LEVEL", "info")
	t.Setenv("INGESTION_SERVICE_HTTP_PORT", "8081")

	dir := t.TempDir()
	defaults := filepath.Join(dir, "defaults.yaml")
	if err := os.WriteFile(defaults, []byte("environment: local\nlog_level: info\nhttp_port: 8080\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	cfg, err := config.Load[config.Base]("ingestion-service", defaults)
	if err != nil {
		t.Fatalf("load: %v", err)
	}

	if cfg.ServiceName != "ingestion-service" {
		t.Fatalf("service_name = %q", cfg.ServiceName)
	}
	if cfg.HTTPPort != 8081 {
		t.Fatalf("http_port = %d, want 8081", cfg.HTTPPort)
	}
}