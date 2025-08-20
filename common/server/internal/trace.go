package internal

import (
	"context"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/people257/poor-guy-shop/common/server/config"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.32.0"
	"go.opentelemetry.io/otel/trace"
	"go.opentelemetry.io/otel/trace/noop"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func logTraceID(ctx context.Context) logging.Fields {
	methodName, ok := grpc.Method(ctx)
	if ok && methodName == grpc_health_v1.Health_Check_FullMethodName {
		return nil
	}

	if span := trace.SpanContextFromContext(ctx); span.IsSampled() {
		return logging.Fields{"traceID", span.TraceID().String()}
	}
	return nil
}

func NewTraceExporter(ctx context.Context, cfg *config.ObservabilityConfig) (sdktrace.SpanExporter, func()) {
	if !cfg.Trace.Enable {
		return nil, func() {}
	}

	exporter, err := otlptracegrpc.New(ctx,
		otlptracegrpc.WithEndpoint(cfg.Address),
		otlptracegrpc.WithInsecure(),
		otlptracegrpc.WithCompressor("gzip"),
		otlptracegrpc.WithHeaders(cfg.OTLPHeaders),
	)
	if err != nil {
		panic(err)
	}
	cleanUp := func() {
		ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
		defer cancel()
		if err := exporter.Shutdown(ctx); err != nil {
			zap.L().Error("failed to shutdown trace exporter", zap.Error(err))
		}
	}
	return exporter, cleanUp
}

func NewSampler(cfg *config.ObservabilityConfig) sdktrace.Sampler {
	if cfg.Trace.Enable {
		return sdktrace.AlwaysSample()
	}
	return sdktrace.NeverSample()
}

func NewTracerProvider(
	sampler sdktrace.Sampler,
	exporter sdktrace.SpanExporter,
	serverCfg *config.ServerConfig,
	cfg *config.ObservabilityConfig,
) (trace.TracerProvider, func()) {
	cleanUp := func() {}

	if cfg.Trace.Enable {
		res := resource.NewSchemaless(
			semconv.ServiceName(serverCfg.Name),
			semconv.DeploymentEnvironmentName(serverCfg.Env),
		)

		p := sdktrace.NewTracerProvider(
			sdktrace.WithSampler(sampler),
			sdktrace.WithBatcher(exporter),
			sdktrace.WithResource(res),
		)

		otel.SetTracerProvider(p)
		otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(propagation.TraceContext{}, propagation.Baggage{}))

		cleanUp = func() {
			ctx, cancel := context.WithTimeout(context.Background(), cfg.Timeout)
			defer cancel()
			if err := p.Shutdown(ctx); err != nil {
				zap.L().Error("failed to shutdown trace provider", zap.Error(err))
			}
		}

		return p, cleanUp
	} else {
		return noop.NewTracerProvider(), cleanUp
	}
}
