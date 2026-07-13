package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ai-finops/ai-finops/packages/config"
	"github.com/ai-finops/ai-finops/packages/logger"
	svcconfig "github.com/ai-finops/ai-finops/services/outbox-relay/internal/config"
	"github.com/ai-finops/ai-finops/services/outbox-relay/internal/infrastructure/postgres"
	"github.com/ai-finops/ai-finops/services/outbox-relay/internal/relay"
)

func main() {
	cfg, err := config.Load[svcconfig.Config]("outbox-relay", "config/defaults.yaml", "config/config.local.yaml")
	if err != nil {
		slog.Error("failed to load config", slog.String("error", err.Error()))
		os.Exit(1)
	}

	log := logger.New(cfg.ServiceName, cfg.LogLevel)
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	pool, queries, err := postgres.Open(ctx, cfg.DatabaseURL)
	if err != nil {
		log.Error("failed to open postgres", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer pool.Close()

	rel, err := relay.New(pool, queries, cfg, log)
	if err != nil {
		log.Error("failed to create relay", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer func() { _ = rel.Close() }()

	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})
	healthSrv := &http.Server{
		Addr:    fmt.Sprintf(":%d", cfg.HTTPPort),
		Handler: mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	go func() {
		log.Info("health server starting", slog.String("addr", healthSrv.Addr))
		if err := healthSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Error("health server failed", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	go func() {
		log.Info("outbox relay starting")
		if err := rel.Run(ctx); err != nil && ctx.Err() == nil {
			log.Error("outbox relay stopped", slog.String("error", err.Error()))
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	log.Info("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = healthSrv.Shutdown(shutdownCtx)
	log.Info("outbox relay stopped")
}