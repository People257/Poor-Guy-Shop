package inventory

import (
	"context"

	"github.com/google/uuid"
)

// Repository 库存仓储接口
type Repository interface {
	// GetBySkuID 根据SKU ID获取库存
	GetBySkuID(ctx context.Context, skuID uuid.UUID) (*Inventory, error)

	// BatchGetBySkuIDs 批量获取库存
	BatchGetBySkuIDs(ctx context.Context, skuIDs []uuid.UUID) ([]*Inventory, error)

	// Create 创建库存记录
	Create(ctx context.Context, inventory *Inventory) error

	// Update 更新库存记录
	Update(ctx context.Context, inventory *Inventory) error

	// UpdateWithVersion 乐观锁更新库存记录
	UpdateWithVersion(ctx context.Context, inventory *Inventory, version int32) error

	// Delete 删除库存记录
	Delete(ctx context.Context, skuID uuid.UUID) error

	// List 分页查询库存列表
	List(ctx context.Context, offset, limit int) ([]*Inventory, int64, error)

	// ListLowStock 查询库存不足的商品
	ListLowStock(ctx context.Context, offset, limit int) ([]*Inventory, int64, error)

	// ListOutOfStock 查询售罄的商品
	ListOutOfStock(ctx context.Context, offset, limit int) ([]*Inventory, int64, error)
}

// LogRepository 库存日志仓储接口
type LogRepository interface {
	// Create 创建库存变动日志
	Create(ctx context.Context, log *InventoryLog) error

	// GetBySkuID 根据SKU ID获取变动日志
	GetBySkuID(ctx context.Context, skuID uuid.UUID, offset, limit int) ([]*InventoryLog, int64, error)

	// GetByOrderID 根据订单ID获取变动日志
	GetByOrderID(ctx context.Context, orderID uuid.UUID) ([]*InventoryLog, error)

	// GetByType 根据变动类型获取日志
	GetByType(ctx context.Context, changeType InventoryChangeType, offset, limit int) ([]*InventoryLog, int64, error)
}
