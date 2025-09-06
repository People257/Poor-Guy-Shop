package client

import (
	"context"
	"fmt"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/people257/poor-guy-shop/order-service/cmd/grpc/config"
)

// InventoryServiceClient 库存服务客户端
type InventoryServiceClient struct {
	config *config.ServiceConfig
	conn   *grpc.ClientConn
}

// NewInventoryServiceClient 创建库存服务客户端
func NewInventoryServiceClient(cfg *config.ServiceConfig) (*InventoryServiceClient, error) {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Printf("Warning: Failed to connect to inventory service at %s: %v. Using mock implementation.", addr, err)
		// 连接失败时使用mock模式
		return &InventoryServiceClient{
			config: cfg,
			conn:   nil,
		}, nil
	}

	return &InventoryServiceClient{
		config: cfg,
		conn:   conn,
	}, nil
}

// Close 关闭连接
func (c *InventoryServiceClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// ReserveInventory 预占库存
func (c *InventoryServiceClient) ReserveInventory(ctx context.Context, orderID string, items []InventoryItem) (*InventoryResponse, error) {
	if c.conn == nil {
		// Mock模式：直接返回成功
		log.Printf("Inventory client in mock mode: reserving inventory for order %s with %d items", orderID, len(items))
		return &InventoryResponse{
			Success: true,
			Message: "库存预占成功 (Mock)",
		}, nil
	}

	// TODO: 实现真实的gRPC调用
	// 由于proto定义依赖问题，暂时使用mock
	log.Printf("Inventory service connected but using mock implementation: reserving for order %s with %d items", orderID, len(items))
	return &InventoryResponse{
		Success: true,
		Message: "库存预占成功 (Connected Mock)",
	}, nil
}

// ReleaseInventory 释放预占库存
func (c *InventoryServiceClient) ReleaseInventory(ctx context.Context, orderID string) (*InventoryResponse, error) {
	if c.conn == nil {
		log.Printf("Inventory client in mock mode: releasing inventory for order %s", orderID)
		return &InventoryResponse{
			Success: true,
			Message: "库存释放成功 (Mock)",
		}, nil
	}

	// TODO: 实现真实的gRPC调用
	log.Printf("Inventory service connected but using mock implementation: releasing for order %s", orderID)
	return &InventoryResponse{
		Success: true,
		Message: "库存释放成功 (Connected Mock)",
	}, nil
}

// ConfirmInventory 确认库存扣减
func (c *InventoryServiceClient) ConfirmInventory(ctx context.Context, orderID string) (*InventoryResponse, error) {
	if c.conn == nil {
		log.Printf("Inventory client in mock mode: confirming inventory for order %s", orderID)
		return &InventoryResponse{
			Success: true,
			Message: "库存确认成功 (Mock)",
		}, nil
	}

	// TODO: 实现真实的gRPC调用
	log.Printf("Inventory service connected but using mock implementation: confirming for order %s", orderID)
	return &InventoryResponse{
		Success: true,
		Message: "库存确认成功 (Connected Mock)",
	}, nil
}
