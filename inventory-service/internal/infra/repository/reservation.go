package repository

import (
	"context"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/people257/poor-guy-shop/inventory-service/gen/gen/model"
	"github.com/people257/poor-guy-shop/inventory-service/gen/gen/query"
	"github.com/people257/poor-guy-shop/inventory-service/internal/domain/inventory"
)

// ReservationRepository 预占记录仓储实现
type ReservationRepository struct {
	db    *gorm.DB
	query *query.Query
}

// NewReservationRepository 创建预占记录仓储
func NewReservationRepository(db *gorm.DB, query *query.Query) *ReservationRepository {
	return &ReservationRepository{
		db:    db,
		query: query,
	}
}

// Create 创建预占记录
func (r *ReservationRepository) Create(ctx context.Context, res *inventory.InventoryReservation) error {
	reservationModel := r.domainToModel(res)

	q := r.query.InventoryReservation
	if err := q.WithContext(ctx).Create(reservationModel); err != nil {
		return err
	}

	// 更新领域对象的ID
	if reservationModel.ID != "" {
		parsedID, err := uuid.Parse(reservationModel.ID)
		if err == nil {
			res.ID = parsedID
		}
	}

	return nil
}

// Update 更新预占记录
func (r *ReservationRepository) Update(ctx context.Context, res *inventory.InventoryReservation) error {
	reservationModel := r.domainToModel(res)

	q := r.query.InventoryReservation
	_, err := q.WithContext(ctx).Where(q.ID.Eq(reservationModel.ID)).Updates(reservationModel)
	return err
}

// UpdateWithVersion 乐观锁更新预占记录
func (r *ReservationRepository) UpdateWithVersion(ctx context.Context, res *inventory.InventoryReservation, version int32) error {
	reservationModel := r.domainToModel(res)

	q := r.query.InventoryReservation
	result, err := q.WithContext(ctx).
		Where(q.ID.Eq(reservationModel.ID), q.Version.Eq(version)).
		Updates(reservationModel)

	if err != nil {
		return err
	}

	if result.RowsAffected == 0 {
		return gorm.ErrRecordNotFound // 乐观锁冲突
	}

	return nil
}

// GetByID 根据ID获取预占记录
func (r *ReservationRepository) GetByID(ctx context.Context, id uuid.UUID) (*inventory.InventoryReservation, error) {
	q := r.query.InventoryReservation
	reservationModel, err := q.WithContext(ctx).Where(q.ID.Eq(id.String())).First()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, inventory.ErrReservationNotFound
		}
		return nil, err
	}

	return r.modelToDomain(reservationModel), nil
}

// GetByOrderID 根据订单ID获取预占记录列表
func (r *ReservationRepository) GetByOrderID(ctx context.Context, orderID uuid.UUID) ([]*inventory.InventoryReservation, error) {
	q := r.query.InventoryReservation
	reservationModels, err := q.WithContext(ctx).
		Where(q.OrderID.Eq(orderID.String())).
		Order(q.CreatedAt.Desc()).
		Find()
	if err != nil {
		return nil, err
	}

	reservations := make([]*inventory.InventoryReservation, len(reservationModels))
	for i, model := range reservationModels {
		reservations[i] = r.modelToDomain(model)
	}

	return reservations, nil
}

// GetBySkuID 根据SKU ID获取预占记录列表
func (r *ReservationRepository) GetBySkuID(ctx context.Context, skuID uuid.UUID, offset, limit int) ([]*inventory.InventoryReservation, int64, error) {
	q := r.query.InventoryReservation

	// 查询总数
	total, err := q.WithContext(ctx).Where(q.SkuID.Eq(skuID.String())).Count()
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	reservationModels, err := q.WithContext(ctx).
		Where(q.SkuID.Eq(skuID.String())).
		Offset(offset).
		Limit(limit).
		Order(q.CreatedAt.Desc()).
		Find()
	if err != nil {
		return nil, 0, err
	}

	reservations := make([]*inventory.InventoryReservation, len(reservationModels))
	for i, model := range reservationModels {
		reservations[i] = r.modelToDomain(model)
	}

	return reservations, total, nil
}

