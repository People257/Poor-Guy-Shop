package refund

import (
	"context"

	"github.com/google/uuid"
)

// Repository 退款仓储接口
type Repository interface {
	// Create 创建退款记录
	Create(ctx context.Context, refund *Refund) error

	// GetByID 根据ID获取退款记录
	GetByID(ctx context.Context, id uuid.UUID) (*Refund, error)

	// GetByPaymentOrderID 根据支付订单ID获取退款记录
	GetByPaymentOrderID(ctx context.Context, paymentOrderID uuid.UUID) ([]*Refund, error)

	// Update 更新退款记录
	Update(ctx context.Context, refund *Refund) error

	// List 分页查询退款记录
	List(ctx context.Context, page, pageSize int) ([]*Refund, int64, error)
}
