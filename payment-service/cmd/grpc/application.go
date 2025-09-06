package main

import (
	"context"

	"github.com/people257/poor-guy-shop/common/server"
	"github.com/people257/poor-guy-shop/payment-service/api/payment"
	pb "github.com/people257/poor-guy-shop/payment-service/gen/proto/proto/payment"
	"google.golang.org/grpc"
)

// Application 应用程序结构
type Application struct {
	Server *server.Server
}

// NewApplication 创建应用程序
func NewApplication(
	srv *server.Server,
	paymentHandler *payment.GrpcHandler,
) *Application {
	// 注册gRPC服务
	srv.RegisterServer(func(grpcServer *grpc.Server) {
		pb.RegisterPaymentServiceServer(grpcServer, paymentHandler)
	})

	return &Application{
		Server: srv,
	}
}

// Run 运行应用程序
func (app *Application) Run(ctx context.Context) error {
	return app.Server.Run(ctx)
}
