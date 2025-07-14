package internal

import (
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"github.com/people257/poor-guy-shop/common/server/config"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"

	"google.golang.org/grpc/status"
)

func NewGrpcServer(cfg *config.GrpcServerConfig, logger *zap.Logger) (*grpc.Server, func()) {
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			recovery.UnaryServerInterceptor(recovery.WithRecoveryHandler(panicHandler)),
		),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)

}

func panicHandler(p any) (err error) {
	zap.L().Error("gRPC server panicked", zap.Any("panic", p), zap.Stack("stack"))
	return status.Errorf(codes.Internal, "gRPC server panicked")
}
