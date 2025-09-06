package order

import (
	"github.com/shopspring/decimal"
)

// OrderStatus 订单状态
type OrderStatus int

const (
	OrderStatusUnknown        OrderStatus = 0
	OrderStatusPendingPayment OrderStatus = 1 // 待付款
	OrderStatusPaid           OrderStatus = 2 // 已付款
	OrderStatusShipped        OrderStatus = 3 // 已发货
	OrderStatusDelivered      OrderStatus = 4 // 已收货
	OrderStatusCancelled      OrderStatus = 5 // 已取消
	OrderStatusRefunded       OrderStatus = 6 // 已退款
)

// PaymentMethod 支付方式
type PaymentMethod int

const (
	PaymentMethodUnknown PaymentMethod = 0
	PaymentMethodAlipay  PaymentMethod = 1 // 支付宝
	PaymentMethodWechat  PaymentMethod = 2 // 微信支付
	PaymentMethodBalance PaymentMethod = 3 // 余额支付
)

// PaymentStatus 支付状态
type PaymentStatus int

const (
	PaymentStatusUnpaid    PaymentStatus = 0 // 未支付
	PaymentStatusPaid      PaymentStatus = 1 // 已支付
	PaymentStatusRefunding PaymentStatus = 2 // 退款中
	PaymentStatusRefunded  PaymentStatus = 3 // 已退款
)

// Order 订单实体（匹配数据库模型）
type Order struct {
	ID             string          `json:"id"`
	OrderNo        string          `json:"order_no"`
	UserID         string          `json:"user_id"`
	Status         int32           `json:"status"`
	TotalAmount    decimal.Decimal `json:"total_amount"`
	DiscountAmount decimal.Decimal `json:"discount_amount"`
	ShippingFee    decimal.Decimal `json:"shipping_fee"`
	ActualAmount   decimal.Decimal `json:"actual_amount"`
	PaymentMethod  string          `json:"payment_method"`
	PaymentStatus  int32           `json:"payment_status"`
	PaymentTime    string          `json:"payment_time"`
	DeliveryTime   string          `json:"delivery_time"`
	ReceiveTime    string          `json:"receive_time"`
	CancelTime     string          `json:"cancel_time"`
	CancelReason   string          `json:"cancel_reason"`
	Remark         string          `json:"remark"`
	CreatedAt      string          `json:"created_at"`
	UpdatedAt      string          `json:"updated_at"`
	Version        int32           `json:"version"`
}

// OrderItem 订单商品项实体（匹配数据库模型）
type OrderItem struct {
	ID          string          `json:"id"`
	OrderID     string          `json:"order_id"`
	ProductID   string          `json:"product_id"`
	SkuID       string          `json:"sku_id"`
	ProductName string          `json:"product_name"`
	SkuName     string          `json:"sku_name"`
	Price       decimal.Decimal `json:"price"`
	Quantity    int32           `json:"quantity"`
	CreatedAt   string          `json:"created_at"`
	UpdatedAt   string          `json:"updated_at"`
}

// OrderAddress 订单收货地址实体（匹配数据库模型）
type OrderAddress struct {
	ID            string `json:"id"`
	OrderID       string `json:"order_id"`
	ReceiverName  string `json:"receiver_name"`
	ReceiverPhone string `json:"receiver_phone"`
	Province      string `json:"province"`
	City          string `json:"city"`
	District      string `json:"district"`
	DetailAddress string `json:"detail_address"`
	PostalCode    string `json:"postal_code"`
	CreatedAt     string `json:"created_at"`
	UpdatedAt     string `json:"updated_at"`
}

// OrderStatusLog 订单状态变更日志（匹配数据库模型）
type OrderStatusLog struct {
	ID        string `json:"id"`
	OrderID   string `json:"order_id"`
	Status    int32  `json:"status"`
	Remark    string `json:"remark"`
	CreatedAt string `json:"created_at"`
}

// OrderPayment 订单支付记录（匹配数据库模型）
type OrderPayment struct {
	ID            string          `json:"id"`
	OrderID       string          `json:"order_id"`
	PaymentNo     string          `json:"payment_no"`
	PaymentMethod string          `json:"payment_method"`
	Amount        decimal.Decimal `json:"amount"`
	Status        int32           `json:"status"`
	PaidAt        string          `json:"paid_at"`
	CreatedAt     string          `json:"created_at"`
	UpdatedAt     string          `json:"updated_at"`
}

// 简单的业务方法
// CanCancel 检查订单是否可以取消
func (o *Order) CanCancel() bool {
	return o.Status == int32(OrderStatusPendingPayment) || o.Status == int32(OrderStatusPaid)
}

// CanPay 检查订单是否可以支付
func (o *Order) CanPay() bool {
	return o.Status == int32(OrderStatusPendingPayment)
}

// CanConfirm 检查订单是否可以确认收货
func (o *Order) CanConfirm() bool {
	return o.Status == int32(OrderStatusShipped)
}

// IsCompleted 检查订单是否完成
func (o *Order) IsCompleted() bool {
	return o.Status == int32(OrderStatusDelivered)
}

// IsCancelled 检查订单是否已取消
func (o *Order) IsCancelled() bool {
	return o.Status == int32(OrderStatusCancelled)
}

// IsPaid 检查订单是否已支付
func (o *Order) IsPaid() bool {
	return o.PaymentStatus == int32(PaymentStatusPaid)
}
