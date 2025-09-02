package address

import (
	"context"
)

// Repository 地址仓储接口
type Repository interface {
	// Create 创建地址
	Create(ctx context.Context, address *Address) error

	// Update 更新地址
	Update(ctx context.Context, address *Address) error

	// Delete 删除地址
	Delete(ctx context.Context, id string) error

	// GetByID 根据ID获取地址
	GetByID(ctx context.Context, id string) (*Address, error)

	// GetByUserID 获取用户的所有地址
	GetByUserID(ctx context.Context, userID string) ([]*Address, error)

	// GetDefaultByUserID 获取用户默认地址
	GetDefaultByUserID(ctx context.Context, userID string) (*Address, error)

	// CountByUserID 统计用户地址数量
	CountByUserID(ctx context.Context, userID string) (int, error)

	// UnsetAllDefault 取消用户所有默认地址
	UnsetAllDefault(ctx context.Context, userID string) error

	// SetDefault 设置默认地址
	SetDefault(ctx context.Context, userID, addressID string) error
}
