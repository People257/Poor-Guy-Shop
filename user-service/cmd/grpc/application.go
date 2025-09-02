package main

import (
	"context"

	"github.com/people257/poor-guy-shop/common/server"
	"github.com/people257/poor-guy-shop/user-service/api/address"
	"github.com/people257/poor-guy-shop/user-service/api/auth"
	"github.com/people257/poor-guy-shop/user-service/api/info"
	"github.com/people257/poor-guy-shop/user-service/cmd/grpc/internal/config"
	addresspb "github.com/people257/poor-guy-shop/user-service/gen/proto/user/address"
	authpb "github.com/people257/poor-guy-shop/user-service/gen/proto/user/auth"
	infopb "github.com/people257/poor-guy-shop/user-service/gen/proto/user/info"
	"google.golang.org/grpc"
)

type Application struct {
	Server *server.Server
}

func NewApplication(
	cfg *config.Config,
	s *server.Server,
	authServer *auth.AuthServer,
	infoServer *info.InfoServer,
	addressServer *address.AddressServer,

) *Application {

	s.RegisterServer(func(s *grpc.Server) {
		// 注册gRPC服务
		authpb.RegisterAuthServiceServer(s, authServer)
		infopb.RegisterInfoServiceServer(s, infoServer)
		addresspb.RegisterAddressServiceServer(s, addressServer)
	})

	return &Application{
		Server: s,
	}
}

func (s *Application) Run(ctx context.Context) error {
	return s.Server.Run(ctx)
}
