//go:build wireinject

package main

import (
	"context"

	"github.com/people257/poor-guy-shop/user-service/api"
	"github.com/people257/poor-guy-shop/user-service/cmd/grpc/config"

	"github.com/people257/poor-guy-shop/user-service/cmd/grpc/internal"
	"github.com/people257/poor-guy-shop/user-service/internal/application"
	"github.com/people257/poor-guy-shop/user-service/internal/domain"

	"github.com/people257/poor-guy-shop/user-service/internal/infra"

	"github.com/google/wire"
	"github.com/people257/poor-guy-shop/common/server"
)

func InitializeApplication(ctx context.Context, configPath string) (*Application, func()) {
	panic(wire.Build(
		config.MustLoad,
		config.GetGrpcServerConfig,

		config.ConfigProviderSet,

		application.AppProviderSet,
		api.APIProviderSet,
		infra.InfraProviderSet,
		domain.DomainServiceProviderSet,
		internal.InternalProviderSet,

		server.InitializeServer,

		NewApplication,
	))
}
