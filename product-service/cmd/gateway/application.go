package main

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/people257/poor-guy-shop/common/auth"
	"github.com/people257/poor-guy-shop/common/gateway"
	addresspb "github.com/people257/poor-guy-shop/user-service/gen/proto/user/address"
	authpb "github.com/people257/poor-guy-shop/user-service/gen/proto/user/auth"
	infopb "github.com/people257/poor-guy-shop/user-service/gen/proto/user/info"
	"google.golang.org/grpc"
)

type Application struct {
	Gateway *gateway.Gateway
}

func NewApplication(
	gw *gateway.Gateway,
	authClient authpb.AuthServiceClient,
) *Application {
	_ = gw.RegisterHandler(func(gwmux *runtime.ServeMux, conn *grpc.ClientConn) error {
		var err error
		err = authpb.RegisterAuthServiceHandler(context.Background(), gwmux, conn)
		if err != nil {
			return err
		}
		err = infopb.RegisterInfoServiceHandler(context.Background(), gwmux, conn)
		if err != nil {
			return err
		}
		err = addresspb.RegisterAddressServiceHandler(context.Background(), gwmux, conn)
		if err != nil {
			return err
		}
		return err
	})

	e := gw.Echo
	e.Use(auth.BuildMetadataMiddleware(authClient))

	return &Application{
		Gateway: gw,
	}
}

func (s *Application) Run(ctx context.Context) error {
	return s.Gateway.Run(ctx)
}
