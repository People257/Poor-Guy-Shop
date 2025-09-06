package client

import (
	"context"
	"log"
	"time"

	"github.com/people257/poor-guy-shop/order-service/cmd/grpc/config"
)

// Manager 客户端管理器
type Manager struct {
	UserClient      *UserServiceClient
	ProductClient   *ProductServiceClient
	PaymentClient   *PaymentServiceClient
	InventoryClient *InventoryServiceClient
}

// NewManager 创建客户端管理器
func NewManager(cfg *config.ServicesConfig) (*Manager, error) {
	manager := &Manager{}

	// 创建用户服务客户端
	userClient, err := NewUserServiceClient(&cfg.UserService)
	if err != nil {
		log.Printf("Failed to create user service client: %v", err)
		return nil, err
	}
	manager.UserClient = userClient

	// 创建商品服务客户端
	productClient, err := NewProductServiceClient(&cfg.ProductService)
	if err != nil {
		log.Printf("Failed to create product service client: %v", err)
		return nil, err
	}
	manager.ProductClient = productClient

	// 创建支付服务客户端
	paymentClient, err := NewPaymentServiceClient(&cfg.PaymentService)
	if err != nil {
		log.Printf("Failed to create payment service client: %v", err)
		return nil, err
	}
	manager.PaymentClient = paymentClient

	// 创建库存服务客户端
	inventoryClient, err := NewInventoryServiceClient(&cfg.InventoryService)
	if err != nil {
		log.Printf("Failed to create inventory service client: %v", err)
		return nil, err
	}
	manager.InventoryClient = inventoryClient

	return manager, nil
}

// Close 关闭所有客户端连接
func (m *Manager) Close() error {
	var errs []error

	if err := m.UserClient.Close(); err != nil {
		errs = append(errs, err)
	}

	if err := m.ProductClient.Close(); err != nil {
		errs = append(errs, err)
	}

	if err := m.PaymentClient.Close(); err != nil {
		errs = append(errs, err)
	}

	if err := m.InventoryClient.Close(); err != nil {
		errs = append(errs, err)
	}

	if len(errs) > 0 {
		log.Printf("Errors while closing clients: %v", errs)
		return errs[0] // 返回第一个错误
	}

	return nil
}

// HealthCheck 检查所有服务的健康状态
func (m *Manager) HealthCheck(ctx context.Context) map[string]bool {
	status := make(map[string]bool)

	// 创建超时上下文
	timeoutCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	// 检查用户服务
	status["user"] = m.checkUserService(timeoutCtx)

	// 检查商品服务
	status["product"] = m.checkProductService(timeoutCtx)

	// 检查支付服务
	status["payment"] = m.checkPaymentService(timeoutCtx)

	// 检查库存服务
	status["inventory"] = m.checkInventoryService(timeoutCtx)

	return status
}

func (m *Manager) checkUserService(ctx context.Context) bool {
	// TODO: 实现用户服务健康检查
	return true
}

func (m *Manager) checkProductService(ctx context.Context) bool {
	// TODO: 实现商品服务健康检查
	return true
}

func (m *Manager) checkPaymentService(ctx context.Context) bool {
	// 简单检查：尝试验证一个不存在的订单
	_, err := m.PaymentClient.VerifyPayment(ctx, "health-check")
	return err == nil
}

func (m *Manager) checkInventoryService(ctx context.Context) bool {
	// 简单检查：尝试查询一个不存在的订单库存
	_, err := m.InventoryClient.ReleaseInventory(ctx, "health-check")
	return err == nil
}
