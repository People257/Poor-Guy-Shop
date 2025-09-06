package payment

import (
	"context"

	"github.com/shopspring/decimal"
)

// PaymentClient 支付服务客户端接口
type PaymentClient interface {
	// CreatePayment 创建支付
	CreatePayment(ctx context.Context, req *CreatePaymentRequest) (*CreatePaymentResponse, error)

	// CreateQRPayment 创建扫码支付
	CreateQRPayment(ctx context.Context, req *CreatePaymentRequest) (*CreatePaymentResponse, error)

	// VerifyCallback 验证支付回调
	VerifyCallback(ctx context.Context, params map[string]string) (*CallbackResult, error)

	// QueryPayment 查询支付状态
	QueryPayment(ctx context.Context, outTradeNo string) (*QueryPaymentResult, error)

	// CreateRefund 创建退款
	CreateRefund(ctx context.Context, req *CreateRefundRequest) (*CreateRefundResponse, error)

	// QueryRefund 查询退款状态
	QueryRefund(ctx context.Context, outTradeNo, outRefundNo string) (*QueryRefundResult, error)
}

// CreatePaymentRequest 创建支付请求
type CreatePaymentRequest struct {
	OutTradeNo  string          `json:"out_trade_no"`
	Amount      decimal.Decimal `json:"amount"`
	Subject     string          `json:"subject"`
	Description string          `json:"description"`
	NotifyURL   string          `json:"notify_url"`
	ReturnURL   string          `json:"return_url"`
}

// CreatePaymentResponse 创建支付响应
type CreatePaymentResponse struct {
	PaymentURL    string            `json:"payment_url"`
	QRCode        string            `json:"qr_code"`
	PaymentParams map[string]string `json:"payment_params"`
}

// CallbackResult 支付回调结果
type CallbackResult struct {
	IsValid           bool              `json:"is_valid"`
	IsPaid            bool              `json:"is_paid"`
	OutTradeNo        string            `json:"out_trade_no"`
	ThirdPartyTradeNo string            `json:"third_party_trade_no"`
	Amount            string            `json:"amount"`
	TradeStatus       string            `json:"trade_status"`
	ErrorMsg          string            `json:"error_msg"`
	RawData           map[string]string `json:"raw_data"`
}

// QueryPaymentResult 查询支付结果
type QueryPaymentResult struct {
	IsPaid            bool        `json:"is_paid"`
	OutTradeNo        string      `json:"out_trade_no"`
	ThirdPartyTradeNo string      `json:"third_party_trade_no"`
	Amount            string      `json:"amount"`
	Status            string      `json:"status"`
	ErrorMsg          string      `json:"error_msg"`
	RawResult         interface{} `json:"raw_result"`
}

// CreateRefundRequest 创建退款请求
type CreateRefundRequest struct {
	OutTradeNo  string          `json:"out_trade_no"`
	OutRefundNo string          `json:"out_refund_no"`
	Amount      decimal.Decimal `json:"amount"`
	Reason      string          `json:"reason"`
}

// CreateRefundResponse 创建退款响应
type CreateRefundResponse struct {
	IsSuccess          bool        `json:"is_success"`
	OutRefundNo        string      `json:"out_refund_no"`
	ThirdPartyRefundNo string      `json:"third_party_refund_no"`
	RefundAmount       string      `json:"refund_amount"`
	Status             string      `json:"status"`
	ErrorMsg           string      `json:"error_msg"`
	RawResult          interface{} `json:"raw_result"`
}

// QueryRefundResult 查询退款结果
type QueryRefundResult struct {
	IsSuccess          bool        `json:"is_success"`
	OutRefundNo        string      `json:"out_refund_no"`
	ThirdPartyRefundNo string      `json:"third_party_refund_no"`
	RefundAmount       string      `json:"refund_amount"`
	Status             string      `json:"status"`
	ErrorMsg           string      `json:"error_msg"`
	RawResult          interface{} `json:"raw_result"`
}
