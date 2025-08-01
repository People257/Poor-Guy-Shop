//go:build wireinject

package main

import (
	"context"
	"github.com/people257/poor-guy-shop/common/gateway"
	"github.com/people257/poor-guy-shop/project-template/cmd/gateway/internal/config"

	"github.com/google/wire"
)

func InitializeApplication(ctx context.Context, configPath string) (*Application, func()) {
	panic(wire.Build(
		config.MustLoad,
		config.GetGatewayConfig,

		gateway.InitializeGateway,

		NewApplication,
	))
}
