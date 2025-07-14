//go:build wireinject

package server

import (
	"context"
	"github.com/google/wire"
	"github.com/people257/poor-guy-shop/common/server/config"
	"github.com/people257/poor-guy-shop/common/server/internal"
)

func InitializeServer(ctx context.Context, cfg *config.GrpcServerConfig) (*Sever, func()) {
	panic(wire.Build(

		config.ConfigProviderSet,

		internal.NewConsulClient,
		internal.NewRegister,

		internal.NewZapLogger,
		internal.NewGrpcServer,

		newServer,
	))
}
