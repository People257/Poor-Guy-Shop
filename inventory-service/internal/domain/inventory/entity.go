package inventory

import (
	"time"

	"github.com/google/uuid"
)

// InventoryChangeType 库存变动类型
type InventoryChangeType string

const (
	InventoryChangeTypeIn      InventoryChangeType = "in"      // 入库
	InventoryChangeTypeOut     InventoryChangeType = "out"     // 出库
	InventoryChangeTypeReserve InventoryChangeType = "reserve" // 预占
	InventoryChangeTypeRelease InventoryChangeType = "release" // 释放预占
	InventoryChangeTypeAdjust  InventoryChangeType = "adjust"  // 库存调整
)

// InventoryStatus 库存状态
type InventoryStatus string

const (
	InventoryStatusNormal     InventoryStatus = "normal"       // 正常
	InventoryStatusLowStock   InventoryStatus = "low_stock"    // 库存不足
	InventoryStatusOutOfStock InventoryStatus = "out_of_stock" // 售罄
)

// ReservationStatus 预占状态
type ReservationStatus string

const (
	ReservationStatusReserved  ReservationStatus = "reserved"  // 预占中
	ReservationStatusConfirmed ReservationStatus = "confirmed" // 已确认
	ReservationStatusReleased  ReservationStatus = "released"  // 已释放
	ReservationStatusExpired   ReservationStatus = "expired"   // 已过期
)

// Inventory 库存实体
type Inventory struct {
	ID                uuid.UUID `json:"id"`
	SkuID             uuid.UUID `json:"sku_id"`
	AvailableQuantity int32     `json:"available_quantity"`
	ReservedQuantity  int32     `json:"reserved_quantity"`
	TotalQuantity     int32     `json:"total_quantity"`
	AlertQuantity     int32     `json:"alert_quantity"`
	CreatedAt         time.Time `json:"created_at"`
	UpdatedAt         time.Time `json:"updated_at"`
	Version           int32     `json:"version"`
}

// InventoryLog 库存变动日志实体
type InventoryLog struct {
	ID             uuid.UUID           `json:"id"`
	SkuID          uuid.UUID           `json:"sku_id"`
	Type           InventoryChangeType `json:"type"`
	Quantity       int32               `json:"quantity"`
	BeforeQuantity int32               `json:"before_quantity"`
	AfterQuantity  int32               `json:"after_quantity"`
	Reason         string              `json:"reason"`
	OrderID        *uuid.UUID          `json:"order_id"`
	OperatorID     *uuid.UUID          `json:"operator_id"`
	CreatedAt      time.Time           `json:"created_at"`
}

// InventoryReservation 库存预占记录实体
type InventoryReservation struct {
	ID          uuid.UUID         `json:"id"`
	SkuID       uuid.UUID         `json:"sku_id"`
	OrderID     uuid.UUID         `json:"order_id"`
	Quantity    int32             `json:"quantity"`
	Status      ReservationStatus `json:"status"`
	CreatedAt   time.Time         `json:"created_at"`
	UpdatedAt   time.Time         `json:"updated_at"`
	ExpiresAt   *time.Time        `json:"expires_at"`
	ConfirmedAt *time.Time        `json:"confirmed_at"`
	ReleasedAt  *time.Time        `json:"released_at"`
	Version     int32             `json:"version"`
}

// NewInventory 创建新的库存记录
func NewInventory(skuID uuid.UUID, totalQuantity, alertQuantity int32) *Inventory {
	now := time.Now()
	return &Inventory{
		ID:                uuid.New(),
		SkuID:             skuID,
		AvailableQuantity: totalQuantity,
		ReservedQuantity:  0,
		TotalQuantity:     totalQuantity,
		AlertQuantity:     alertQuantity,
		CreatedAt:         now,
		UpdatedAt:         now,
		Version:           1,
	}
}

