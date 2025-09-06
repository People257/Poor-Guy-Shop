package client

import "context"

// 定义简化的接口，避免跨服务proto依赖

// InventoryItem 库存商品项
type InventoryItem struct {
	SkuID    string `json:"sku_id"`
	Quantity int32  `json:"quantity"`
}

// InventoryResponse 库存响应
type InventoryResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// PaymentRequest 支付请求
type PaymentRequest struct {
	OrderID       string `json:"order_id"`
	Amount        string `json:"amount"`
	PaymentMethod string `json:"payment_method"`
	Subject       string `json:"subject"`
	Description   string `json:"description"`
	NotifyURL     string `json:"notify_url"`
	ReturnURL     string `json:"return_url"`
}

// PaymentResponse 支付响应
type PaymentResponse struct {
	Success       bool              `json:"success"`
	PaymentURL    string            `json:"payment_url"`
	QRCode        string            `json:"qr_code"`
	PaymentParams map[string]string `json:"payment_params"`
}

// InventoryServiceInterface 库存服务接口
type InventoryServiceInterface interface {
	ReserveInventory(ctx context.Context, orderID string, items []InventoryItem) (*InventoryResponse, error)
	ReleaseInventory(ctx context.Context, orderID string) (*InventoryResponse, error)
	ConfirmInventory(ctx context.Context, orderID string) (*InventoryResponse, error)
}

// PaymentServiceInterface 支付服务接口
type PaymentServiceInterface interface {
	CreatePayment(ctx context.Context, req *PaymentRequest) (*PaymentResponse, error)
	VerifyPayment(ctx context.Context, orderID string) (bool, error)
}
