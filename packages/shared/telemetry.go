package shared

import (
	"context"
	"fmt"
	"log/slog"
	"runtime/debug"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.26.0"
)

const defaultServiceVersion = "0.1.0"

// InitTelemetry configures OTLP trace and metric exporters. When endpoint is empty,
// telemetry is disabled and the returned shutdown function is a no-op.
func InitTelemetry(ctx context.Context, serviceName, environment, endpoint string, log *slog.Logger) (func(context.Context) error, error) {
	if endpoint == "" {
		if log != nil {
			log.Info("telemetry disabled", slog.String("service", serviceName))
		}
		return func(context.Context) error { return nil }, nil
	}

	res, err := resource.Merge(
		resource.Default(),
		resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceName(serviceName),
			semconv.ServiceVersion(serviceVersion()),
			semconv.DeploymentEnvironment(environment),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("create otel resource: %w", err)
	}

	grpcEndpoint := normalizeGRPCEndpoint(endpoint)

	traceExporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(grpcEndpoint),
		otlptracegrpc.WithInsecure(),
	)
	if err != nil {
		return nil, fmt.Errorf("create trace exporter: %w", err)
	}

	tracerProvider := trace.NewTracerProvider(
		trace.WithBatcher(traceExporter),
		trace.WithResource(res),
	)
	otel.SetTracerProvider(tracerProvider)

	metricExporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint(grpcEndpoint),
		otlpmetricgrpc.WithInsecure(),
	)
	if err != nil {
		_ = tracerProvider.Shutdown(ctx)
		return nil, fmt.Errorf("create metric exporter: %w", err)
	}

	meterProvider := metric.NewMeterProvider(
		metric.WithReader(metric.NewPeriodicReader(metricExporter)),
		metric.WithResource(res),
	)
	otel.SetMeterProvider(meterProvider)

	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	if log != nil {
		log.Info("telemetry enabled",
			slog.String("service", serviceName),
			slog.String("environment", environment),
			slog.String("endpoint", grpcEndpoint),
		)
	}

	return func(shutdownCtx context.Context) error {
		shutdownCtx, cancel := context.WithTimeout(shutdownCtx, 10*time.Second)
		defer cancel()

		var shutdownErr error
		if err := tracerProvider.Shutdown(shutdownCtx); err != nil {
			shutdownErr = err
		}
		if err := meterProvider.Shutdown(shutdownCtx); err != nil {
			if shutdownErr != nil {
				return fmt.Errorf("shutdown telemetry: trace=%v, metric=%w", shutdownErr, err)
			}
			return fmt.Errorf("shutdown meter provider: %w", err)
		}
		return shutdownErr
	}, nil
}

func normalizeGRPCEndpoint(endpoint string) string {
	endpoint = strings.TrimSpace(endpoint)
	endpoint = strings.TrimPrefix(endpoint, "https://")
	endpoint = strings.TrimPrefix(endpoint, "http://")
	return endpoint
}

func serviceVersion() string {
	if info, ok := debug.ReadBuildInfo(); ok {
		if info.Main.Version != "" && info.Main.Version != "(devel)" {
			return info.Main.Version
		}
	}
	return defaultServiceVersion
}