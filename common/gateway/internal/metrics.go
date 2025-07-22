package internal

import (
	"context"
	"github.com/people257/poor-guy-shop/common/gateway/config"
	"go.opentelemetry.io/contrib/instrumentation/runtime"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlpmetric/otlpmetricgrpc"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/noop"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.32.0"
	"go.uber.org/zap"
	"os"
	"time"
)

// NewMeterProvider creates a new meter provider
func NewMeterProvider(cfg *config.ObservabilityConfig,
	serverCfg *config.ServerConfig,
	exporter sdkmetric.Exporter,
) (metric.MeterProvider, func()) {
	cleanUp := func() {}

	if cfg.Metrics.Enable {
		hostName, _ := os.Hostname()
		res := resource.NewSchemaless(
			semconv.ServiceName(serverCfg.Name),
			semconv.DeploymentEnvironmentName(serverCfg.Env),
			semconv.HostName(hostName),
			semconv.ProcessPID(os.Getpid()),
		)
		mp := sdkmetric.NewMeterProvider(
			sdkmetric.WithResource(res),
			sdkmetric.WithReader(sdkmetric.NewPeriodicReader(
				exporter,
				sdkmetric.WithProducer(runtime.NewProducer()),
				sdkmetric.WithInterval(30*time.Second),
			)),
		)
		otel.SetMeterProvider(mp)

		os.Setenv("OTEL_GO_X_DEPRECATED_RUNTIME_METRICS", "true")
		err := runtime.Start()
		if err != nil {
			zap.L().Error("failed to start runtime metrics collector", zap.Error(err))
		}

		cleanUp = func() {
			ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
			defer cancel()
			if err := mp.Shutdown(ctx); err != nil {
				zap.L().Error("failed to shutdown metrics provider", zap.Error(err))
			}
		}

		return mp, cleanUp
	} else {
		return noop.NewMeterProvider(), cleanUp
	}
}

func NewMetricExporter(ctx context.Context, cfg *config.ObservabilityConfig) sdkmetric.Exporter {
	if !cfg.Metrics.Enable {
		return nil
	}
	exporter, err := otlpmetricgrpc.New(ctx,
		otlpmetricgrpc.WithEndpoint(cfg.Address),
		otlpmetricgrpc.WithInsecure(),
		otlpmetricgrpc.WithCompressor("gzip"),
		otlpmetricgrpc.WithHeaders(cfg.OTLPHeaders),
	)
	if err != nil {
		panic(err)
	}
	return exporter
}
