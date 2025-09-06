package order

import (
	"context"
	"fmt"
	"time"

	"github.com/shopspring/decimal"

	"github.com/people257/poor-guy-shop/order-service/internal/domain/order"
	"github.com/people257/poor-guy-shop/order-service/internal/infra/client"
)

// Service 订单应用服务
type Service struct {
	orderRepo       order.Repository
	orderDS         order.DomainService
	userClient      *client.UserServiceClient
	productClient   *client.ProductServiceClient
	paymentClient   *client.PaymentServiceClient
	inventoryClient *client.InventoryServiceClient
}

// NewService 创建订单应用服务
func NewService(
	orderRepo order.Repository,
	orderDS order.DomainService,
	userClient *client.UserServiceClient,
	productClient *client.ProductServiceClient,
	paymentClient *client.PaymentServiceClient,
	inventoryClient *client.InventoryServiceClient,
) *Service {
	return &Service{
		orderRepo:       orderRepo,
		orderDS:         orderDS,
		userClient:      userClient,
		productClient:   productClient,
		paymentClient:   paymentClient,
		inventoryClient: inventoryClient,
	}
}

// CreateOrderRequest 创建订单请求
type CreateOrderRequest struct {
	UserID         string                    `json:"user_id"`
	Items          []CreateOrderItemRequest  `json:"items"`
	Address        CreateOrderAddressRequest `json:"address"`
	PaymentMethod  string                    `json:"payment_method"`
	Remark         string                    `json:"remark"`
	DiscountAmount decimal.Decimal           `json:"discount_amount"`
	ShippingFee    decimal.Decimal           `json:"shipping_fee"`
}

// CreateOrderItemRequest 创建订单商品项请求
type CreateOrderItemRequest struct {
	ProductID   string          `json:"product_id"`
	SkuID       string          `json:"sku_id"`
	Quantity    int32           `json:"quantity"`
	Price       decimal.Decimal `json:"price"`
	ProductName string          `json:"product_name"`
	SkuName     string          `json:"sku_name"`
}

// CreateOrderAddressRequest 创建订单地址请求
type CreateOrderAddressRequest struct {
	ReceiverName  string `json:"receiver_name"`
	ReceiverPhone string `json:"receiver_phone"`
	Province      string `json:"province"`
	City          string `json:"city"`
	District      string `json:"district"`
	DetailAddress string `json:"detail_address"`
	PostalCode    string `json:"postal_code"`
}

