package server

import (
	"context"
	"fmt"
	"github.com/people257/poor-guy-shop/common/server/config"
	"github.com/people257/poor-guy-shop/common/server/internal"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
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

func newServer(
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

func (s *Sever) Run(ctx context.Context) error {
	errGroup, ctx := errgroup.WithContext(ctx)

	// 启动grpc 服务
	errGroup.Go(func() error {
		return s.startGrpcSever()
	})

	// 启动服务注册与心跳 goroutine
	errGroup.Go(func() error {
		if err := s.Register.RegisterService(); err != nil {
			return fmt.Errorf("failed to register service: %v", err)
		}

		// 上下文取消时,注销服务
		<-ctx.Done()
		return s.Register.DeregisterService()
	})

	zap.L().Info("Sever is running")
	return errGroup.Wait()
}

func (s *Sever) RegisterServer(fn func(s *grpc.Server)) {
	fn(s.GrpcServer)
}
