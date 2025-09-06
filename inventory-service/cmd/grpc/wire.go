//go:build wireinject

package main

import (
	"github.com/google/wire"

	"github.com/people257/poor-guy-shop/inventory-service/api"
	"github.com/people257/poor-guy-shop/inventory-service/cmd/grpc/config"
	"github.com/people257/poor-guy-shop/inventory-service/cmd/grpc/internal"
	"github.com/people257/poor-guy-shop/inventory-service/internal/application"
	"github.com/people257/poor-guy-shop/inventory-service/internal/infra"
)

// InitializeApplication 初始化应用程序
func InitializeApplication(configPath string) (*Application, error) {
	wire.Build(
		// Config
		config.ProviderSet,

		// Database
		internal.NewDatabase,
		internal.NewGormDB,
		internal.NewQuery,

		// Infrastructure
		infra.ProviderSet,

		// Application
		application.ProviderSet,

		// API
		api.ProviderSet,

		// Application
		NewApplication,
	)
	return nil, nil
}
