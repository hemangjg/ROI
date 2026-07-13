package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ai-finops/ai-finops/services/analytics-service/internal/di"
	"github.com/ai-finops/ai-finops/packages/shared"
)

func main() {
	app, cleanup, err := di.Initialize()
	if err != nil {
		slog.Error("failed to initialize service", slog.String("error", err.Error()))
		os.Exit(1)
	}
	defer cleanup()

	httpSrv := app.HTTPServer()
	errCh := make(chan error, 2)

	go func() {
		app.Log.Info("grpc server starting", slog.String("addr", app.GRPCServer.Addr()))
		errCh <- app.GRPCServer.Serve()
	}()

	go func() {
		app.Log.Info("http server starting", slog.String("addr", httpSrv.Addr))
		if err := httpSrv.ListenAndServe(); err != nil {
			errCh <- err
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		if err != nil && err.Error() != "http: Server closed" {
			app.Log.Error("server stopped with error", slog.String("error", err.Error()))
			os.Exit(1)
		}
	case sig := <-sigCh:
		app.Log.Info("shutdown signal received", slog.String("signal", sig.String()))
	}

	ctx, cancel := shared.ShutdownContext(15 * time.Second)
	defer cancel()

	app.GRPCServer.GracefulStop()

	if err := httpSrv.Shutdown(ctx); err != nil {
		app.Log.Error("http shutdown failed", slog.String("error", err.Error()))
		os.Exit(1)
	}

	app.Log.Info("server stopped gracefully")
}