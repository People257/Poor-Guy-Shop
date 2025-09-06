package inventory

import (
	"context"
	"errors"
	"time"

	"github.com/google/uuid"
)

// DomainService 库存领域服务
type DomainService struct {
	inventoryRepo Repository
	logRepo       LogRepository
}

// NewDomainService 创建库存领域服务
func NewDomainService(inventoryRepo Repository, logRepo LogRepository) *DomainService {
	return &DomainService{
		inventoryRepo: inventoryRepo,
		logRepo:       logRepo,
	}
}

// CreateInventory 创建库存记录
func (s *DomainService) CreateInventory(ctx context.Context, skuID uuid.UUID, totalQuantity, alertQuantity int32, operatorID *uuid.UUID) (*Inventory, error) {
	// 检查SKU是否已存在库存记录
	existingInventory, err := s.inventoryRepo.GetBySkuID(ctx, skuID)
	if err == nil && existingInventory != nil {
		return nil, errors.New("inventory already exists for this SKU")
	}

	// 创建新的库存记录
	inventory := NewInventory(skuID, totalQuantity, alertQuantity)

	if err := s.inventoryRepo.Create(ctx, inventory); err != nil {
		return nil, err
	}

	// 记录库存变动日志
	log := NewInventoryLog(
		skuID,
		InventoryChangeTypeIn,
		totalQuantity,
		0,
		totalQuantity,
		"初始库存",
		nil,
		operatorID,
	)

	if err := s.logRepo.Create(ctx, log); err != nil {
		// 日志记录失败不影响主流程，但应该记录错误
		// TODO: 添加日志记录
	}

	return inventory, nil
}

// UpdateInventoryQuantity 更新库存数量
func (s *DomainService) UpdateInventoryQuantity(ctx context.Context, skuID uuid.UUID, changeType InventoryChangeType, quantity int32, reason string, orderID, operatorID *uuid.UUID) (*Inventory, error) {
	if quantity <= 0 {
		return nil, ErrInvalidQuantity
	}

	// 获取当前库存
	inventory, err := s.inventoryRepo.GetBySkuID(ctx, skuID)
	if err != nil {
		return nil, err
	}
	if inventory == nil {
		return nil, ErrInventoryNotFound
	}

	// 记录变动前的数量
	var beforeQuantity int32
	switch changeType {
	case InventoryChangeTypeIn, InventoryChangeTypeOut, InventoryChangeTypeAdjust:
		beforeQuantity = inventory.AvailableQuantity
	case InventoryChangeTypeReserve, InventoryChangeTypeRelease:
		beforeQuantity = inventory.AvailableQuantity
	}

	// 更新库存数量
	if err := inventory.UpdateQuantity(changeType, quantity); err != nil {
		return nil, err
	}

	// 使用乐观锁更新库存
	if err := s.inventoryRepo.UpdateWithVersion(ctx, inventory, inventory.Version-1); err != nil {
		return nil, err
	}

	// 记录库存变动日志
	var logQuantity int32
	switch changeType {
	case InventoryChangeTypeIn, InventoryChangeTypeAdjust:
		logQuantity = quantity
	case InventoryChangeTypeOut, InventoryChangeTypeReserve:
		logQuantity = -quantity
	case InventoryChangeTypeRelease:
		logQuantity = quantity
	}

	log := NewInventoryLog(
		skuID,
		changeType,
		logQuantity,
		beforeQuantity,
		inventory.AvailableQuantity,
		reason,
		orderID,
		operatorID,
	)

	if err := s.logRepo.Create(ctx, log); err != nil {
		// 日志记录失败不影响主流程，但应该记录错误
		// TODO: 添加日志记录
	}

	return inventory, nil
}

