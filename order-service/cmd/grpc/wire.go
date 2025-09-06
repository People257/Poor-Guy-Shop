//go:build wireinject
// +build wireinject

package main

import (
	"context"

	"github.com/google/wire"

	"github.com/people257/poor-guy-shop/common/server"

	"github.com/people257/poor-guy-shop/order-service/api"
	appconfig "github.com/people257/poor-guy-shop/order-service/cmd/grpc/config"
	"github.com/people257/poor-guy-shop/order-service/cmd/grpc/internal"
	"github.com/people257/poor-guy-shop/order-service/internal/application"
	"github.com/people257/poor-guy-shop/order-service/internal/domain"
	"github.com/people257/poor-guy-shop/order-service/internal/infra"
)

// InitializeApplication 初始化应用程序
func InitializeApplication(ctx context.Context, configPath string) (*Application, func(), error) {
	panic(wire.Build(
		// 配置相关
		appconfig.MustLoad,
		appconfig.GetGrpcServerConfig,
		appconfig.GetDBConfig,
		appconfig.GetServicesConfig,

		// 基础设施
		internal.NewDatabase,
		internal.NewGormDB,
		internal.NewQuery,
		server.InitializeServer,

		// 各层Provider
		infra.ProviderSet,
		domain.ProviderSet,
		application.ProviderSet,
		api.ProviderSet,

		// 应用程序
		NewApplication,
	))
}
