package client

import (
	"context"
	"sync"

	"github.com/people257/poor-guy-shop/inventory-service/cmd/grpc/config"
	"github.com/people257/poor-guy-shop/inventory-service/internal/infra/client/order"
	"github.com/people257/poor-guy-shop/inventory-service/internal/infra/client/product"
	"github.com/people257/poor-guy-shop/inventory-service/internal/infra/client/user"
)

// Manager 客户端管理器
type Manager struct {
	OrderClient   order.Client
	ProductClient product.Client
	UserClient    user.Client

	mu      sync.RWMutex
	closed  bool
	closers []func() error
}

// NewManager 创建客户端管理器
func NewManager(servicesConfig *config.ServicesConfig) (*Manager, error) {
	var closers []func() error

	// 创建订单服务客户端
	orderClient, err := order.NewGrpcClient(servicesConfig)
	if err != nil {
		return nil, err
	}
	if closer, ok := orderClient.(interface{ Close() error }); ok {
		closers = append(closers, closer.Close)
	}

	// 创建商品服务客户端
	productClient, err := product.NewGrpcClient(servicesConfig)
	if err != nil {
		// 清理已创建的客户端
		for _, closer := range closers {
			closer()
		}
		return nil, err
	}
	if closer, ok := productClient.(interface{ Close() error }); ok {
		closers = append(closers, closer.Close)
	}

	// 创建用户服务客户端
	userClient, err := user.NewGrpcClient(servicesConfig)
	if err != nil {
		// 清理已创建的客户端
		for _, closer := range closers {
			closer()
		}
		return nil, err
	}
	if closer, ok := userClient.(interface{ Close() error }); ok {
		closers = append(closers, closer.Close)
	}

	return &Manager{
		OrderClient:   orderClient,
		ProductClient: productClient,
		UserClient:    userClient,
		closers:       closers,
	}, nil
}

// Close 关闭所有客户端连接
func (m *Manager) Close() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.closed {
		return nil
	}

	var lastErr error
	for _, closer := range m.closers {
		if err := closer(); err != nil {
			lastErr = err
		}
	}

	m.closed = true
	return lastErr
}

// HealthCheck 健康检查
func (m *Manager) HealthCheck(ctx context.Context) map[string]error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	if m.closed {
		return map[string]error{
			"manager": context.Canceled,
		}
	}

	results := make(map[string]error)

	// 检查订单服务
	if _, err := m.OrderClient.GetOrderStatus(ctx, "health-check"); err != nil {
		results["order"] = err
	}

	// 检查商品服务
	if _, err := m.ProductClient.ValidateProducts(ctx, nil); err != nil {
		results["product"] = err
	}

	// 检查用户服务
	if _, err := m.UserClient.ValidateUser(ctx, [16]byte{}); err != nil {
		results["user"] = err
	}

	return results
}
