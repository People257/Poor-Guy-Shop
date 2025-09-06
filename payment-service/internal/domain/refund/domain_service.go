package refund

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// DomainService 退款领域服务
type DomainService struct {
	refundRepo Repository
}

// NewDomainService 创建退款领域服务
func NewDomainService(refundRepo Repository) *DomainService {
	return &DomainService{
		refundRepo: refundRepo,
	}
}

// CreateRefund 创建退款记录
func (s *DomainService) CreateRefund(ctx context.Context, paymentOrderID uuid.UUID, amount decimal.Decimal, reason string) (*Refund, error) {
	// 验证退款金额
	if amount.LessThanOrEqual(decimal.Zero) {
		return nil, fmt.Errorf("refund amount must be greater than zero")
	}

	// 检查是否已存在相同的退款申请
	existingRefunds, err := s.refundRepo.GetByPaymentOrderID(ctx, paymentOrderID)
	if err != nil {
		return nil, fmt.Errorf("failed to check existing refunds: %w", err)
	}

	// 计算已退款金额
	totalRefunded := decimal.Zero
	for _, refund := range existingRefunds {
		if refund.Status == RefundStatusSuccess {
			totalRefunded = totalRefunded.Add(refund.Amount)
		}
	}

	// 检查是否有正在处理的退款
	for _, refund := range existingRefunds {
		if refund.Status == RefundStatusPending {
			return nil, fmt.Errorf("there is already a pending refund for this payment order")
		}
	}

	// 创建退款记录
	refund := NewRefund(paymentOrderID, amount, reason)

	// 保存到数据库
	if err := s.refundRepo.Create(ctx, refund); err != nil {
		return nil, fmt.Errorf("failed to create refund: %w", err)
	}

	return refund, nil
}

// ProcessRefundCallback 处理退款回调
func (s *DomainService) ProcessRefundCallback(ctx context.Context, refundID uuid.UUID, isSuccess bool, thirdPartyRefundID, response string) (*Refund, error) {
	refund, err := s.refundRepo.GetByID(ctx, refundID)
	if err != nil {
		return nil, fmt.Errorf("refund not found: %w", err)
	}

	// 检查退款状态
	if refund.Status != RefundStatusPending {
		return refund, fmt.Errorf("refund status is not pending: %s", refund.Status)
	}

	// 更新退款状态
	if isSuccess {
		refund.MarkAsSuccess(thirdPartyRefundID, response)
	} else {
		refund.MarkAsFailed(response)
	}

	// 保存更新
	if err := s.refundRepo.Update(ctx, refund); err != nil {
		return nil, fmt.Errorf("failed to update refund: %w", err)
	}

	return refund, nil
}

// GetRefundsByPaymentOrder 获取支付订单的所有退款记录
func (s *DomainService) GetRefundsByPaymentOrder(ctx context.Context, paymentOrderID uuid.UUID) ([]*Refund, error) {
	refunds, err := s.refundRepo.GetByPaymentOrderID(ctx, paymentOrderID)
	if err != nil {
		return nil, fmt.Errorf("failed to get refunds: %w", err)
	}

	return refunds, nil
}

// CalculateTotalRefunded 计算已退款总金额
func (s *DomainService) CalculateTotalRefunded(ctx context.Context, paymentOrderID uuid.UUID) (decimal.Decimal, error) {
	refunds, err := s.refundRepo.GetByPaymentOrderID(ctx, paymentOrderID)
	if err != nil {
		return decimal.Zero, fmt.Errorf("failed to get refunds: %w", err)
	}

	total := decimal.Zero
	for _, refund := range refunds {
		if refund.Status == RefundStatusSuccess {
			total = total.Add(refund.Amount)
		}
	}

	return total, nil
}
