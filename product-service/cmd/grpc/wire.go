//go:build wireinject
// +build wireinject

package main

import (
	"context"

	"github.com/google/wire"

	"github.com/people257/poor-guy-shop/common/server"
	"github.com/people257/poor-guy-shop/product-service/api"
	"github.com/people257/poor-guy-shop/product-service/cmd/grpc/config"
	"github.com/people257/poor-guy-shop/product-service/cmd/grpc/internal"
	"github.com/people257/poor-guy-shop/product-service/internal/application"
	"github.com/people257/poor-guy-shop/product-service/internal/domain"
	"github.com/people257/poor-guy-shop/product-service/internal/infra"
)

// InitializeApplication 初始化应用程序
func InitializeApplication(ctx context.Context, configPath string) (*Application, func()) {
	wire.Build(
		// 配置提供者
		config.ConfigProviderSet,

		// 内部依赖
		internal.InternalProviderSet,

		// 基础设施层
		infra.ProviderSet,

		// 领域层
		domain.ProviderSet,

		// 应用层
		application.ProviderSet,

		// API层
		api.ProviderSet,

		// 服务器
		server.InitializeServer,

		// 应用程序
		NewApplication,
	)
	return nil, nil
}
