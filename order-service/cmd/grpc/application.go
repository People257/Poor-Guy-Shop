package main

import (
	"context"

	"github.com/people257/poor-guy-shop/common/server"

	"github.com/people257/poor-guy-shop/order-service/api/cart"
	"github.com/people257/poor-guy-shop/order-service/api/order"
	pb_cart "github.com/people257/poor-guy-shop/order-service/gen/proto/order/cart"
	pb_order "github.com/people257/poor-guy-shop/order-service/gen/proto/order/order"
	"google.golang.org/grpc"
)

// Application 应用程序结构
type Application struct {
	Server *server.Server
}

// NewApplication 创建应用程序实例
func NewApplication(
	srv *server.Server,
	orderHandler *order.GrpcHandler,
	cartHandler *cart.GrpcHandler,
) *Application {
	// 注册gRPC服务
	srv.RegisterServer(func(grpcServer *grpc.Server) {
		pb_order.RegisterOrderServiceServer(grpcServer, orderHandler)
		pb_cart.RegisterCartServiceServer(grpcServer, cartHandler)
	})

	return &Application{
		Server: srv,
	}
}

// Run 运行应用程序
func (app *Application) Run(ctx context.Context) error {
	return app.Server.Run(ctx)
}
