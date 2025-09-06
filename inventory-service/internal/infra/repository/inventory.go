package repository

import (
	"context"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/people257/poor-guy-shop/inventory-service/gen/gen/model"
	"github.com/people257/poor-guy-shop/inventory-service/gen/gen/query"
	"github.com/people257/poor-guy-shop/inventory-service/internal/domain/inventory"
)

// InventoryRepository 库存仓储实现
type InventoryRepository struct {
	db    *gorm.DB
	query *query.Query
}

// NewInventoryRepository 创建库存仓储
func NewInventoryRepository(db *gorm.DB, query *query.Query) *InventoryRepository {
	return &InventoryRepository{
		db:    db,
		query: query,
	}
}

// GetBySkuID 根据SKU ID获取库存
func (r *InventoryRepository) GetBySkuID(ctx context.Context, skuID uuid.UUID) (*inventory.Inventory, error) {
	q := r.query.Inventory
	inventoryModel, err := q.WithContext(ctx).Where(q.SkuID.Eq(skuID.String())).First()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, inventory.ErrInventoryNotFound
		}
		return nil, err
	}

	return r.modelToDomain(inventoryModel), nil
}

// BatchGetBySkuIDs 批量获取库存
func (r *InventoryRepository) BatchGetBySkuIDs(ctx context.Context, skuIDs []uuid.UUID) ([]*inventory.Inventory, error) {
	skuIDStrings := make([]string, len(skuIDs))
	for i, skuID := range skuIDs {
		skuIDStrings[i] = skuID.String()
	}

	q := r.query.Inventory
	inventoryModels, err := q.WithContext(ctx).Where(q.SkuID.In(skuIDStrings...)).Find()
	if err != nil {
		return nil, err
	}

	inventories := make([]*inventory.Inventory, len(inventoryModels))
	for i, model := range inventoryModels {
		inventories[i] = r.modelToDomain(model)
	}

	return inventories, nil
}

// Create 创建库存记录
func (r *InventoryRepository) Create(ctx context.Context, inv *inventory.Inventory) error {
	inventoryModel := r.domainToModel(inv)

	q := r.query.Inventory
	if err := q.WithContext(ctx).Create(inventoryModel); err != nil {
		return err
	}

	// 更新领域对象的ID
	if inventoryModel.ID != "" {
		parsedID, err := uuid.Parse(inventoryModel.ID)
		if err == nil {
			inv.ID = parsedID
		}
	}

	return nil
}

// Update 更新库存记录
func (r *InventoryRepository) Update(ctx context.Context, inv *inventory.Inventory) error {
	inventoryModel := r.domainToModel(inv)

	q := r.query.Inventory
	_, err := q.WithContext(ctx).Where(q.ID.Eq(inventoryModel.ID)).Updates(inventoryModel)
	return err
}

// UpdateWithVersion 乐观锁更新库存记录
func (r *InventoryRepository) UpdateWithVersion(ctx context.Context, inv *inventory.Inventory, version int32) error {
	inventoryModel := r.domainToModel(inv)

	q := r.query.Inventory
	result, err := q.WithContext(ctx).
		Where(q.ID.Eq(inventoryModel.ID), q.Version.Eq(version)).
		Updates(inventoryModel)

	if err != nil {
		return err
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound // 乐观锁冲突
	}

	return nil
}

// Delete 删除库存记录
func (r *InventoryRepository) Delete(ctx context.Context, skuID uuid.UUID) error {
	q := r.query.Inventory
	_, err := q.WithContext(ctx).Where(q.SkuID.Eq(skuID.String())).Delete()
	return err
}

// List 分页查询库存列表
func (r *InventoryRepository) List(ctx context.Context, offset, limit int) ([]*inventory.Inventory, int64, error) {
	q := r.query.Inventory

	// 查询总数
	total, err := q.WithContext(ctx).Count()
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	inventoryModels, err := q.WithContext(ctx).
		Offset(offset).
		Limit(limit).
		Order(q.UpdatedAt.Desc()).
		Find()
	if err != nil {
		return nil, 0, err
	}

	inventories := make([]*inventory.Inventory, len(inventoryModels))
	for i, model := range inventoryModels {
		inventories[i] = r.modelToDomain(model)
	}

	return inventories, total, nil
}