// GetExpiredReservations 获取过期的预占记录
func (r *ReservationRepository) GetExpiredReservations(ctx context.Context, limit int) ([]*inventory.InventoryReservation, error) {
	q := r.query.InventoryReservation
	now := time.Now()

	reservationModels, err := q.WithContext(ctx).
		Where(
			q.Status.Eq(string(inventory.ReservationStatusReserved)),
			q.ExpiresAt.Lt(now),
		).
		Limit(limit).
		Order(q.ExpiresAt.Asc()).
		Find()
	if err != nil {
		return nil, err
	}

	reservations := make([]*inventory.InventoryReservation, len(reservationModels))
	for i, model := range reservationModels {
		reservations[i] = r.modelToDomain(model)
	}

	return reservations, nil
}

// GetByStatus 根据状态获取预占记录
func (r *ReservationRepository) GetByStatus(ctx context.Context, status inventory.ReservationStatus, offset, limit int) ([]*inventory.InventoryReservation, int64, error) {
	q := r.query.InventoryReservation

	// 查询总数
	total, err := q.WithContext(ctx).Where(q.Status.Eq(string(status))).Count()
	if err != nil {
		return nil, 0, err
	}

	// 分页查询
	reservationModels, err := q.WithContext(ctx).
		Where(q.Status.Eq(string(status))).
		Offset(offset).
		Limit(limit).
		Order(q.CreatedAt.Desc()).
		Find()
	if err != nil {
		return nil, 0, err
	}

	reservations := make([]*inventory.InventoryReservation, len(reservationModels))
	for i, model := range reservationModels {
		reservations[i] = r.modelToDomain(model)
	}

	return reservations, total, nil
}

// Delete 删除预占记录
func (r *ReservationRepository) Delete(ctx context.Context, id uuid.UUID) error {
	q := r.query.InventoryReservation
	_, err := q.WithContext(ctx).Where(q.ID.Eq(id.String())).Delete()
	return err
}

// modelToDomain 将数据库模型转换为领域对象
func (r *ReservationRepository) modelToDomain(model *model.InventoryReservation) *inventory.InventoryReservation {
	id, _ := uuid.Parse(model.ID)
	skuID, _ := uuid.Parse(model.SkuID)
	orderID, _ := uuid.Parse(model.OrderID)

	var expiresAt *time.Time
	if model.ExpiresAt != nil {
		expiresAt = model.ExpiresAt
	}

	var confirmedAt *time.Time
	if model.ConfirmedAt != nil {
		confirmedAt = model.ConfirmedAt
	}

	var releasedAt *time.Time
	if model.ReleasedAt != nil {
		releasedAt = model.ReleasedAt
	}

	return &inventory.InventoryReservation{
		ID:          id,
		SkuID:       skuID,
		OrderID:     orderID,
		Quantity:    model.Quantity,
		Status:      inventory.ReservationStatus(model.Status),
		CreatedAt:   model.CreatedAt,
		UpdatedAt:   model.UpdatedAt,
		ExpiresAt:   expiresAt,
		ConfirmedAt: confirmedAt,
		ReleasedAt:  releasedAt,
		Version:     model.Version,
	}
}

// domainToModel 将领域对象转换为数据库模型
func (r *ReservationRepository) domainToModel(res *inventory.InventoryReservation) *model.InventoryReservation {
	return &model.InventoryReservation{
		ID:          res.ID.String(),
		SkuID:       res.SkuID.String(),
		OrderID:     res.OrderID.String(),
		Quantity:    res.Quantity,
		Status:      string(res.Status),
		CreatedAt:   res.CreatedAt,
		UpdatedAt:   res.UpdatedAt,
		ExpiresAt:   res.ExpiresAt,
		ConfirmedAt: res.ConfirmedAt,
		ReleasedAt:  res.ReleasedAt,
		Version:     res.Version,
	}
}
