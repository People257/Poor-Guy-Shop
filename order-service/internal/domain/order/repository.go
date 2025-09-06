package order

import (
	"context"
)

// Repository 订单仓储接口
type Repository interface {
	// 创建订单（包括订单项和地址）
	Create(ctx context.Context, order *Order, items []*OrderItem, address *OrderAddress) error

	// 根据ID获取订单
	GetByID(ctx context.Context, id string) (*Order, error)

	// 根据订单号获取订单
	GetByOrderNo(ctx context.Context, orderNo string) (*Order, error)

	// 根据用户ID获取订单列表
	ListByUserID(ctx context.Context, userID string, status int32, page, pageSize int32) ([]*Order, int64, error)

	// 更新订单
	Update(ctx context.Context, order *Order) error

	// 删除订单（软删除）
	Delete(ctx context.Context, id string) error

	// 获取订单商品项
	GetOrderItems(ctx context.Context, orderID string) ([]*OrderItem, error)

	// 获取订单地址
	GetOrderAddress(ctx context.Context, orderID string) (*OrderAddress, error)

	// 创建状态日志
	CreateStatusLog(ctx context.Context, orderID string, status int32, remark string) error

	// 获取状态日志
	GetStatusLogs(ctx context.Context, orderID string) ([]*OrderStatusLog, error)
}

// OrderItemRepository 订单商品项仓储接口
type OrderItemRepository interface {
	// 创建订单商品项
	Create(ctx context.Context, item *OrderItem) error

	// 批量创建订单商品项
	BatchCreate(ctx context.Context, items []OrderItem) error

	// 根据订单ID获取商品项列表
	ListByOrderID(ctx context.Context, orderID string) ([]OrderItem, error)

	// 更新订单商品项
	Update(ctx context.Context, item *OrderItem) error

	// 删除订单商品项
	Delete(ctx context.Context, id string) error
}

// OrderAddressRepository 订单地址仓储接口
type OrderAddressRepository interface {
	// 创建订单地址
	Create(ctx context.Context, address *OrderAddress) error

	// 根据订单ID获取地址
	GetByOrderID(ctx context.Context, orderID string) (*OrderAddress, error)

	// 更新订单地址
	Update(ctx context.Context, address *OrderAddress) error

	// 删除订单地址
	Delete(ctx context.Context, id string) error
}

// OrderStatusLogRepository 订单状态日志仓储接口
type OrderStatusLogRepository interface {
	// 创建状态日志
	Create(ctx context.Context, log *OrderStatusLog) error

	// 根据订单ID获取状态日志列表
	ListByOrderID(ctx context.Context, orderID string) ([]OrderStatusLog, error)
}

// OrderPaymentRepository 订单支付记录仓储接口
type OrderPaymentRepository interface {
	// 创建支付记录
	Create(ctx context.Context, payment *OrderPayment) error

	// 根据支付流水号获取支付记录
	GetByPaymentNo(ctx context.Context, paymentNo string) (*OrderPayment, error)

	// 根据订单ID获取支付记录列表
	ListByOrderID(ctx context.Context, orderID string) ([]OrderPayment, error)

	// 更新支付记录
	Update(ctx context.Context, payment *OrderPayment) error
}

// OrderStats 订单统计信息
type OrderStats struct {
	TotalOrders     int64 `json:"total_orders"`
	PendingOrders   int64 `json:"pending_orders"`
	PaidOrders      int64 `json:"paid_orders"`
	ShippedOrders   int64 `json:"shipped_orders"`
	DeliveredOrders int64 `json:"delivered_orders"`
	CancelledOrders int64 `json:"cancelled_orders"`
}
