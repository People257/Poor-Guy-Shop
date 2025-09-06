package payment

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/people257/poor-guy-shop/payment-service/internal/domain/payment"
	infraPayment "github.com/people257/poor-guy-shop/payment-service/internal/infra/payment"
	"github.com/shopspring/decimal"
)

// Service 支付应用服务
type Service struct {
	paymentDS     *payment.DomainService
	paymentClient infraPayment.PaymentClient
}

// NewService 创建支付应用服务
func NewService(
	paymentDS *payment.DomainService,
	paymentClient infraPayment.PaymentClient,
) *Service {
	return &Service{
		paymentDS:     paymentDS,
		paymentClient: paymentClient,
	}
}

// CreatePaymentOrderRequest 创建支付订单请求
type CreatePaymentOrderRequest struct {
	OrderID       string                `json:"order_id"`
	Amount        string                `json:"amount"`
	PaymentMethod payment.PaymentMethod `json:"payment_method"`
	Subject       string                `json:"subject"`
	Description   string                `json:"description"`
	NotifyURL     string                `json:"notify_url"`
	ReturnURL     string                `json:"return_url"`
}

// CreatePaymentOrderResponse 创建支付订单响应
type CreatePaymentOrderResponse struct {
	PaymentOrder  *payment.PaymentOrder `json:"payment_order"`
	PaymentURL    string                `json:"payment_url"`
	QRCode        string                `json:"qr_code"`
	PaymentParams map[string]string     `json:"payment_params"`
}

// CreatePaymentOrder 创建支付订单
func (s *Service) CreatePaymentOrder(ctx context.Context, userID string, req CreatePaymentOrderRequest) (*CreatePaymentOrderResponse, error) {
	// 解析金额
	amount, err := decimal.NewFromString(req.Amount)
	if err != nil {
		return nil, fmt.Errorf("invalid amount: %w", err)
	}

	// 解析UUID
	orderUUID, err := uuid.Parse(req.OrderID)
	if err != nil {
		return nil, fmt.Errorf("invalid order ID: %w", err)
	}

	userUUID, err := uuid.Parse(userID)
	if err != nil {
		return nil, fmt.Errorf("invalid user ID: %w", err)
	}

	// 创建支付订单
	paymentOrder, err := s.paymentDS.CreatePaymentOrder(
		ctx,
		orderUUID,
		userUUID,
		amount,
		req.PaymentMethod,
		req.Subject,
		req.Description,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create payment order: %w", err)
	}

	response := &CreatePaymentOrderResponse{
		PaymentOrder: paymentOrder,
	}

	// 直接标记支付成功（测试模式）
	// 模拟支付成功回调
	_, err = s.paymentDS.ProcessPaymentCallback(
		ctx,
		paymentOrder.ID.String(),
		true, // 直接设置为支付成功
		"test_mock_payment_success",
	)
	if err != nil {
		return nil, fmt.Errorf("failed to mock payment success: %w", err)
	}

	// 设置模拟的支付参数
	response.PaymentURL = fmt.Sprintf("http://mock-payment-url.com/pay/%s", paymentOrder.ID.String())
	response.QRCode = "mock_qr_code_data"
	response.PaymentParams = map[string]string{
		"mock_mode":    "true",
		"auto_success": "true",
		"payment_id":   paymentOrder.ID.String(),
	}

	return response, nil
}

// GetPaymentOrder 获取支付订单
func (s *Service) GetPaymentOrder(ctx context.Context, paymentID string) (*payment.PaymentOrder, error) {
	paymentUUID, err := uuid.Parse(paymentID)
	if err != nil {
		return nil, fmt.Errorf("invalid payment ID: %w", err)
	}

	paymentOrder, err := s.paymentDS.ValidatePaymentStatus(ctx, paymentUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to get payment order: %w", err)
	}

	return paymentOrder, nil
}

// HandlePaymentCallbackRequest 处理支付回调请求
type HandlePaymentCallbackRequest struct {
	Provider string            `json:"provider"`
	Params   map[string]string `json:"params"`
}

// HandlePaymentCallbackResponse 处理支付回调响应
type HandlePaymentCallbackResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
}

// HandlePaymentCallback 处理支付回调
func (s *Service) HandlePaymentCallback(ctx context.Context, req HandlePaymentCallbackRequest) (*HandlePaymentCallbackResponse, error) {
	switch req.Provider {
	case "alipay":
		return s.handleAlipayCallback(ctx, req.Params)
	case "wechat":
		// TODO: 实现微信支付回调
		return nil, fmt.Errorf("wechat callback not implemented yet")
	default:
		return nil, fmt.Errorf("unsupported payment provider: %s", req.Provider)
	}
}

// handleAlipayCallback 处理支付宝回调
func (s *Service) handleAlipayCallback(ctx context.Context, params map[string]string) (*HandlePaymentCallbackResponse, error) {
	// 验证回调
	callbackResult, err := s.paymentClient.VerifyCallback(ctx, params)
	if err != nil {
		return &HandlePaymentCallbackResponse{
			Success: false,
			Message: fmt.Sprintf("failed to verify callback: %v", err),
		}, nil
	}

	if !callbackResult.IsValid {
		return &HandlePaymentCallbackResponse{
			Success: false,
			Message: callbackResult.ErrorMsg,
		}, nil
	}

	// 处理支付结果
	_, err = s.paymentDS.ProcessPaymentCallback(
		ctx,
		callbackResult.OutTradeNo,
		callbackResult.IsPaid,
		callbackResult.ThirdPartyTradeNo,
	)
	if err != nil {
		return &HandlePaymentCallbackResponse{
			Success: false,
			Message: fmt.Sprintf("failed to process callback: %v", err),
		}, nil
	}

	return &HandlePaymentCallbackResponse{
		Success: true,
		Message: "callback processed successfully",
	}, nil
}

// VerifyPaymentStatus 验证支付状态
func (s *Service) VerifyPaymentStatus(ctx context.Context, orderID string) (*payment.PaymentOrder, error) {
	orderUUID, err := uuid.Parse(orderID)
	if err != nil {
		return nil, fmt.Errorf("invalid order ID: %w", err)
	}

	paymentOrder, err := s.paymentDS.ValidatePaymentStatus(ctx, orderUUID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify payment status: %w", err)
	}

	// 如果支付订单状态为待支付，主动查询第三方支付状态
	if paymentOrder.Status == payment.PaymentStatusPending && paymentOrder.ThirdPartyOrderID != "" {
		switch paymentOrder.PaymentMethod {
		case payment.PaymentMethodAlipay:
			queryResult, err := s.paymentClient.QueryPayment(ctx, paymentOrder.ID.String())
			if err == nil && queryResult.IsPaid {
				// 更新支付状态
				_, err = s.paymentDS.ProcessPaymentCallback(
					ctx,
					paymentOrder.ID.String(),
					true,
					queryResult.ThirdPartyTradeNo,
				)
				if err != nil {
					return nil, fmt.Errorf("failed to update payment status: %w", err)
				}
				paymentOrder.Status = payment.PaymentStatusSuccess
			}
		}
	}

	return paymentOrder, nil
}
