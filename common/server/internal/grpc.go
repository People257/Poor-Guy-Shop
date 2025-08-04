package internal

import (
	"buf.build/go/protovalidate"
	"github.com/people257/poor-guy-shop/common/server/config"
	"github.com/people257/poor-guy-shop/common/server/internal/interceptor"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc/filters"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc/status"
)

func NewGrpcServer(cfg *config.ServerConfig, logger *zap.Logger) (*grpc.Server, func()) {
	logOpts := []logging.Option{
		logging.WithLogOnEvents(logging.StartCall),
		logging.WithFieldsFromContext(logTraceID),
	}

	validator, err := protovalidate.New(protovalidate.WithFailFast())
	if err != nil {
		panic(err)
	}

	s := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler(
			otelgrpc.WithFilter(filters.Not(filters.HealthCheck())),
		)),
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(InterceptorLogger(logger), logOpts...),
			interceptor.ValidateUnaryServerInterceptor(validator),
			recovery.UnaryServerInterceptor(
				recovery.WithRecoveryHandler(panicHandler),
			),
		),
		grpc.Creds(insecure.NewCredentials()),
	)

	if cfg.Env == config.EnvDev {
		reflection.Register(s)
	}

	cleanUp := func() {
		s.GracefulStop()
	}

	return s, cleanUp
}

func panicHandler(p any) error {
	zap.L().Error("panic triggered", zap.Any("panic", p))
	return status.Error(codes.Internal, "Internal Server Error")
}