// ListLowStock 查询库存不足的商品
func (r *InventoryRepository) ListLowStock(ctx context.Context, offset, limit int) ([]*inventory.Inventory, int64, error) {
	// 使用原始SQL查询库存不足的商品
	var total int64
	err := r.db.WithContext(ctx).Model(&model.Inventory{}).
		Where("available_quantity <= alert_quantity AND available_quantity > 0").
		Count(&total).Error
	if err != nil {
		return nil, 0, err
	}

	var inventoryModels []*model.Inventory
	err = r.db.WithContext(ctx).
		Where("available_quantity <= alert_quantity AND available_quantity > 0").
		Offset(offset).
		Limit(limit).
		Order("available_quantity ASC").
		Find(&inventoryModels).Error
	if err != nil {
		return nil, 0, err
	}

	inventories := make([]*inventory.Inventory, len(inventoryModels))
	for i, model := range inventoryModels {
		inventories[i] = r.modelToDomain(model)
	}

	return inventories, total, nil
}

// ListOutOfStock 查询售罄的商品
func (r *InventoryRepository) ListOutOfStock(ctx context.Context, offset, limit int) ([]*inventory.Inventory, int64, error) {
	q := r.query.Inventory

	// 查询总数
	total, err := q.WithContext(ctx).Where(q.AvailableQuantity.Eq(0)).Count()
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	inventoryModels, err := q.WithContext(ctx).
		Where(q.AvailableQuantity.Eq(0)).
		Offset(offset).
		Limit(limit).
		Order(q.UpdatedAt.Desc()).
		Find()
	if err != nil {
		return nil, 0, err
	}

	inventories := make([]*inventory.Inventory, len(inventoryModels))
	for i, model := range inventoryModels {
		inventories[i] = r.modelToDomain(model)
	}

	return inventories, total, nil
}

// modelToDomain 将数据库模型转换为领域对象
func (r *InventoryRepository) modelToDomain(model *model.Inventory) *inventory.Inventory {
	id, _ := uuid.Parse(model.ID)
	skuID, _ := uuid.Parse(model.SkuID)

	return &inventory.Inventory{
		ID:                id,
		SkuID:             skuID,
		AvailableQuantity: model.AvailableQuantity,
		ReservedQuantity:  model.ReservedQuantity,
		TotalQuantity:     model.TotalQuantity,
		AlertQuantity:     model.AlertQuantity,
		CreatedAt:         model.CreatedAt,
		UpdatedAt:         model.UpdatedAt,
		Version:           model.Version,
	}
}

// domainToModel 将领域对象转换为数据库模型
func (r *InventoryRepository) domainToModel(inv *inventory.Inventory) *model.Inventory {
	return &model.Inventory{
		ID:                inv.ID.String(),
		SkuID:             inv.SkuID.String(),
		AvailableQuantity: inv.AvailableQuantity,
		ReservedQuantity:  inv.ReservedQuantity,
		TotalQuantity:     inv.TotalQuantity,
		AlertQuantity:     inv.AlertQuantity,
		CreatedAt:         inv.CreatedAt,
		UpdatedAt:         inv.UpdatedAt,
		Version:           inv.Version,
	}
}

// InventoryLogRepository 库存日志仓储实现
type InventoryLogRepository struct {
	db    *gorm.DB
	query *query.Query
}

// NewInventoryLogRepository 创建库存日志仓储
func NewInventoryLogRepository(db *gorm.DB, query *query.Query) *InventoryLogRepository {
	return &InventoryLogRepository{
		db:    db,
		query: query,
	}
}

// Create 创建库存变动日志
func (r *InventoryLogRepository) Create(ctx context.Context, log *inventory.InventoryLog) error {
	logModel := r.domainToLogModel(log)

	q := r.query.InventoryLog
	if err := q.WithContext(ctx).Create(logModel); err != nil {
		return err
	}

	// 更新领域对象的ID
	if logModel.ID != "" {
		parsedID, err := uuid.Parse(logModel.ID)
		if err == nil {
			log.ID = parsedID
		}
	}

	return nil
}

