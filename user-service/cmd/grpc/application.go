package main

import (
	"context"
	"google.golang.org/grpc"
	"poor-guy-shop/common/server"
	"poor-guy-shop/user-service/cmd/grpc/internal/config"
)

type Application struct {
	Server *server.Server
}

func NewApplication(
	cfg *config.Config,
	s *server.Server,

) *Application {

	s.RegisterServer(func(s *grpc.Server) {

	})

	return &Application{
		Server: s,
	}
}

func (s *Application) Run(ctx context.Context) error {
	return s.Server.Run(ctx)
}
