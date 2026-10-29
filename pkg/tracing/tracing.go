package tracing

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.43.0"
)

const _defaultShutdownTimeout = 5 * time.Second

type ShutdownFunc func(context.Context) error

type Config struct {
	Enabled     bool
	ServiceName string
	Version     string
	Prettier    bool
	SampleRate  float64
}

func New(ctx context.Context, cfg Config) (ShutdownFunc, error) {
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	if !cfg.Enabled {
		return func(context.Context) error { return nil }, nil
	}

	exporterOpts := []stdouttrace.Option{stdouttrace.WithWriter(os.Stdout)}
	if cfg.Prettier {
		exporterOpts = append(exporterOpts, stdouttrace.WithPrettyPrint())
	}

	exporter, err := stdouttrace.New(exporterOpts...)
	if err != nil {
		return nil, fmt.Errorf("tracing - New - stdouttrace.New: %w", err)
	}

	res, err := resource.New(
		ctx,
		resource.WithFromEnv(),
		resource.WithTelemetrySDK(),
		resource.WithHost(),
		resource.WithAttributes(
			semconv.ServiceName(cfg.ServiceName),
			semconv.ServiceVersion(cfg.Version),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("tracing - New - resource.New: %w", err)
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter),
		sdktrace.WithResource(res),
		sdktrace.WithSampler(sdktrace.ParentBased(sdktrace.TraceIDRatioBased(cfg.SampleRate))),
	)

	otel.SetTracerProvider(tp)

	shutdown := func(shutdownCtx context.Context) error {
		timeoutCtx, cancel := context.WithTimeout(shutdownCtx, _defaultShutdownTimeout)
		defer cancel()

		if err := tp.Shutdown(timeoutCtx); err != nil {
			return fmt.Errorf("tracing - shutdown - tp.Shutdown: %w", err)
		}

		return nil
	}

	return shutdown, nil
}
