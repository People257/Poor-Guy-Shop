//go:build wireinject
// +build wireinject

package main

import (
	"context"
	"github.com/google/wire"

	"github.com/people257/poor-guy-shop/user-service/cmd/grpc/internal"
	"github.com/people257/poor-guy-shop/user-service/cmd/grpc/internal/config"
	"github.com/people257/poor-guy-shop/user-service/internal/application"
	"github.com/people257/poor-guy-shop/user-service/internal/infra"
)

func InitializeApplication(ctx context.Context, configPath string) (*Application, func()) {
	panic(wire.Build(
		config.MustLoad,
		config.GetGrpcServerConfig,

		config.ConfigProviderSet,
		application.AppProviderSet,
		api.HandlerProviderSet,
		infra.InfraProviderSet,
		internal.InternalProviderSet,

		server.InitializeServer,

		NewApplication,
	))
}
