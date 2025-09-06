package refund

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/people257/poor-guy-shop/payment-service/internal/domain/payment"
	"github.com/people257/poor-guy-shop/payment-service/internal/domain/refund"
	infraPayment "github.com/people257/poor-guy-shop/payment-service/internal/infra/payment"
	"github.com/shopspring/decimal"
)

// Service 退款应用服务
type Service struct {
	refundDS      *refund.DomainService
	refundRepo    refund.Repository
	paymentRepo   payment.Repository
	paymentClient infraPayment.PaymentClient
}

// NewService 创建退款应用服务
func NewService(
	refundDS *refund.DomainService,
	refundRepo refund.Repository,
	paymentRepo payment.Repository,
	paymentClient infraPayment.PaymentClient,
) *Service {
	return &Service{
		refundDS:      refundDS,
		refundRepo:    refundRepo,
		paymentRepo:   paymentRepo,
		paymentClient: paymentClient,
	}
}

// CreateRefundRequest 创建退款请求
type CreateRefundRequest struct {
	PaymentID string `json:"payment_id"`
	Amount    string `json:"amount"`
	Reason    string `json:"reason"`
}

// CreateRefundResponse 创建退款响应
type CreateRefundResponse struct {
	Refund *refund.Refund `json:"refund"`
}

// CreateRefund 创建退款
func (s *Service) CreateRefund(ctx context.Context, req CreateRefundRequest) (*CreateRefundResponse, error) {
	// 解析退款金额
	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		return nil, fmt.Errorf("invalid amount: %w", err)
	}

	// 获取支付订单
	paymentID, err := uuid.Parse(req.PaymentID)
	if err != nil {
		return nil, fmt.Errorf("invalid payment ID: %w", err)
	}

	paymentOrder, err := s.paymentRepo.GetByID(ctx, paymentID)
	if err != nil {
		return nil, fmt.Errorf("payment order not found: %w", err)
	}

	// 检查支付订单是否可以退款
	if !paymentOrder.CanRefund() {
		return nil, fmt.Errorf("payment order cannot be refunded, status: %s", paymentOrder.Status)
	}

	// 检查退款金额是否超过支付金额
	if amount.GreaterThan(paymentOrder.Amount) {
		return nil, fmt.Errorf("refund amount cannot exceed payment amount")
	}

	// 计算已退款金额
	totalRefunded, err := s.refundDS.CalculateTotalRefunded(ctx, paymentID)
	if err != nil {
		return nil, fmt.Errorf("failed to calculate total refunded: %w", err)
	}

	// 检查剩余可退款金额
	remainingAmount := paymentOrder.Amount.Sub(totalRefunded)
	if amount.GreaterThan(remainingAmount) {
		return nil, fmt.Errorf("refund amount exceeds remaining refundable amount: %s", remainingAmount.String())
	}

	// 创建退款记录
	refundEntity, err := s.refundDS.CreateRefund(ctx, paymentID, amount, req.Reason)
	if err != nil {
		return nil, fmt.Errorf("failed to create refund: %w", err)
	}

	// 调用第三方退款
	switch paymentOrder.PaymentMethod {
	case payment.PaymentMethodAlipay:
		err = s.processAlipayRefund(ctx, paymentOrder, refundEntity)
	case payment.PaymentMethodWechat:
		// TODO: 实现微信退款
		err = fmt.Errorf("wechat refund not implemented yet")
	default:
		err = fmt.Errorf("unsupported payment method for refund: %s", paymentOrder.PaymentMethod)
	}

	if err != nil {
		// 退款失败，更新退款状态
		refundEntity.MarkAsFailed(err.Error())
		if _, updateErr := s.refundDS.ProcessRefundCallback(ctx, refundEntity.ID, false, "", err.Error()); updateErr != nil {
			return nil, fmt.Errorf("failed to update refund status: %w", updateErr)
		}
		return nil, fmt.Errorf("failed to process refund: %w", err)
	}

	return &CreateRefundResponse{
		Refund: refundEntity,
	}, nil
}

// processAlipayRefund 处理支付宝退款
func (s *Service) processAlipayRefund(ctx context.Context, paymentOrder *payment.PaymentOrder, refundEntity *refund.Refund) error {
	refundReq := &infraPayment.CreateRefundRequest{
		OutTradeNo:  paymentOrder.ID.String(),
		OutRefundNo: refundEntity.ID.String(),
		Amount:      refundEntity.Amount,
		Reason:      refundEntity.Reason,
	}

	refundResp, err := s.paymentClient.CreateRefund(ctx, refundReq)
	if err != nil {
		return fmt.Errorf("failed to create alipay refund: %w", err)
	}

	if refundResp.IsSuccess {
		// 退款成功
		_, err = s.refundDS.ProcessRefundCallback(
			ctx,
			refundEntity.ID,
			true,
			refundResp.ThirdPartyRefundNo,
			"refund success",
		)
		if err != nil {
			return fmt.Errorf("failed to update refund status: %w", err)
		}

		// 更新支付订单状态
		if refundEntity.Amount.Equal(paymentOrder.Amount) {
			// 全额退款
			paymentOrder.MarkAsRefunded()
		} else {
			// 部分退款
			paymentOrder.MarkAsPartialRefunded()
		}

		if err := s.paymentRepo.Update(ctx, paymentOrder); err != nil {
			return fmt.Errorf("failed to update payment order status: %w", err)
		}
	} else {
		// 退款失败
		_, err = s.refundDS.ProcessRefundCallback(
			ctx,
			refundEntity.ID,
			false,
			"",
			refundResp.ErrorMsg,
		)
		if err != nil {
			return fmt.Errorf("failed to update refund status: %w", err)
		}
		return fmt.Errorf("alipay refund failed: %s", refundResp.ErrorMsg)
	}

	return nil
}

// GetRefund 获取退款记录
func (s *Service) GetRefund(ctx context.Context, refundID string) (*refund.Refund, error) {
	refundUUID, err := uuid.Parse(refundID)
	if err != nil {
		return nil, fmt.Errorf("invalid refund ID: %w", err)
	}

	refundEntity, err := s.refundRepo.GetByID(ctx, refundUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get refund: %w", err)
	}

	return refundEntity, nil
}

// GetRefundsByPaymentOrder 获取支付订单的退款记录
func (s *Service) GetRefundsByPaymentOrder(ctx context.Context, paymentOrderID string) ([]*refund.Refund, error) {
	paymentUUID, err := uuid.Parse(paymentOrderID)
	if err != nil {
		return nil, fmt.Errorf("invalid payment order ID: %w", err)
	}

	refunds, err := s.refundDS.GetRefundsByPaymentOrder(ctx, paymentUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get refunds: %w", err)
	}

	return refunds, nil
}
