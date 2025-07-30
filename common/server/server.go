package server

import (
	"context"
	"errors"
	"fmt"
	capi "github.com/hashicorp/consul/api"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"net"
	"net/http"
	"poor-guy-shop/common/server/config"
	"poor-guy-shop/common/server/internal"
	"time"
)

type Server struct {
	Config            *config.GrpcServerConfig
	GrpcServer        *grpc.Server
	metricsHttpServer *http.Server
	Register          *internal.Register
}

func newServer(
	cfg *config.GrpcServerConfig,
	s *grpc.Server,
	hs *http.Server,
	register *internal.Register,

	_ trace.TracerProvider,
	_ metric.MeterProvider,
	_ *capi.Client,
) *Server {
	healthcheck := health.NewServer()
	grpc_health_v1.RegisterHealthServer(s, healthcheck)

	return &Server{
		GrpcServer:        s,
		metricsHttpServer: hs,
		Register:          register,
		Config:            cfg,
	}
}

func (s *Server) startGrpcServer() error {
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

func (s *Server) startObservabilityHttp() error {
	isMetricsEnable := s.Config.Observability.Metrics.Enable
	isPprofEnable := s.Config.Observability.Pprof.Enable

	if !isMetricsEnable && !isPprofEnable {
		return nil
	}

	zap.L().Info("starting observability http server", zap.Uint16("port", s.Config.Observability.Port),
		zap.Bool("metrics", isMetricsEnable),
		zap.Bool("pprof", isPprofEnable))
	if err := s.metricsHttpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		return err
	}

	return nil
}

func (s *Server) registerToRegistry(ctx context.Context) error {
	err := s.Register.RegisterService()
	if err != nil {
		return fmt.Errorf("failed to register service: %w", err)
	}

	go func() {
		t := time.NewTicker(10 * time.Second)
		defer t.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-t.C:
				err := s.Register.CheckAndReRegisterService()
				if err != nil {
					zap.L().Error("failed to re-register service", zap.Error(err))
				}
			}
		}
	}()

	<-ctx.Done()

	if err := s.Register.DeregisterService(); err != nil {
		zap.L().Error("failed to deregister service", zap.Error(err))
	}

	return nil
}

func (s *Server) Run(ctx context.Context) error {
	errGroup, ctx := errgroup.WithContext(ctx)

	errGroup.Go(func() error {
		return s.startGrpcServer()
	})

	errGroup.Go(func() error {
		return s.startObservabilityHttp()
	})

	errGroup.Go(func() error {
		return s.registerToRegistry(ctx)
	})

	return errGroup.Wait()
}

func (s *Server) RegisterServer(fn func(s *grpc.Server)) {
	fn(s.GrpcServer)
}
