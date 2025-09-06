package payment

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// DomainService 支付领域服务
type DomainService struct {
	paymentRepo Repository
}

// NewDomainService 创建支付领域服务
func NewDomainService(paymentRepo Repository) *DomainService {
	return &DomainService{
		paymentRepo: paymentRepo,
	}
}

// CreatePaymentOrder 创建支付订单
func (s *DomainService) CreatePaymentOrder(ctx context.Context, orderID, userID uuid.UUID, amount decimal.Decimal, paymentMethod PaymentMethod, subject, description string) (*PaymentOrder, error) {
	// 检查是否已存在支付订单
	existingPayment, err := s.paymentRepo.GetByOrderID(ctx, orderID)
	if err == nil && existingPayment != nil {
		// 如果已存在且状态为待支付，返回现有订单
		if existingPayment.Status == PaymentStatusPending && !existingPayment.IsExpired() {
			return existingPayment, nil
		}
		// 如果已存在但已过期或失败，创建新的支付订单
		if existingPayment.Status == PaymentStatusPending && existingPayment.IsExpired() {
			existingPayment.MarkAsFailed("payment expired")
			if err := s.paymentRepo.Update(ctx, existingPayment); err != nil {
				return nil, fmt.Errorf("failed to update expired payment: %w", err)
			}
		}
	}

	// 验证金额
	if amount.LessThanOrEqual(decimal.Zero) {
		return nil, fmt.Errorf("payment amount must be greater than zero")
	}

	// 创建新的支付订单
	paymentOrder := NewPaymentOrder(orderID, userID, amount, paymentMethod, subject, description)

	// 保存到数据库
	if err := s.paymentRepo.Create(ctx, paymentOrder); err != nil {
		return nil, fmt.Errorf("failed to create payment order: %w", err)
	}

	// 记录日志
	log := NewPaymentLog(paymentOrder.ID, "create", fmt.Sprintf("amount: %s, method: %s", amount.String(), paymentMethod), "", "")
	if err := s.paymentRepo.CreateLog(ctx, log); err != nil {
		// 日志记录失败不影响主流程
		fmt.Printf("failed to create payment log: %v\n", err)
	}

	return paymentOrder, nil
}

// ProcessPaymentCallback 处理支付回调
func (s *DomainService) ProcessPaymentCallback(ctx context.Context, orderIdentifier string, isSuccess bool, response string) (*PaymentOrder, error) {
	// 尝试解析为UUID，如果成功则按ID查找，否则按第三方订单ID查找
	var paymentOrder *PaymentOrder
	var err error

	if paymentID, parseErr := uuid.Parse(orderIdentifier); parseErr == nil {
		paymentOrder, err = s.paymentRepo.GetByID(ctx, paymentID)
	} else {
		paymentOrder, err = s.paymentRepo.GetByThirdPartyOrderID(ctx, orderIdentifier)
	}

	if err != nil {
		return nil, fmt.Errorf("payment order not found: %w", err)
	}

	// 检查支付订单状态
	if paymentOrder.Status != PaymentStatusPending {
		return paymentOrder, fmt.Errorf("payment order status is not pending: %s", paymentOrder.Status)
	}

	// 更新支付订单状态
	if isSuccess {
		paymentOrder.MarkAsPaid(orderIdentifier, response)
	} else {
		paymentOrder.MarkAsFailed(response)
	}

	// 保存更新
	if err := s.paymentRepo.Update(ctx, paymentOrder); err != nil {
		return nil, fmt.Errorf("failed to update payment order: %w", err)
	}

	// 记录日志
	action := "callback_success"
	if !isSuccess {
		action = "callback_failed"
	}
	log := NewPaymentLog(paymentOrder.ID, action, "", response, "")
	if err := s.paymentRepo.CreateLog(ctx, log); err != nil {
		fmt.Printf("failed to create payment log: %v\n", err)
	}

	return paymentOrder, nil
}

// CancelPaymentOrder 取消支付订单
func (s *DomainService) CancelPaymentOrder(ctx context.Context, paymentOrderID uuid.UUID) error {
	paymentOrder, err := s.paymentRepo.GetByID(ctx, paymentOrderID)
	if err != nil {
		return fmt.Errorf("payment order not found: %w", err)
	}

	// 只有待支付状态的订单可以取消
	if paymentOrder.Status != PaymentStatusPending {
		return fmt.Errorf("cannot cancel payment order with status: %s", paymentOrder.Status)
	}

	paymentOrder.MarkAsCancelled()

	if err := s.paymentRepo.Update(ctx, paymentOrder); err != nil {
		return fmt.Errorf("failed to cancel payment order: %w", err)
	}

	// 记录日志
	log := NewPaymentLog(paymentOrder.ID, "cancel", "", "", "")
	if err := s.paymentRepo.CreateLog(ctx, log); err != nil {
		fmt.Printf("failed to create payment log: %v\n", err)
	}

	return nil
}

// ValidatePaymentStatus 验证支付状态
func (s *DomainService) ValidatePaymentStatus(ctx context.Context, orderID uuid.UUID) (*PaymentOrder, error) {
	paymentOrder, err := s.paymentRepo.GetByOrderID(ctx, orderID)
	if err != nil {
		return nil, fmt.Errorf("payment order not found for order: %s", orderID)
	}

	return paymentOrder, nil
}
