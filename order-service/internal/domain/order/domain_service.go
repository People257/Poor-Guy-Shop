package order

import (
	"context"
	"fmt"
	"time"
)

// DomainService 订单领域服务
type DomainService interface {
	// 创建订单
	CreateOrder(ctx context.Context, order *Order, items []*OrderItem, address *OrderAddress) (*Order, error)

	// 更新订单状态
	UpdateOrderStatus(ctx context.Context, order *Order, status int32, reason string) error

	// 支付订单
	PayOrder(ctx context.Context, order *Order, paymentMethod string) error

	// 生成订单号
	GenerateOrderNo() string
}

// domainService 订单领域服务实现
type domainService struct {
	orderRepo Repository
}

// NewDomainService 创建订单领域服务
func NewDomainService(orderRepo Repository) DomainService {
	return &domainService{
		orderRepo: orderRepo,
	}
}

// CreateOrder 创建订单
func (ds *domainService) CreateOrder(ctx context.Context, order *Order, items []*OrderItem, address *OrderAddress) (*Order, error) {
	// 生成订单号
	order.OrderNo = ds.GenerateOrderNo()

	// 设置初始状态
	order.Status = int32(OrderStatusPendingPayment)
	order.PaymentStatus = int32(PaymentStatusUnpaid)

	// 调用仓储创建订单
	if err := ds.orderRepo.Create(ctx, order, items, address); err != nil {
		return nil, fmt.Errorf("创建订单失败: %w", err)
	}

	return order, nil
}

// UpdateOrderStatus 更新订单状态
func (ds *domainService) UpdateOrderStatus(ctx context.Context, order *Order, status int32, reason string) error {
	// 验证状态转换是否合法
	if !ds.isValidStatusTransition(order.Status, status) {
		return fmt.Errorf("无效的状态转换: %d -> %d", order.Status, status)
	}

	// 更新订单状态
	oldStatus := order.Status
	order.Status = status
	order.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")

	// 根据状态设置相应的时间字段
	switch status {
	case int32(OrderStatusCancelled):
		order.CancelTime = time.Now().Format("2006-01-02 15:04:05")
		order.CancelReason = reason
	case int32(OrderStatusShipped):
		order.DeliveryTime = time.Now().Format("2006-01-02 15:04:05")
	case int32(OrderStatusDelivered):
		order.ReceiveTime = time.Now().Format("2006-01-02 15:04:05")
	}

	// 更新数据库
	if err := ds.orderRepo.Update(ctx, order); err != nil {
		return fmt.Errorf("更新订单状态失败: %w", err)
	}

	// 记录状态日志
	logReason := reason
	if logReason == "" {
		logReason = fmt.Sprintf("状态从 %d 更新为 %d", oldStatus, status)
	}

	if err := ds.orderRepo.CreateStatusLog(ctx, order.ID, status, logReason); err != nil {
		// 日志记录失败不影响主流程
		// 这里可以记录错误日志
	}

	return nil
}

// PayOrder 支付订单
func (ds *domainService) PayOrder(ctx context.Context, order *Order, paymentMethod string) error {
	// 检查订单状态
	if order.Status != int32(OrderStatusPendingPayment) {
		return fmt.Errorf("订单状态不是待付款，无法支付")
	}

	// 更新支付信息
	order.PaymentMethod = paymentMethod
	order.PaymentStatus = int32(PaymentStatusPaid)
	order.PaymentTime = time.Now().Format("2006-01-02 15:04:05")
	order.Status = int32(OrderStatusPaid)
	order.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")

	// 更新数据库
	if err := ds.orderRepo.Update(ctx, order); err != nil {
		return fmt.Errorf("更新订单支付信息失败: %w", err)
	}

	// 记录状态日志
	if err := ds.orderRepo.CreateStatusLog(ctx, order.ID, int32(OrderStatusPaid), "订单支付成功"); err != nil {
		// 日志记录失败不影响主流程
	}

	return nil
}

// GenerateOrderNo 生成订单号
func (ds *domainService) GenerateOrderNo() string {
	// 简单的订单号生成规则：ORD + 时间戳 + 随机数
	now := time.Now()
	return fmt.Sprintf("ORD%s%03d", now.Format("20060102150405"), now.Nanosecond()%1000)
}

// isValidStatusTransition 检查状态转换是否合法
func (ds *domainService) isValidStatusTransition(from, to int32) bool {
	// 定义合法的状态转换
	validTransitions := map[int32][]int32{
		int32(OrderStatusPendingPayment): {int32(OrderStatusPaid), int32(OrderStatusCancelled)},
		int32(OrderStatusPaid):           {int32(OrderStatusShipped), int32(OrderStatusCancelled), int32(OrderStatusRefunded)},
		int32(OrderStatusShipped):        {int32(OrderStatusDelivered), int32(OrderStatusCancelled)},
		int32(OrderStatusDelivered):      {int32(OrderStatusCancelled)}, // 已收货还可以申请退货
		int32(OrderStatusCancelled):      {},                            // 已取消不能再转换
		int32(OrderStatusRefunded):       {},                            // 已退款不能再转换
	}

	allowedStates, exists := validTransitions[from]
	if !exists {
		return false
	}

	for _, allowedState := range allowedStates {
		if allowedState == to {
			return true
		}
	}

	return false
}
