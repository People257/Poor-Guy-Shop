package inventory

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/people257/poor-guy-shop/inventory-service/internal/domain/inventory"
)

// Service 库存应用服务
type Service struct {
	inventoryDomain *inventory.DomainService
	inventoryRepo   inventory.Repository
	logRepo         inventory.LogRepository
}

// NewService 创建库存应用服务
func NewService(inventoryDomain *inventory.DomainService, inventoryRepo inventory.Repository, logRepo inventory.LogRepository) *Service {
	return &Service{
		inventoryDomain: inventoryDomain,
		inventoryRepo:   inventoryRepo,
		logRepo:         logRepo,
	}
}

// GetInventory 获取库存信息
func (s *Service) GetInventory(ctx context.Context, skuID uuid.UUID) (*inventory.Inventory, error) {
	return s.inventoryRepo.GetBySkuID(ctx, skuID)
}

// BatchGetInventory 批量获取库存信息
func (s *Service) BatchGetInventory(ctx context.Context, skuIDs []uuid.UUID) ([]*inventory.Inventory, error) {
	return s.inventoryRepo.BatchGetBySkuIDs(ctx, skuIDs)
}

// CreateInventory 创建库存记录
func (s *Service) CreateInventory(ctx context.Context, skuID uuid.UUID, totalQuantity, alertQuantity int32, operatorID *uuid.UUID) (*inventory.Inventory, error) {
	return s.inventoryDomain.CreateInventory(ctx, skuID, totalQuantity, alertQuantity, operatorID)
}

// UpdateInventoryQuantity 更新库存数量
func (s *Service) UpdateInventoryQuantity(ctx context.Context, skuID uuid.UUID, changeType inventory.InventoryChangeType, quantity int32, reason string, orderID, operatorID *uuid.UUID) (*inventory.Inventory, error) {
	return s.inventoryDomain.UpdateInventoryQuantity(ctx, skuID, changeType, quantity, reason, orderID, operatorID)
}

// IncrementInventory 增加库存（入库）
func (s *Service) IncrementInventory(ctx context.Context, skuID uuid.UUID, quantity int32, reason string, operatorID *uuid.UUID) (*inventory.Inventory, error) {
	return s.inventoryDomain.UpdateInventoryQuantity(ctx, skuID, inventory.InventoryChangeTypeIn, quantity, reason, nil, operatorID)
}

// DecrementInventory 减少库存（出库）
func (s *Service) DecrementInventory(ctx context.Context, skuID uuid.UUID, quantity int32, reason string, operatorID *uuid.UUID) (*inventory.Inventory, error) {
	return s.inventoryDomain.UpdateInventoryQuantity(ctx, skuID, inventory.InventoryChangeTypeOut, quantity, reason, nil, operatorID)
}

// AdjustInventory 调整库存
func (s *Service) AdjustInventory(ctx context.Context, skuID uuid.UUID, newQuantity int32, reason string, operatorID *uuid.UUID) (*inventory.Inventory, error) {
	return s.inventoryDomain.UpdateInventoryQuantity(ctx, skuID, inventory.InventoryChangeTypeAdjust, newQuantity, reason, nil, operatorID)
}

// CheckInventoryAvailability 检查库存可用性
func (s *Service) CheckInventoryAvailability(ctx context.Context, items []inventory.ReserveItem) (bool, []uuid.UUID, error) {
	return s.inventoryDomain.CheckInventoryAvailability(ctx, items)
}

// ReserveInventory 预占库存
func (s *Service) ReserveInventory(ctx context.Context, orderID uuid.UUID, items []inventory.ReserveItem, expiresAt *time.Time) ([]*inventory.InventoryReservation, error) {
	return s.inventoryDomain.BatchReserveInventory(ctx, orderID, items, expiresAt)
}

// GetInventoryLogs 获取库存变动日志
func (s *Service) GetInventoryLogs(ctx context.Context, skuID uuid.UUID, page, pageSize int) ([]*inventory.InventoryLog, int64, error) {
	offset := (page - 1) * pageSize
	return s.logRepo.GetBySkuID(ctx, skuID, offset, pageSize)
}

// GetInventoryLogsByOrderID 根据订单ID获取库存变动日志
func (s *Service) GetInventoryLogsByOrderID(ctx context.Context, orderID uuid.UUID) ([]*inventory.InventoryLog, error) {
	return s.logRepo.GetByOrderID(ctx, orderID)
}

// ListInventory 分页查询库存列表
func (s *Service) ListInventory(ctx context.Context, page, pageSize int) ([]*inventory.Inventory, int64, error) {
	offset := (page - 1) * pageSize
	return s.inventoryRepo.List(ctx, offset, pageSize)
}

// ListLowStockInventory 查询库存不足的商品
func (s *Service) ListLowStockInventory(ctx context.Context, page, pageSize int) ([]*inventory.Inventory, int64, error) {
	offset := (page - 1) * pageSize
	return s.inventoryRepo.ListLowStock(ctx, offset, pageSize)
}

// ListOutOfStockInventory 查询售罄的商品
func (s *Service) ListOutOfStockInventory(ctx context.Context, page, pageSize int) ([]*inventory.Inventory, int64, error) {
	offset := (page - 1) * pageSize
	return s.inventoryRepo.ListOutOfStock(ctx, offset, pageSize)
}

// DeleteInventory 删除库存记录
func (s *Service) DeleteInventory(ctx context.Context, skuID uuid.UUID) error {
	return s.inventoryRepo.Delete(ctx, skuID)
}