// CreateOrder 创建订单
func (s *Service) CreateOrder(ctx context.Context, req CreateOrderRequest) (*order.Order, error) {
	// 1. 验证商品信息和库存
	for _, item := range req.Items {
		// 获取商品信息
		product, err := s.productClient.GetProduct(ctx, item.ProductID)
		if err != nil {
			return nil, fmt.Errorf("获取商品信息失败: %w", err)
		}

		// 检查商品状态
		if product.Status != 2 { // ACTIVE
			return nil, fmt.Errorf("商品 %s 不可购买", product.Name)
		}

		// 获取SKU信息
		sku, err := s.productClient.GetProductSKU(ctx, item.ProductID, item.SkuID)
		if err != nil {
			return nil, fmt.Errorf("获取商品SKU信息失败: %w", err)
		}

		// 检查SKU状态
		if !sku.IsActive {
			return nil, fmt.Errorf("商品SKU %s 不可购买", sku.Name)
		}

		// 检查库存 - 暂时注释掉，改用库存服务检查
		// if sku.StockQuantity < item.Quantity {
		// 	return nil, fmt.Errorf("商品 %s 库存不足，当前库存: %d，需要: %d",
		// 		sku.Name, sku.StockQuantity, item.Quantity)
		// }
	}

	// 2. 构建订单实体
	orderEntity := &order.Order{
		UserID:         req.UserID,
		Status:         int32(order.OrderStatusPendingPayment),
		PaymentMethod:  req.PaymentMethod,
		PaymentStatus:  int32(order.PaymentStatusUnpaid),
		Remark:         req.Remark,
		DiscountAmount: req.DiscountAmount,
		ShippingFee:    req.ShippingFee,
		CreatedAt:      time.Now().Format("2006-01-02 15:04:05"),
		UpdatedAt:      time.Now().Format("2006-01-02 15:04:05"),
	}

	// 3. 构建订单商品项
	var items []*order.OrderItem
	var totalAmount decimal.Decimal
	for _, item := range req.Items {
		orderItem := &order.OrderItem{
			ProductID:   item.ProductID,
			SkuID:       item.SkuID,
			Quantity:    item.Quantity,
			Price:       item.Price,
			ProductName: item.ProductName,
			SkuName:     item.SkuName,
			CreatedAt:   time.Now().Format("2006-01-02 15:04:05"),
			UpdatedAt:   time.Now().Format("2006-01-02 15:04:05"),
		}
		items = append(items, orderItem)

		// 计算小计
		subtotal := item.Price.Mul(decimal.NewFromInt32(item.Quantity))
		totalAmount = totalAmount.Add(subtotal)
	}

	// 4. 计算实付金额
	actualAmount := totalAmount.Add(req.ShippingFee).Sub(req.DiscountAmount)
	orderEntity.TotalAmount = totalAmount
	orderEntity.ActualAmount = actualAmount

	// 5. 构建订单地址
	orderAddress := &order.OrderAddress{
		ReceiverName:  req.Address.ReceiverName,
		ReceiverPhone: req.Address.ReceiverPhone,
		Province:      req.Address.Province,
		City:          req.Address.City,
		District:      req.Address.District,
		DetailAddress: req.Address.DetailAddress,
		PostalCode:    req.Address.PostalCode,
		CreatedAt:     time.Now().Format("2006-01-02 15:04:05"),
		UpdatedAt:     time.Now().Format("2006-01-02 15:04:05"),
	}

	// 6. 使用领域服务创建订单
	createdOrder, err := s.orderDS.CreateOrder(ctx, orderEntity, items, orderAddress)
	if err != nil {
		return nil, fmt.Errorf("创建订单失败: %w", err)
	}

	// 7. 预占库存
	err = s.reserveInventoryForOrder(ctx, createdOrder.ID, req.Items)
	if err != nil {
		// TODO: 这里应该回滚订单创建，但目前简化处理
		return nil, fmt.Errorf("库存预占失败: %w", err)
	}

	return createdOrder, nil
}

// GetOrderRequest 获取订单请求
type GetOrderRequest struct {
	OrderID string `json:"order_id"`
	UserID  string `json:"user_id"`
}

// GetOrder 获取订单详情
func (s *Service) GetOrder(ctx context.Context, req GetOrderRequest) (*order.Order, error) {
	orderEntity, err := s.orderRepo.GetByID(ctx, req.OrderID)
	if err != nil {
		return nil, fmt.Errorf("获取订单失败: %w", err)
	}

	// 检查订单是否属于该用户
	if orderEntity.UserID != req.UserID {
		return nil, order.ErrOrderNotFound
	}

	return orderEntity, nil
}

// ListOrdersRequest 获取订单列表请求
type ListOrdersRequest struct {
	UserID   string `json:"user_id"`
	Status   int32  `json:"status"`
	Page     int32  `json:"page"`
	PageSize int32  `json:"page_size"`
}

// ListOrdersResponse 获取订单列表响应
type ListOrdersResponse struct {
	Orders     []*order.Order `json:"orders"`
	Total      int64          `json:"total"`
	Page       int32          `json:"page"`
	PageSize   int32          `json:"page_size"`
	TotalPages int32          `json:"total_pages"`
}

// ListOrders 获取订单列表
func (s *Service) ListOrders(ctx context.Context, req ListOrdersRequest) (*ListOrdersResponse, error) {
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	orders, total, err := s.orderRepo.ListByUserID(ctx, req.UserID, req.Status, req.Page, req.PageSize)
	if err != nil {
		return nil, fmt.Errorf("获取订单列表失败: %w", err)
	}

	totalPages := int32((total + int64(req.PageSize) - 1) / int64(req.PageSize))

	return &ListOrdersResponse{
		Orders:     orders,
		Total:      total,
		Page:       req.Page,
		PageSize:   req.PageSize,
		TotalPages: totalPages,
	}, nil
}