// GetBySkuID 根据SKU ID获取变动日志
func (r *InventoryLogRepository) GetBySkuID(ctx context.Context, skuID uuid.UUID, offset, limit int) ([]*inventory.InventoryLog, int64, error) {
	q := r.query.InventoryLog

	// 查询总数
	total, err := q.WithContext(ctx).Where(q.SkuID.Eq(skuID.String())).Count()
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	logModels, err := q.WithContext(ctx).
		Where(q.SkuID.Eq(skuID.String())).
		Offset(offset).
		Limit(limit).
		Order(q.CreatedAt.Desc()).
		Find()
	if err != nil {
		return nil, 0, err
	}

	logs := make([]*inventory.InventoryLog, len(logModels))
	for i, model := range logModels {
		logs[i] = r.logModelToDomain(model)
	}

	return logs, total, nil
}

// GetByOrderID 根据订单ID获取变动日志
func (r *InventoryLogRepository) GetByOrderID(ctx context.Context, orderID uuid.UUID) ([]*inventory.InventoryLog, error) {
	q := r.query.InventoryLog
	logModels, err := q.WithContext(ctx).
		Where(q.OrderID.Eq(orderID.String())).
		Order(q.CreatedAt.Desc()).
		Find()
	if err != nil {
		return nil, err
	}

	logs := make([]*inventory.InventoryLog, len(logModels))
	for i, model := range logModels {
		logs[i] = r.logModelToDomain(model)
	}

	return logs, nil
}

// GetByType 根据变动类型获取日志
func (r *InventoryLogRepository) GetByType(ctx context.Context, changeType inventory.InventoryChangeType, offset, limit int) ([]*inventory.InventoryLog, int64, error) {
	q := r.query.InventoryLog

	// 查询总数
	total, err := q.WithContext(ctx).Where(q.Type.Eq(string(changeType))).Count()
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	logModels, err := q.WithContext(ctx).
		Where(q.Type.Eq(string(changeType))).
		Offset(offset).
		Limit(limit).
		Order(q.CreatedAt.Desc()).
		Find()
	if err != nil {
		return nil, 0, err
	}

	logs := make([]*inventory.InventoryLog, len(logModels))
	for i, model := range logModels {
		logs[i] = r.logModelToDomain(model)
	}

	return logs, total, nil
}

// logModelToDomain 将日志数据库模型转换为领域对象
func (r *InventoryLogRepository) logModelToDomain(model *model.InventoryLog) *inventory.InventoryLog {
	id, _ := uuid.Parse(model.ID)
	skuID, _ := uuid.Parse(model.SkuID)

	var orderID *uuid.UUID
	if model.OrderID != nil {
		if parsedOrderID, err := uuid.Parse(*model.OrderID); err == nil {
			orderID = &parsedOrderID
		}
	}

	var operatorID *uuid.UUID
	if model.OperatorID != nil {
		if parsedOperatorID, err := uuid.Parse(*model.OperatorID); err == nil {
			operatorID = &parsedOperatorID
		}
	}

	var reason string
	if model.Reason != nil {
		reason = *model.Reason
	}

	return &inventory.InventoryLog{
		ID:             id,
		SkuID:          skuID,
		Type:           inventory.InventoryChangeType(model.Type),
		Quantity:       model.Quantity,
		BeforeQuantity: model.BeforeQuantity,
		AfterQuantity:  model.AfterQuantity,
		Reason:         reason,
		OrderID:        orderID,
		OperatorID:     operatorID,
		CreatedAt:      model.CreatedAt,
	}
}

// domainToLogModel 将日志领域对象转换为数据库模型
func (r *InventoryLogRepository) domainToLogModel(log *inventory.InventoryLog) *model.InventoryLog {
	var orderID *string
	if log.OrderID != nil {
		orderIDStr := log.OrderID.String()
		orderID = &orderIDStr
	}

	var operatorID *string
	if log.OperatorID != nil {
		operatorIDStr := log.OperatorID.String()
		operatorID = &operatorIDStr
	}

	var reason *string
	if log.Reason != "" {
		reason = &log.Reason
	}

	return &model.InventoryLog{
		ID:             log.ID.String(),
		SkuID:          log.SkuID.String(),
		Type:           string(log.Type),
		Quantity:       log.Quantity,
		BeforeQuantity: log.BeforeQuantity,
		AfterQuantity:  log.AfterQuantity,
		Reason:         reason,
		OrderID:        orderID,
		OperatorID:     operatorID,
		CreatedAt:      log.CreatedAt,
	}
}
