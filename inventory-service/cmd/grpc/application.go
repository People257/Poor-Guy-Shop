package main

import (
	"google.golang.org/grpc"

	"github.com/people257/poor-guy-shop/inventory-service/api/inventory"
	pb "github.com/people257/poor-guy-shop/inventory-service/gen/proto/proto/inventory"
)

// Application 应用程序
type Application struct {
	inventoryServer *inventory.Server
}

// NewApplication 创建应用程序
func NewApplication(inventoryServer *inventory.Server) *Application {
	return &Application{
		inventoryServer: inventoryServer,
	}
}

// RegisterServices 注册gRPC服务
func (a *Application) RegisterServices(s *grpc.Server) {
	pb.RegisterInventoryServiceServer(s, a.inventoryServer)
}
