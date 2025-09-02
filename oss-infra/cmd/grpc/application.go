package main

import (
	"context"

	"github.com/people257/poor-guy-shop/common/server"
	"github.com/people257/poor-guy-shop/oss-infra/api/file"
	"github.com/people257/poor-guy-shop/oss-infra/cmd/grpc/internal/config"
	filepb "github.com/people257/poor-guy-shop/oss-infra/gen/proto/oss/file"
	"google.golang.org/grpc"
)

type Application struct {
	Server *server.Server
}

func NewApplication(
	cfg *config.Config,
	s *server.Server,
	fileHandler *file.Handler,
) *Application {

	s.RegisterServer(func(s *grpc.Server) {
		filepb.RegisterFileServiceServer(s, fileHandler)
	})

	return &Application{
		Server: s,
	}
}

func (s *Application) Run(ctx context.Context) error {
	return s.Server.Run(ctx)
}
