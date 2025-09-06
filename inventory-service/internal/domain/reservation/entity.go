package reservation

import (
	"context"

	"github.com/google/uuid"
	"github.com/people257/poor-guy-shop/inventory-service/internal/domain/inventory"
)

// Repository 预占记录仓储接口
type Repository interface {
	// Create 创建预占记录
	Create(ctx context.Context, reservation *inventory.InventoryReservation) error

	// Update 更新预占记录
	Update(ctx context.Context, reservation *inventory.InventoryReservation) error

	// UpdateWithVersion 乐观锁更新预占记录
	UpdateWithVersion(ctx context.Context, reservation *inventory.InventoryReservation, version int32) error

	// GetByID 根据ID获取预占记录
	GetByID(ctx context.Context, id uuid.UUID) (*inventory.InventoryReservation, error)

	// GetByOrderID 根据订单ID获取预占记录列表
	GetByOrderID(ctx context.Context, orderID uuid.UUID) ([]*inventory.InventoryReservation, error)

	// GetBySkuID 根据SKU ID获取预占记录列表
	GetBySkuID(ctx context.Context, skuID uuid.UUID, offset, limit int) ([]*inventory.InventoryReservation, int64, error)

	// GetExpiredReservations 获取过期的预占记录
	GetExpiredReservations(ctx context.Context, limit int) ([]*inventory.InventoryReservation, error)

	// GetByStatus 根据状态获取预占记录
	GetByStatus(ctx context.Context, status inventory.ReservationStatus, offset, limit int) ([]*inventory.InventoryReservation, int64, error)

	// Delete 删除预占记录
	Delete(ctx context.Context, id uuid.UUID) error
}
