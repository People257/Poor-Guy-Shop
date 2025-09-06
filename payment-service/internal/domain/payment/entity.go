package payment

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// PaymentMethod 支付方式
type PaymentMethod string

const (
	PaymentMethodAlipay   PaymentMethod = "alipay"
	PaymentMethodWechat   PaymentMethod = "wechat"
	PaymentMethodBankCard PaymentMethod = "bank_card"
	PaymentMethodBalance  PaymentMethod = "balance"
)

// PaymentStatus 支付状态
type PaymentStatus string

const (
	PaymentStatusPending         PaymentStatus = "pending"
	PaymentStatusSuccess         PaymentStatus = "success"
	PaymentStatusFailed          PaymentStatus = "failed"
	PaymentStatusCancelled       PaymentStatus = "cancelled"
	PaymentStatusRefunded        PaymentStatus = "refunded"
	PaymentStatusPartialRefunded PaymentStatus = "partial_refunded"
)

// PaymentOrder 支付订单实体
type PaymentOrder struct {
	ID                 uuid.UUID       `json:"id"`
	OrderID            uuid.UUID       `json:"order_id"`
	UserID             uuid.UUID       `json:"user_id"`
	Amount             decimal.Decimal `json:"amount"`
	PaymentMethod      PaymentMethod   `json:"payment_method"`
	Status             PaymentStatus   `json:"status"`
	ThirdPartyOrderID  string          `json:"third_party_order_id"`
	ThirdPartyResponse string          `json:"third_party_response"`
	Subject            string          `json:"subject"`
	Description        string          `json:"description"`
	NotifyURL          string          `json:"notify_url"`
	ReturnURL          string          `json:"return_url"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
	PaidAt             *time.Time      `json:"paid_at"`
	ExpiredAt          *time.Time      `json:"expired_at"`
}

// PaymentLog 支付日志实体
type PaymentLog struct {
	ID             uuid.UUID `json:"id"`
	PaymentOrderID uuid.UUID `json:"payment_order_id"`
	Action         string    `json:"action"`
	RequestData    string    `json:"request_data"`
	ResponseData   string    `json:"response_data"`
	ErrorMessage   string    `json:"error_message"`
	CreatedAt      time.Time `json:"created_at"`
}

// NewPaymentOrder 创建新的支付订单
func NewPaymentOrder(orderID, userID uuid.UUID, amount decimal.Decimal, paymentMethod PaymentMethod, subject, description string) *PaymentOrder {
	now := time.Now()
	expiredAt := now.Add(30 * time.Minute) // 30分钟过期

	return &PaymentOrder{
		ID:            uuid.New(),
		OrderID:       orderID,
		UserID:        userID,
		Amount:        amount,
		PaymentMethod: paymentMethod,
		Status:        PaymentStatusPending,
		Subject:       subject,
		Description:   description,
		CreatedAt:     now,
		UpdatedAt:     now,
		ExpiredAt:     &expiredAt,
	}
}

// MarkAsPaid 标记为已支付
func (p *PaymentOrder) MarkAsPaid(thirdPartyOrderID, thirdPartyResponse string) {
	now := time.Now()
	p.Status = PaymentStatusSuccess
	p.ThirdPartyOrderID = thirdPartyOrderID
	p.ThirdPartyResponse = thirdPartyResponse
	p.PaidAt = &now
	p.UpdatedAt = now
}

// MarkAsFailed 标记为支付失败
func (p *PaymentOrder) MarkAsFailed(errorMessage string) {
	p.Status = PaymentStatusFailed
	p.ThirdPartyResponse = errorMessage
	p.UpdatedAt = time.Now()
}

// MarkAsCancelled 标记为已取消
func (p *PaymentOrder) MarkAsCancelled() {
	p.Status = PaymentStatusCancelled
	p.UpdatedAt = time.Now()
}

// MarkAsRefunded 标记为已退款
func (p *PaymentOrder) MarkAsRefunded() {
	p.Status = PaymentStatusRefunded
	p.UpdatedAt = time.Now()
}

// MarkAsPartialRefunded 标记为部分退款
func (p *PaymentOrder) MarkAsPartialRefunded() {
	p.Status = PaymentStatusPartialRefunded
	p.UpdatedAt = time.Now()
}

// IsExpired 检查是否过期
func (p *PaymentOrder) IsExpired() bool {
	if p.ExpiredAt == nil {
		return false
	}
	return time.Now().After(*p.ExpiredAt)
}

// CanRefund 检查是否可以退款
func (p *PaymentOrder) CanRefund() bool {
	return p.Status == PaymentStatusSuccess || p.Status == PaymentStatusPartialRefunded
}

// NewPaymentLog 创建支付日志
func NewPaymentLog(paymentOrderID uuid.UUID, action, requestData, responseData, errorMessage string) *PaymentLog {
	return &PaymentLog{
		ID:             uuid.New(),
		PaymentOrderID: paymentOrderID,
		Action:         action,
		RequestData:    requestData,
		ResponseData:   responseData,
		ErrorMessage:   errorMessage,
		CreatedAt:      time.Now(),
	}
}
