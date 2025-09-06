package payment

import (
	"context"

	"github.com/google/uuid"
)

// Repository 支付订单仓储接口
type Repository interface {
	// Create 创建支付订单
	Create(ctx context.Context, paymentOrder *PaymentOrder) error

	// GetByID 根据ID获取支付订单
	GetByID(ctx context.Context, id uuid.UUID) (*PaymentOrder, error)

	// GetByOrderID 根据业务订单ID获取支付订单
	GetByOrderID(ctx context.Context, orderID uuid.UUID) (*PaymentOrder, error)

	// GetByThirdPartyOrderID 根据第三方订单ID获取支付订单
	GetByThirdPartyOrderID(ctx context.Context, thirdPartyOrderID string) (*PaymentOrder, error)

	// Update 更新支付订单
	Update(ctx context.Context, paymentOrder *PaymentOrder) error

	// List 分页查询支付订单
	List(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]*PaymentOrder, int64, error)

	// CreateLog 创建支付日志
	CreateLog(ctx context.Context, log *PaymentLog) error

	// GetLogs 获取支付日志
	GetLogs(ctx context.Context, paymentOrderID uuid.UUID) ([]*PaymentLog, error)
}
