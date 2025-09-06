package payment

import (
	"context"
	"fmt"
	"net/url"

	"github.com/smartwalle/alipay/v3"
)

// AlipayConfig 支付宝配置
type AlipayConfig struct {
	AppID      string `mapstructure:"app_id"`
	PrivateKey string `mapstructure:"private_key"`
	PublicKey  string `mapstructure:"public_key"`
	IsSandbox  bool   `mapstructure:"is_sandbox"`
}

// AlipayService 支付宝支付服务
type AlipayService struct {
	client *alipay.Client
	config *AlipayConfig
}

// NewAlipayClient 创建支付宝支付服务
func NewAlipayClient(config *AlipayConfig) PaymentClient {
	service, err := NewAlipayService(config)
	if err != nil {
		panic(fmt.Sprintf("failed to create alipay client: %v", err))
	}
	return service
}

// NewAlipayService 创建支付宝支付服务
func NewAlipayService(config *AlipayConfig) (*AlipayService, error) {
	var client *alipay.Client
	var err error

	if config.IsSandbox {
		client, err = alipay.New(config.AppID, config.PrivateKey, true)
	} else {
		client, err = alipay.New(config.AppID, config.PrivateKey, false)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to create alipay client: %w", err)
	}

	// 加载支付宝公钥
	if err := client.LoadAliPayPublicKey(config.PublicKey); err != nil {
		return nil, fmt.Errorf("failed to load alipay public key: %w", err)
	}

	return &AlipayService{
		client: client,
		config: config,
	}, nil
}

// CreatePayment 创建支付
func (s *AlipayService) CreatePayment(ctx context.Context, req *CreatePaymentRequest) (*CreatePaymentResponse, error) {
	// 创建支付宝支付请求
	pay := alipay.TradePagePay{
		Trade: alipay.Trade{
			Subject:     req.Subject,
			OutTradeNo:  req.OutTradeNo,
			TotalAmount: req.Amount.String(),
			ProductCode: "FAST_INSTANT_TRADE_PAY",
		},
	}

	if req.Description != "" {
		pay.Body = req.Description
	}
	if req.ReturnURL != "" {
		pay.ReturnURL = req.ReturnURL
	}
	if req.NotifyURL != "" {
		pay.NotifyURL = req.NotifyURL
	}

	// 生成支付URL
	payURL, err := s.client.TradePagePay(pay)
	if err != nil {
		return nil, fmt.Errorf("failed to create alipay payment: %w", err)
	}

	return &CreatePaymentResponse{
		PaymentURL: payURL.String(),
		QRCode:     "", // 网页支付不需要二维码
		PaymentParams: map[string]string{
			"payment_url": payURL.String(),
		},
	}, nil
}

// CreateQRPayment 创建扫码支付
func (s *AlipayService) CreateQRPayment(ctx context.Context, req *CreatePaymentRequest) (*CreatePaymentResponse, error) {
	// 创建支付宝扫码支付请求
	pay := alipay.TradePreCreate{
		Trade: alipay.Trade{
			Subject:     req.Subject,
			OutTradeNo:  req.OutTradeNo,
			TotalAmount: req.Amount.String(),
		},
	}

	if req.Description != "" {
		pay.Body = req.Description
	}
	if req.NotifyURL != "" {
		pay.NotifyURL = req.NotifyURL
	}

	// 创建预支付订单
	result, err := s.client.TradePreCreate(ctx, pay)
	if err != nil {
		return nil, fmt.Errorf("failed to create alipay qr payment: %w", err)
	}

	if !result.IsSuccess() {
		return nil, fmt.Errorf("alipay qr payment failed: %s", result.Msg)
	}

	return &CreatePaymentResponse{
		PaymentURL: "",
		QRCode:     result.QRCode,
		PaymentParams: map[string]string{
			"qr_code": result.QRCode,
		},
	}, nil
}

