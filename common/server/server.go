package server

import (
	"fmt"
	"github.com/people257/poor-guy-shop/common/server/config"
	"github.com/people257/poor-guy-shop/common/server/internal"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
)

type Sever struct {
	Config     *config.GrpcServerConfig
	GrpcServer *grpc.Server
	Register   *internal.Register
}

func newSever(
	cfg *config.GrpcServerConfig,
	s *grpc.Server,
	register *internal.Register,
) *Sever {
	healthcheck := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, healthcheck)

	return &Sever{
		Config:     cfg,
		GrpcServer: s,
		Register:   register,
	}
}

func (s *Sever) startGrpcSever() error {
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.Config.Server.Port))
	if err != nil {
		return err
	}
	zap.L().Info("starting grpc server", zap.Uint16("port", s.Config.Server.Port))
	if err := s.GrpcServer.Serve(lis); err != nil {
		return err
	}

	return nil
}
