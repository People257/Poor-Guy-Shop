package product

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/people257/poor-guy-shop/inventory-service/cmd/grpc/config"
)

// ProductInfo 商品信息
type ProductInfo struct {
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	SKU         string    `json:"sku"`
	Price       float64   `json:"price"`
	Status      string    `json:"status"`
	CategoryID  uuid.UUID `json:"category_id"`
	Description string    `json:"description"`
}

// Client 商品服务客户端接口
type Client interface {
	// GetProduct 获取商品信息
	GetProduct(ctx context.Context, skuID uuid.UUID) (*ProductInfo, error)

	// BatchGetProducts 批量获取商品信息
	BatchGetProducts(ctx context.Context, skuIDs []uuid.UUID) ([]*ProductInfo, error)

	// ValidateProducts 验证商品是否存在且可售
	ValidateProducts(ctx context.Context, skuIDs []uuid.UUID) (map[uuid.UUID]bool, error)
}

// GrpcClient 商品服务gRPC客户端实现
type GrpcClient struct {
	conn   *grpc.ClientConn
	config *config.ServiceConfig
}

// NewGrpcClient 创建商品服务gRPC客户端
func NewGrpcClient(servicesConfig *config.ServicesConfig) (Client, error) {
	serviceConfig := servicesConfig.ProductService

	// 建立gRPC连接
	addr := fmt.Sprintf("%s:%d", serviceConfig.Host, serviceConfig.Port)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to product service: %w", err)
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

// GetProduct 获取商品信息
func (c *GrpcClient) GetProduct(ctx context.Context, skuID uuid.UUID) (*ProductInfo, error) {
	// TODO: 实现商品信息查询
	// 这里需要根据实际的商品服务proto定义来实现
	return &ProductInfo{
		ID:     skuID,
		Name:   "Sample Product",
		SKU:    skuID.String(),
		Price:  99.99,
		Status: "active",
	}, nil
}

// BatchGetProducts 批量获取商品信息
func (c *GrpcClient) BatchGetProducts(ctx context.Context, skuIDs []uuid.UUID) ([]*ProductInfo, error) {
	// TODO: 实现批量商品信息查询
	// 这里需要根据实际的商品服务proto定义来实现
	products := make([]*ProductInfo, len(skuIDs))
	for i, skuID := range skuIDs {
		products[i] = &ProductInfo{
			ID:     skuID,
			Name:   fmt.Sprintf("Product %d", i+1),
			SKU:    skuID.String(),
			Price:  99.99,
			Status: "active",
		}
	}
	return products, nil
}

// ValidateProducts 验证商品是否存在且可售
func (c *GrpcClient) ValidateProducts(ctx context.Context, skuIDs []uuid.UUID) (map[uuid.UUID]bool, error) {
	// TODO: 实现商品验证
	// 这里需要根据实际的商品服务proto定义来实现
	result := make(map[uuid.UUID]bool)
	for _, skuID := range skuIDs {
		result[skuID] = true // 假设所有商品都有效
	}
	return result, nil
}