// VerifyCallback 验证支付回调
func (s *AlipayService) VerifyCallback(ctx context.Context, params map[string]string) (*CallbackResult, error) {
	// 转换为url.Values
	values := make(url.Values)
	for k, v := range params {
		values.Set(k, v)
	}

	// 验证回调签名
	if err := s.client.VerifySign(values); err != nil {
		return &CallbackResult{
			IsValid:  false,
			ErrorMsg: fmt.Sprintf("invalid signature: %v", err),
		}, nil
	}

	// 获取交易状态
	tradeStatus := params["trade_status"]
	outTradeNo := params["out_trade_no"]
	tradeNo := params["trade_no"]
	totalAmount := params["total_amount"]

	result := &CallbackResult{
		IsValid:           true,
		OutTradeNo:        outTradeNo,
		ThirdPartyTradeNo: tradeNo,
		Amount:            totalAmount,
		TradeStatus:       tradeStatus,
		RawData:           params,
	}

	// 判断支付是否成功
	switch tradeStatus {
	case "TRADE_SUCCESS", "TRADE_FINISHED":
		result.IsPaid = true
	case "TRADE_CLOSED":
		result.IsPaid = false
		result.ErrorMsg = "trade closed"
	default:
		result.IsPaid = false
		result.ErrorMsg = fmt.Sprintf("unknown trade status: %s", tradeStatus)
	}

	return result, nil
}

// QueryPayment 查询支付状态
func (s *AlipayService) QueryPayment(ctx context.Context, outTradeNo string) (*QueryPaymentResult, error) {
	query := alipay.TradeQuery{
		OutTradeNo: outTradeNo,
	}

	result, err := s.client.TradeQuery(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query alipay payment: %w", err)
	}

	if !result.IsSuccess() {
		return &QueryPaymentResult{
			Status:    "UNKNOWN",
			ErrorMsg:  result.Msg,
			RawResult: result,
		}, nil
	}

	queryResult := &QueryPaymentResult{
		OutTradeNo:        result.OutTradeNo,
		ThirdPartyTradeNo: result.TradeNo,
		Amount:            result.TotalAmount,
		Status:            string(result.TradeStatus),
		RawResult:         result,
	}

	// 判断支付状态
	switch result.TradeStatus {
	case alipay.TradeStatusSuccess, alipay.TradeStatusFinished:
		queryResult.IsPaid = true
	case alipay.TradeStatusWaitBuyerPay:
		queryResult.IsPaid = false
		queryResult.Status = "PENDING"
	case alipay.TradeStatusClosed:
		queryResult.IsPaid = false
		queryResult.Status = "CLOSED"
	default:
		queryResult.IsPaid = false
		queryResult.Status = "UNKNOWN"
	}

	return queryResult, nil
}

// CreateRefund 创建退款
func (s *AlipayService) CreateRefund(ctx context.Context, req *CreateRefundRequest) (*CreateRefundResponse, error) {
	refund := alipay.TradeRefund{
		OutTradeNo:   req.OutTradeNo,
		RefundAmount: req.Amount.String(),
		RefundReason: req.Reason,
		OutRequestNo: req.OutRefundNo,
	}

	result, err := s.client.TradeRefund(ctx, refund)
	if err != nil {
		return nil, fmt.Errorf("failed to create alipay refund: %w", err)
	}

	response := &CreateRefundResponse{
		OutRefundNo:        req.OutRefundNo,
		ThirdPartyRefundNo: result.TradeNo,
		RefundAmount:       req.Amount.String(),
		RawResult:          result,
	}

	if result.IsSuccess() {
		response.IsSuccess = true
		response.Status = "SUCCESS"
	} else {
		response.IsSuccess = false
		response.Status = "FAILED"
		response.ErrorMsg = result.Msg
	}

	return response, nil
}

// QueryRefund 查询退款状态
func (s *AlipayService) QueryRefund(ctx context.Context, outTradeNo, outRefundNo string) (*QueryRefundResult, error) {
	query := alipay.TradeFastPayRefundQuery{
		OutTradeNo:   outTradeNo,
		OutRequestNo: outRefundNo,
	}

	result, err := s.client.TradeFastPayRefundQuery(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("failed to query alipay refund: %w", err)
	}

	queryResult := &QueryRefundResult{
		OutRefundNo:        outRefundNo,
		ThirdPartyRefundNo: result.TradeNo,
		RefundAmount:       result.RefundAmount,
		RawResult:          result,
	}

	if result.IsSuccess() {
		queryResult.IsSuccess = true
		queryResult.Status = "SUCCESS"
	} else {
		queryResult.IsSuccess = false
		queryResult.Status = "FAILED"
		queryResult.ErrorMsg = result.Msg
	}

	return queryResult, nil
}