// CheckInventoryAvailability 检查库存可用性
func (s *DomainService) CheckInventoryAvailability(ctx context.Context, items []ReserveItem) (bool, []uuid.UUID, error) {
	if len(items) == 0 {
		return true, nil, nil
	}

	// 提取所有SKU ID
	skuIDs := make([]uuid.UUID, len(items))
	itemMap := make(map[uuid.UUID]int32)

	for i, item := range items {
		skuIDs[i] = item.SkuID
		itemMap[item.SkuID] = item.Quantity
	}

	// 批量获取库存信息
	inventories, err := s.inventoryRepo.BatchGetBySkuIDs(ctx, skuIDs)
	if err != nil {
		return false, nil, err
	}

	// 检查每个SKU的库存是否充足
	inventoryMap := make(map[uuid.UUID]*Inventory)
	for _, inventory := range inventories {
		inventoryMap[inventory.SkuID] = inventory
	}

	var insufficientSkus []uuid.UUID
	for skuID, requiredQuantity := range itemMap {
		inventory, exists := inventoryMap[skuID]
		if !exists {
			insufficientSkus = append(insufficientSkus, skuID)
			continue
		}

		if inventory.AvailableQuantity < requiredQuantity {
			insufficientSkus = append(insufficientSkus, skuID)
		}
	}

	return len(insufficientSkus) == 0, insufficientSkus, nil
}

// ReserveItem 预占商品项
type ReserveItem struct {
	SkuID    uuid.UUID
	Quantity int32
}

// BatchReserveInventory 批量预占库存
func (s *DomainService) BatchReserveInventory(ctx context.Context, orderID uuid.UUID, items []ReserveItem, expiresAt *time.Time) ([]*InventoryReservation, error) {
	if len(items) == 0 {
		return nil, ErrInvalidQuantity
	}

	// 首先检查所有库存是否充足
	available, _, err := s.CheckInventoryAvailability(ctx, items)
	if err != nil {
		return nil, err
	}
	if !available {
		return nil, ErrInsufficientInventory
	}

	var reservations []*InventoryReservation

	// 为每个商品创建预占记录并更新库存
	for _, item := range items {
		// 更新库存（预占）
		_, err := s.UpdateInventoryQuantity(
			ctx,
			item.SkuID,
			InventoryChangeTypeReserve,
			item.Quantity,
			"订单预占",
			&orderID,
			nil,
		)
		if err != nil {
			// TODO: 需要回滚已经预占的库存
			return nil, err
		}

		// 创建预占记录
		reservation := NewInventoryReservation(item.SkuID, orderID, item.Quantity, expiresAt)
		reservations = append(reservations, reservation)
	}

	return reservations, nil
}

// ReleaseReservation 释放预占
func (s *DomainService) ReleaseReservation(ctx context.Context, reservation *InventoryReservation) error {
	if !reservation.CanRelease() {
		return ErrReservationAlreadyReleased
	}

	// 释放库存
	_, err := s.UpdateInventoryQuantity(
		ctx,
		reservation.SkuID,
		InventoryChangeTypeRelease,
		reservation.Quantity,
		"释放预占",
		&reservation.OrderID,
		nil,
	)
	if err != nil {
		return err
	}

	// 更新预占记录状态
	reservation.Release()

	return nil
}

// ConfirmReservation 确认预占（实际扣减库存）
func (s *DomainService) ConfirmReservation(ctx context.Context, reservation *InventoryReservation) error {
	if !reservation.CanConfirm() {
		if reservation.IsExpired() {
			return ErrReservationExpired
		}
		return ErrReservationAlreadyConfirmed
	}

	// 确认预占不需要再次更新库存数量，因为预占时已经从可用库存中扣减了
	// 只需要从预占库存转移到实际扣减
	inventory, err := s.inventoryRepo.GetBySkuID(ctx, reservation.SkuID)
	if err != nil {
		return err
	}

	if inventory.ReservedQuantity < reservation.Quantity {
		return ErrInsufficientReservedInventory
	}

	// 从预占库存中扣减，总库存也要相应减少
	inventory.ReservedQuantity -= reservation.Quantity
	inventory.TotalQuantity -= reservation.Quantity
	inventory.UpdatedAt = time.Now()
	inventory.Version++

	if err := s.inventoryRepo.UpdateWithVersion(ctx, inventory, inventory.Version-1); err != nil {
		return err
	}

	// 记录库存变动日志
	log := NewInventoryLog(
		reservation.SkuID,
		InventoryChangeTypeOut,
		-reservation.Quantity,
		inventory.TotalQuantity+reservation.Quantity,
		inventory.TotalQuantity,
		"确认扣减",
		&reservation.OrderID,
		nil,
	)

	if err := s.logRepo.Create(ctx, log); err != nil {
		// 日志记录失败不影响主流程
	}

	// 更新预占记录状态
	reservation.Confirm()

	return nil
}
