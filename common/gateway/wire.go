//go:build wireinject
// +build wireinject

package gateway

import (
	"context"
	"github.com/google/wire"
	"github.com/people257/poor-guy-shop/common/gateway/config"
	"github.com/people257/poor-guy-shop/common/gateway/internal"
)

func InitializeGateway(
	ctx context.Context,
	cfg *config.GatewayConfig,
) (*Gateway, func()) {
	panic(wire.Build(
		config.ConfigProviderSet,

		// Logger
		internal.NewLogExporter,
		internal.NewLoggerProvider,
		internal.NewZapLogger,

		// Tracing
		internal.NewSampler,
		internal.NewTraceExporter,
		internal.NewTracerProvider,

		// Metrics
		internal.NewMetricExporter,
		internal.NewMeterProvider,

		internal.NewGrpcClient,

		internal.NewGatewayMux,
		internal.NewEcho,
		internal.NewObservabilityHttpServer,
		newGateway,
	))
}
