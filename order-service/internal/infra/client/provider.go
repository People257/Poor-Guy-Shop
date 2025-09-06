package client

import (
	"github.com/google/wire"
	"github.com/people257/poor-guy-shop/order-service/cmd/grpc/config"
)

// ClientProviderSet Wire提供器集合
var ClientProviderSet = wire.NewSet(
	NewUserServiceClientFromConfig,
	NewProductServiceClientFromConfig,
	NewPaymentServiceClientFromConfig,
	NewInventoryServiceClientFromConfig,
	NewManagerFromConfig,
)

// NewUserServiceClientFromConfig 从配置创建用户服务客户端
func NewUserServiceClientFromConfig(cfg *config.ServicesConfig) (*UserServiceClient, error) {
	return NewUserServiceClient(&cfg.UserService)
}

// NewProductServiceClientFromConfig 从配置创建产品服务客户端
func NewProductServiceClientFromConfig(cfg *config.ServicesConfig) (*ProductServiceClient, error) {
	return NewProductServiceClient(&cfg.ProductService)
}

// NewPaymentServiceClientFromConfig 从配置创建支付服务客户端
func NewPaymentServiceClientFromConfig(cfg *config.ServicesConfig) (*PaymentServiceClient, error) {
	return NewPaymentServiceClient(&cfg.PaymentService)
}

// NewInventoryServiceClientFromConfig 从配置创建库存服务客户端
func NewInventoryServiceClientFromConfig(cfg *config.ServicesConfig) (*InventoryServiceClient, error) {
	return NewInventoryServiceClient(&cfg.InventoryService)
}

// NewManagerFromConfig 从配置创建客户端管理器
func NewManagerFromConfig(cfg *config.ServicesConfig) (*Manager, error) {
	return NewManager(cfg)
}
