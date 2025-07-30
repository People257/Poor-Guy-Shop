package client

import (
	"fmt"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"poor-guy-shop/common/resolver"
)

func newGrpcClient(cfg *Config, service string) (*grpc.ClientConn, func()) {
	conn, err := grpc.NewClient(
		fmt.Sprintf("consul://%s/%s", cfg.RegistryAddr(), service),
		grpc.WithResolvers(&resolver.Builder{}),
		grpc.WithStatsHandler(otelgrpc.NewClientHandler()),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithDefaultServiceConfig(`{"loadBalancingPolicy": "round_robin"}`),
	)
	if err != nil {
		panic(err)
	}
	cleanUp := func() {
		if err := conn.Close(); err != nil {
			zap.L().Error("failed to close grpc conn", zap.Error(err), zap.String("service", service))
		}
	}
	return conn, cleanUp
}

func NewGrpcClient[T any](cfg *Config,
	service string,
	newFunc func(grpc.ClientConnInterface) T,
) (T, func()) {
	client, cleanUp := newGrpcClient(cfg, service)
	return newFunc(client), cleanUp
}

func NewGrpcClientFromConn[T any](conn *grpc.ClientConn,
	newFunc func(grpc.ClientConnInterface) T,
) (T, func()) {
	return newFunc(conn), func() {}
}

func NewGrpcConn(cfg *Config, service string) (*grpc.ClientConn, func()) {
	return newGrpcClient(cfg, service)
}
