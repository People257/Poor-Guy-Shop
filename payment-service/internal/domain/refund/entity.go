package refund

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// RefundStatus 退款状态
type RefundStatus string

const (
	RefundStatusPending RefundStatus = "pending"
	RefundStatusSuccess RefundStatus = "success"
	RefundStatusFailed  RefundStatus = "failed"
)

// Refund 退款实体
type Refund struct {
	ID                 uuid.UUID       `json:"id"`
	PaymentOrderID     uuid.UUID       `json:"payment_order_id"`
	Amount             decimal.Decimal `json:"amount"`
	Reason             string          `json:"reason"`
	Status             RefundStatus    `json:"status"`
	ThirdPartyRefundID string          `json:"third_party_refund_id"`
	ThirdPartyResponse string          `json:"third_party_response"`
	CreatedAt          time.Time       `json:"created_at"`
	UpdatedAt          time.Time       `json:"updated_at"`
	ProcessedAt        *time.Time      `json:"processed_at"`
}

// NewRefund 创建新的退款记录
func NewRefund(paymentOrderID uuid.UUID, amount decimal.Decimal, reason string) *Refund {
	now := time.Now()
	return &Refund{
		ID:             uuid.New(),
		PaymentOrderID: paymentOrderID,
		Amount:         amount,
		Reason:         reason,
		Status:         RefundStatusPending,
		CreatedAt:      now,
		UpdatedAt:      now,
	}
}

// MarkAsSuccess 标记退款成功
func (r *Refund) MarkAsSuccess(thirdPartyRefundID, thirdPartyResponse string) {
	now := time.Now()
	r.Status = RefundStatusSuccess
	r.ThirdPartyRefundID = thirdPartyRefundID
	r.ThirdPartyResponse = thirdPartyResponse
	r.ProcessedAt = &now
	r.UpdatedAt = now
}

// MarkAsFailed 标记退款失败
func (r *Refund) MarkAsFailed(errorMessage string) {
	r.Status = RefundStatusFailed
	r.ThirdPartyResponse = errorMessage
	r.UpdatedAt = time.Now()
}

// IsProcessed 检查是否已处理
func (r *Refund) IsProcessed() bool {
	return r.Status == RefundStatusSuccess || r.Status == RefundStatusFailed
}
