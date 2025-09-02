//go:build wireinject
// +build wireinject

package main

import (
	"context"

	"github.com/google/wire"
	"github.com/people257/poor-guy-shop/common/server"
	"github.com/people257/poor-guy-shop/oss-infra/api"
	"github.com/people257/poor-guy-shop/oss-infra/cmd/grpc/internal"
	"github.com/people257/poor-guy-shop/oss-infra/cmd/grpc/internal/config"
	"github.com/people257/poor-guy-shop/oss-infra/internal/application"
	"github.com/people257/poor-guy-shop/oss-infra/internal/domain"
	"github.com/people257/poor-guy-shop/oss-infra/internal/infra"
)

func InitializeApplication(ctx context.Context, configPath string) (*Application, func()) {
	panic(wire.Build(
		config.MustLoad,
		config.GetGrpcServerConfig,

		config.ConfigProviderSet,
		internal.InternalProviderSet,
		domain.DomainServiceProviderSet,
		infra.InfraProviderSet,
		application.AppProviderSet,
		api.HandlerProviderSet,

		server.InitializeServer,

		NewApplication,
	))
}
