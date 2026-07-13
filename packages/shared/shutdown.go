package shared

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// ServeWithGracefulShutdown runs an HTTP server until SIGINT/SIGTERM.
func ServeWithGracefulShutdown(log *slog.Logger, srv *http.Server, timeout time.Duration) error {
	errCh := make(chan error, 1)
	go func() {
		log.Info("server starting", slog.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errCh <- err
		}
	}()

	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	select {
	case err := <-errCh:
		return err
	case sig := <-sigCh:
		log.Info("shutdown signal received", slog.String("signal", sig.String()))
	}

	ctx, cancel := ShutdownContext(timeout)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		return err
	}
	log.Info("server stopped gracefully")
	return nil
}

// ShutdownContext returns a timeout context for graceful shutdown.
func ShutdownContext(timeout time.Duration) (context.Context, context.CancelFunc) {
	return context.WithTimeout(context.Background(), timeout)
}