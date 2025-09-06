//go:build wireinject
// +build wireinject

package main

import (
	"context"

	"github.com/google/wire"
	"github.com/people257/poor-guy-shop/common/server"
	"github.com/people257/poor-guy-shop/payment-service/api"
	"github.com/people257/poor-guy-shop/payment-service/cmd/grpc/config"
	"github.com/people257/poor-guy-shop/payment-service/cmd/grpc/internal"
	"github.com/people257/poor-guy-shop/payment-service/internal/application"
	"github.com/people257/poor-guy-shop/payment-service/internal/infra"
)

// InitializeApplication 初始化应用程序
func InitializeApplication(ctx context.Context, configPath string) (*Application, func(), error) {
	wire.Build(
		// 配置
		config.MustLoad,
		config.GetGrpcServerConfig,
		config.GetDBConfig,
		config.GetAlipayConfig,

		// 基础设施
		internal.NewDatabase,
		internal.NewGormDB,
		internal.NewQuery,

		// 服务器
		server.InitializeServer,

		// 各层提供者
		infra.ProviderSet,
		application.ProviderSet,
		api.ProviderSet,

		// 应用程序
		NewApplication,
	)
	return nil, nil, nil
}
