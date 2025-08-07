package main

import (
	"context"
	"github.com/people257/poor-guy-shop/common/server"
	"google.golang.org/grpc"
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
