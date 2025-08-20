//go:build wireinject
// +build wireinject

package server

import (
	"context"

	"github.com/google/wire"
	"github.com/people257/poor-guy-shop/common/server/config"
	"github.com/people257/poor-guy-shop/common/server/internal"
)

func InitializeServer(ctx context.Context, cfg *config.GrpcServerConfig) (*Server, func()) {
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

		internal.NewConsulClient,
		internal.NewRegister,

		internal.NewObservabilityHttpServer,
		internal.NewGrpcServer,
		newServer,
	))
}
