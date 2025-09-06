package client

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/people257/poor-guy-shop/order-service/cmd/grpc/config"
)

// 临时定义，后续需要引入product-service的proto
type ProductService interface {
	GetProduct(ctx context.Context, req *GetProductReq, opts ...grpc.CallOption) (*GetProductResp, error)
	ListProductSKUs(ctx context.Context, req *ListProductSKUsReq, opts ...grpc.CallOption) (*ListProductSKUsResp, error)
}

type GetProductReq struct {
	ID string
}

type GetProductResp struct {
	Product *Product
}

type ListProductSKUsReq struct {
	ProductID string
	IsActive  bool
}

type ListProductSKUsResp struct {
	SKUs []*ProductSKU
}

type Product struct {
	ID        string
	Name      string
	Status    int32
	SalePrice string
}

type ProductSKU struct {
	ID            string
	ProductID     string
	SKUCode       string
	Name          string
	SalePrice     string
	StockQuantity int32
	IsActive      bool
}

// ProductServiceClient 产品服务客户端
type ProductServiceClient struct {
	conn           *grpc.ClientConn
	productService ProductService
}

// NewProductServiceClient 创建产品服务客户端
func NewProductServiceClient(cfg *config.ServiceConfig) (*ProductServiceClient, error) {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to product service: %w", err)
	}

	// 这里需要使用真实的product-service proto client
	// productService := productpb.NewProductServiceClient(conn)

	return &ProductServiceClient{
		conn: conn,
		// productService: productService,
	}, nil
}

// GetProduct 获取产品信息
func (c *ProductServiceClient) GetProduct(ctx context.Context, productID string) (*Product, error) {
	// req := &GetProductReq{ID: productID}
	// resp, err := c.productService.GetProduct(ctx, req)
	// if err != nil {
	//     return nil, fmt.Errorf("failed to get product: %w", err)
	// }
	// return resp.Product, nil

	// 临时返回固定产品信息
	return &Product{
		ID:        productID,
		Name:      "Test Product",
		Status:    2, // ACTIVE
		SalePrice: "99.99",
	}, nil
}

// GetProductSKU 获取产品SKU信息
func (c *ProductServiceClient) GetProductSKU(ctx context.Context, productID string, skuID string) (*ProductSKU, error) {
	// req := &ListProductSKUsReq{
	//     ProductID: productID,
	//     IsActive:  true,
	// }
	// resp, err := c.productService.ListProductSKUs(ctx, req)
	// if err != nil {
	//     return nil, fmt.Errorf("failed to get product SKUs: %w", err)
	// }

	// // 查找指定的SKU
	// for _, sku := range resp.SKUs {
	//     if sku.ID == skuID {
	//         return sku, nil
	//     }
	// }

	// return nil, fmt.Errorf("SKU not found: %s", skuID)

	// 临时返回固定SKU信息
	return &ProductSKU{
		ID:            skuID,
		ProductID:     productID,
		SKUCode:       "TEST-SKU-001",
		Name:          "Test SKU",
		SalePrice:     "99.99",
		StockQuantity: 100,
		IsActive:      true,
	}, nil
}

// CheckStock 检查库存
func (c *ProductServiceClient) CheckStock(ctx context.Context, skuID string, quantity int32) (bool, error) {
	sku, err := c.GetProductSKU(ctx, "", skuID)
	if err != nil {
		return false, err
	}

	return sku.StockQuantity >= quantity, nil
}

// Close 关闭连接
func (c *ProductServiceClient) Close() error {
	return c.conn.Close()
}