// UpdateOrderStatusRequest 更新订单状态请求
type UpdateOrderStatusRequest struct {
	OrderID string `json:"order_id"`
	UserID  string `json:"user_id"`
	Status  int32  `json:"status"`
	Reason  string `json:"reason"`
}

// UpdateOrderStatus 更新订单状态
func (s *Service) UpdateOrderStatus(ctx context.Context, req UpdateOrderStatusRequest) error {
	// 获取订单
	orderEntity, err := s.orderRepo.GetByID(ctx, req.OrderID)
	if err != nil {
		return fmt.Errorf("获取订单失败: %w", err)
	}

	// 检查订单是否属于该用户
	if orderEntity.UserID != req.UserID {
		return order.ErrOrderNotFound
	}

	// 使用领域服务更新状态
	return s.orderDS.UpdateOrderStatus(ctx, orderEntity, req.Status, req.Reason)
}

// CancelOrderRequest 取消订单请求
type CancelOrderRequest struct {
	OrderID string `json:"order_id"`
	UserID  string `json:"user_id"`
	Reason  string `json:"reason"`
}

// CancelOrder 取消订单
func (s *Service) CancelOrder(ctx context.Context, req CancelOrderRequest) error {
	return s.UpdateOrderStatus(ctx, UpdateOrderStatusRequest{
		OrderID: req.OrderID,
		UserID:  req.UserID,
		Status:  int32(order.OrderStatusCancelled),
		Reason:  req.Reason,
	})
}

// PayOrderRequest 支付订单请求
type PayOrderRequest struct {
	OrderID       string `json:"order_id"`
	UserID        string `json:"user_id"`
	PaymentMethod string `json:"payment_method"`
}

// PayOrder 支付订单
func (s *Service) PayOrder(ctx context.Context, req PayOrderRequest) error {
	// 获取订单
	orderEntity, err := s.orderRepo.GetByID(ctx, req.OrderID)
	if err != nil {
		return fmt.Errorf("获取订单失败: %w", err)
	}

	// 检查订单是否属于该用户
	if orderEntity.UserID != req.UserID {
		return order.ErrOrderNotFound
	}

	// 使用领域服务处理支付
	return s.orderDS.PayOrder(ctx, orderEntity, req.PaymentMethod)
}

// reserveInventoryForOrder 为订单预占库存
func (s *Service) reserveInventoryForOrder(ctx context.Context, orderID string, items []CreateOrderItemRequest) error {
	// 构建库存预占请求
	var inventoryItems []client.InventoryItem
	for _, item := range items {
		inventoryItems = append(inventoryItems, client.InventoryItem{
			SkuID:    item.SkuID,
			Quantity: item.Quantity,
		})
	}

	// 调用库存服务预占库存
	resp, err := s.inventoryClient.ReserveInventory(ctx, orderID, inventoryItems)
	if err != nil {
		return fmt.Errorf("调用库存服务失败: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("库存预占失败: %s", resp.Message)
	}

	return nil
}

// createPaymentForOrder 为订单创建支付
func (s *Service) createPaymentForOrder(ctx context.Context, orderID string, amount decimal.Decimal, paymentMethod string) error {
	// 构建支付创建请求
	req := &client.PaymentRequest{
		OrderID:       orderID,
		Amount:        amount.String(),
		PaymentMethod: s.convertToPaymentMethod(paymentMethod),
		Subject:       fmt.Sprintf("订单支付-%s", orderID),
		Description:   "商城订单支付",
		NotifyURL:     "http://localhost:9002/payment/callback", // TODO: 配置化
		ReturnURL:     "http://localhost:8080/order/success",    // TODO: 配置化
	}

	resp, err := s.paymentClient.CreatePayment(ctx, req)
	if err != nil {
		return fmt.Errorf("创建支付订单失败: %w", err)
	}

	// 支付创建成功，由于我们修改了支付逻辑为自动成功，所以这里直接返回
	_ = resp // 暂时不处理支付响应
	return nil
}

// convertToPaymentMethod 转换支付方式
func (s *Service) convertToPaymentMethod(method string) string {
	switch method {
	case "alipay":
		return "alipay"
	case "wechat":
		return "wechat"
	case "balance":
		return "balance"
	default:
		return "alipay" // 默认支付宝
	}
}
