package order

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/people257/poor-guy-shop/inventory-service/cmd/grpc/config"
)

// Client 订单服务客户端接口
type Client interface {
	// GetOrderStatus 获取订单状态
	GetOrderStatus(ctx context.Context, orderID string) (string, error)

	// NotifyInventoryReserved 通知库存已预占
	NotifyInventoryReserved(ctx context.Context, orderID string, success bool, message string) error

	// NotifyInventoryConfirmed 通知库存已确认扣减
	NotifyInventoryConfirmed(ctx context.Context, orderID string, success bool, message string) error
}

// GrpcClient 订单服务gRPC客户端实现
type GrpcClient struct {
	conn   *grpc.ClientConn
	config *config.ServiceConfig
}

// NewGrpcClient 创建订单服务gRPC客户端
func NewGrpcClient(servicesConfig *config.ServicesConfig) (Client, error) {
	serviceConfig := servicesConfig.OrderService

	// 建立gRPC连接
	addr := fmt.Sprintf("%s:%d", serviceConfig.Host, serviceConfig.Port)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to order service: %w", err)
	}

	return &GrpcClient{
		conn:   conn,
		config: &serviceConfig,
	}, nil
}

// Close 关闭连接
func (c *GrpcClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// GetOrderStatus 获取订单状态
func (c *GrpcClient) GetOrderStatus(ctx context.Context, orderID string) (string, error) {
	// TODO: 实现订单状态查询
	// 这里需要根据实际的订单服务proto定义来实现
	return "pending", nil
}

// NotifyInventoryReserved 通知库存已预占
func (c *GrpcClient) NotifyInventoryReserved(ctx context.Context, orderID string, success bool, message string) error {
	// TODO: 实现库存预占通知
	// 这里需要根据实际的订单服务proto定义来实现
	return nil
}

// NotifyInventoryConfirmed 通知库存已确认扣减
func (c *GrpcClient) NotifyInventoryConfirmed(ctx context.Context, orderID string, success bool, message string) error {
	// TODO: 实现库存确认通知
	// 这里需要根据实际的订单服务proto定义来实现
	return nil
}