// UpdateQuantity 更新库存数量
func (i *Inventory) UpdateQuantity(changeType InventoryChangeType, quantity int32) error {
	switch changeType {
	case InventoryChangeTypeIn:
		i.AvailableQuantity += quantity
		i.TotalQuantity += quantity
	case InventoryChangeTypeOut:
		if i.AvailableQuantity < quantity {
			return ErrInsufficientInventory
		}
		i.AvailableQuantity -= quantity
		i.TotalQuantity -= quantity
	case InventoryChangeTypeReserve:
		if i.AvailableQuantity < quantity {
			return ErrInsufficientInventory
		}
		i.AvailableQuantity -= quantity
		i.ReservedQuantity += quantity
	case InventoryChangeTypeRelease:
		if i.ReservedQuantity < quantity {
			return ErrInsufficientReservedInventory
		}
		i.AvailableQuantity += quantity
		i.ReservedQuantity -= quantity
	case InventoryChangeTypeAdjust:
		i.AvailableQuantity = quantity
		i.TotalQuantity = i.AvailableQuantity + i.ReservedQuantity
	}

	i.UpdatedAt = time.Now()
	i.Version++
	return nil
}

// GetStatus 获取库存状态
func (i *Inventory) GetStatus() InventoryStatus {
	if i.AvailableQuantity == 0 {
		return InventoryStatusOutOfStock
	}
	if i.AvailableQuantity <= i.AlertQuantity {
		return InventoryStatusLowStock
	}
	return InventoryStatusNormal
}

// IsLowStock 检查是否库存不足
func (i *Inventory) IsLowStock() bool {
	return i.AvailableQuantity <= i.AlertQuantity
}

// IsOutOfStock 检查是否售罄
func (i *Inventory) IsOutOfStock() bool {
	return i.AvailableQuantity == 0
}

// NewInventoryReservation 创建新的库存预占记录
func NewInventoryReservation(skuID, orderID uuid.UUID, quantity int32, expiresAt *time.Time) *InventoryReservation {
	now := time.Now()
	return &InventoryReservation{
		ID:        uuid.New(),
		SkuID:     skuID,
		OrderID:   orderID,
		Quantity:  quantity,
		Status:    ReservationStatusReserved,
		CreatedAt: now,
		UpdatedAt: now,
		ExpiresAt: expiresAt,
		Version:   1,
	}
}

// Confirm 确认预占（转为实际扣减）
func (r *InventoryReservation) Confirm() {
	now := time.Now()
	r.Status = ReservationStatusConfirmed
	r.ConfirmedAt = &now
	r.UpdatedAt = now
	r.Version++
}

// Release 释放预占
func (r *InventoryReservation) Release() {
	now := time.Now()
	r.Status = ReservationStatusReleased
	r.ReleasedAt = &now
	r.UpdatedAt = now
	r.Version++
}

// MarkExpired 标记为过期
func (r *InventoryReservation) MarkExpired() {
	now := time.Now()
	r.Status = ReservationStatusExpired
	r.ReleasedAt = &now
	r.UpdatedAt = now
	r.Version++
}

// IsExpired 检查是否过期
func (r *InventoryReservation) IsExpired() bool {
	if r.ExpiresAt == nil {
		return false
	}
	return time.Now().After(*r.ExpiresAt)
}

// CanConfirm 检查是否可以确认
func (r *InventoryReservation) CanConfirm() bool {
	return r.Status == ReservationStatusReserved && !r.IsExpired()
}

// CanRelease 检查是否可以释放
func (r *InventoryReservation) CanRelease() bool {
	return r.Status == ReservationStatusReserved
}

// NewInventoryLog 创建库存变动日志
func NewInventoryLog(skuID uuid.UUID, changeType InventoryChangeType, quantity, beforeQuantity, afterQuantity int32, reason string, orderID, operatorID *uuid.UUID) *InventoryLog {
	return &InventoryLog{
		ID:             uuid.New(),
		SkuID:          skuID,
		Type:           changeType,
		Quantity:       quantity,
		BeforeQuantity: beforeQuantity,
		AfterQuantity:  afterQuantity,
		Reason:         reason,
		OrderID:        orderID,
		OperatorID:     operatorID,
		CreatedAt:      time.Now(),
	}
}
