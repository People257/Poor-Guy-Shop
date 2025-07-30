package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/people257/poor-guy-shop/common/gateway"
	"google.golang.org/grpc"
)

type Application struct {
	Gateway *gateway.Gateway
}

func NewApplication(
	gw *gateway.Gateway,
) *Application {
	_ = gw.RegisterHandler(func(gwmux *runtime.ServeMux, conn *grpc.ClientConn) error {
		var err error

		return err
	})

	return &Application{
		Gateway: gw,
	}
}

func (s *Application) Run(ctx context.Context) error {
	return s.Gateway.Run(ctx)
}
