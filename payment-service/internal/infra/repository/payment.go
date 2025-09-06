package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/people257/poor-guy-shop/payment-service/gen/gen/model"
	"github.com/people257/poor-guy-shop/payment-service/gen/gen/query"
	"github.com/people257/poor-guy-shop/payment-service/internal/domain/payment"
	"gorm.io/gorm"
)

// PaymentRepository 支付订单仓储实现
type PaymentRepository struct {
	db *gorm.DB
	q  *query.Query
}

// NewPaymentRepository 创建支付订单仓储
func NewPaymentRepository(db *gorm.DB, q *query.Query) payment.Repository {
	return &PaymentRepository{
		db: db,
		q:  q,
	}
}

// Create 创建支付订单
func (r *PaymentRepository) Create(ctx context.Context, paymentOrder *payment.PaymentOrder) error {
	// 转换为GORM模型
	modelPayment := r.entityToModel(paymentOrder)

	if err := r.q.PaymentOrder.WithContext(ctx).Create(modelPayment); err != nil {
		return fmt.Errorf("failed to create payment order: %w", err)
	}

	// 更新ID
	paymentOrder.ID = uuid.MustParse(modelPayment.ID)
	return nil
}

// GetByID 根据ID获取支付订单
func (r *PaymentRepository) GetByID(ctx context.Context, id uuid.UUID) (*payment.PaymentOrder, error) {
	modelPayment, err := r.q.PaymentOrder.WithContext(ctx).Where(r.q.PaymentOrder.ID.Eq(id.String())).First()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("payment order not found")
		}
		return nil, fmt.Errorf("failed to get payment order: %w", err)
	}

	return r.modelToEntity(modelPayment), nil
}

// GetByOrderID 根据业务订单ID获取支付订单
func (r *PaymentRepository) GetByOrderID(ctx context.Context, orderID uuid.UUID) (*payment.PaymentOrder, error) {
	modelPayment, err := r.q.PaymentOrder.WithContext(ctx).Where(r.q.PaymentOrder.OrderID.Eq(orderID.String())).First()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("payment order not found")
		}
		return nil, fmt.Errorf("failed to get payment order: %w", err)
	}

	return r.modelToEntity(modelPayment), nil
}

// GetByThirdPartyOrderID 根据第三方订单ID获取支付订单
func (r *PaymentRepository) GetByThirdPartyOrderID(ctx context.Context, thirdPartyOrderID string) (*payment.PaymentOrder, error) {
	modelPayment, err := r.q.PaymentOrder.WithContext(ctx).Where(r.q.PaymentOrder.ThirdPartyOrderID.Eq(thirdPartyOrderID)).First()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("payment order not found")
		}
		return nil, fmt.Errorf("failed to get payment order: %w", err)
	}

	return r.modelToEntity(modelPayment), nil
}

// Update 更新支付订单
func (r *PaymentRepository) Update(ctx context.Context, paymentOrder *payment.PaymentOrder) error {
	updates := map[string]interface{}{
		"status":               string(paymentOrder.Status),
		"third_party_order_id": paymentOrder.ThirdPartyOrderID,
		"third_party_response": paymentOrder.ThirdPartyResponse,
		"updated_at":           paymentOrder.UpdatedAt,
		"paid_at":              paymentOrder.PaidAt,
	}

	_, err := r.q.PaymentOrder.WithContext(ctx).Where(r.q.PaymentOrder.ID.Eq(paymentOrder.ID.String())).Updates(updates)
	if err != nil {
		return fmt.Errorf("failed to update payment order: %w", err)
	}

	return nil
}

// List 分页查询支付订单
func (r *PaymentRepository) List(ctx context.Context, userID uuid.UUID, page, pageSize int) ([]*payment.PaymentOrder, int64, error) {
	offset := (page - 1) * pageSize

	// 查询总数
	count, err := r.q.PaymentOrder.WithContext(ctx).Where(r.q.PaymentOrder.UserID.Eq(userID.String())).Count()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count payment orders: %w", err)
	}

	// 查询列表
	models, err := r.q.PaymentOrder.WithContext(ctx).
		Where(r.q.PaymentOrder.UserID.Eq(userID.String())).
		Order(r.q.PaymentOrder.CreatedAt.Desc()).
		Limit(pageSize).
		Offset(offset).
		Find()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list payment orders: %w", err)
	}

	entities := make([]*payment.PaymentOrder, len(models))
	for i, modelPayment := range models {
		entities[i] = r.modelToEntity(modelPayment)
	}

	return entities, count, nil
}

// CreateLog 创建支付日志
func (r *PaymentRepository) CreateLog(ctx context.Context, log *payment.PaymentLog) error {
	modelLog := r.logEntityToModel(log)

	if err := r.q.PaymentLog.WithContext(ctx).Create(modelLog); err != nil {
		return fmt.Errorf("failed to create payment log: %w", err)
	}

	log.ID = uuid.MustParse(modelLog.ID)
	return nil
}

// GetLogs 获取支付日志
func (r *PaymentRepository) GetLogs(ctx context.Context, paymentOrderID uuid.UUID) ([]*payment.PaymentLog, error) {
	models, err := r.q.PaymentLog.WithContext(ctx).
		Where(r.q.PaymentLog.PaymentOrderID.Eq(paymentOrderID.String())).
		Order(r.q.PaymentLog.CreatedAt.Desc()).
		Find()
	if err != nil {
		return nil, fmt.Errorf("failed to get payment logs: %w", err)
	}

	logs := make([]*payment.PaymentLog, len(models))
	for i, modelLog := range models {
		logs[i] = r.logModelToEntity(modelLog)
	}

	return logs, nil
}

// entityToModel 将领域实体转换为GORM模型
func (r *PaymentRepository) entityToModel(entity *payment.PaymentOrder) *model.PaymentOrder {
	modelPayment := &model.PaymentOrder{
		ID:            entity.ID.String(),
		OrderID:       entity.OrderID.String(),
		UserID:        entity.UserID.String(),
		Amount:        entity.Amount,
		PaymentMethod: string(entity.PaymentMethod),
		Status:        string(entity.Status),
		CreatedAt:     entity.CreatedAt,
		UpdatedAt:     entity.UpdatedAt,
		PaidAt:        entity.PaidAt,
		ExpiredAt:     entity.ExpiredAt,
	}

	if entity.ThirdPartyOrderID != "" {
		modelPayment.ThirdPartyOrderID = &entity.ThirdPartyOrderID
	}
	if entity.ThirdPartyResponse != "" {
		modelPayment.ThirdPartyResponse = &entity.ThirdPartyResponse
	}
	if entity.Subject != "" {
		modelPayment.Subject = &entity.Subject
	}
	if entity.Description != "" {
		modelPayment.Description = &entity.Description
	}
	if entity.NotifyURL != "" {
		modelPayment.NotifyURL = &entity.NotifyURL
	}
	if entity.ReturnURL != "" {
		modelPayment.ReturnURL = &entity.ReturnURL
	}

	return modelPayment
}

// modelToEntity 将GORM模型转换为领域实体
func (r *PaymentRepository) modelToEntity(modelPayment *model.PaymentOrder) *payment.PaymentOrder {
	entity := &payment.PaymentOrder{
		ID:            uuid.MustParse(modelPayment.ID),
		OrderID:       uuid.MustParse(modelPayment.OrderID),
		UserID:        uuid.MustParse(modelPayment.UserID),
		Amount:        modelPayment.Amount,
		PaymentMethod: payment.PaymentMethod(modelPayment.PaymentMethod),
		Status:        payment.PaymentStatus(modelPayment.Status),
		CreatedAt:     modelPayment.CreatedAt,
		UpdatedAt:     modelPayment.UpdatedAt,
		PaidAt:        modelPayment.PaidAt,
		ExpiredAt:     modelPayment.ExpiredAt,
	}

	if modelPayment.ThirdPartyOrderID != nil {
		entity.ThirdPartyOrderID = *modelPayment.ThirdPartyOrderID
	}
	if modelPayment.ThirdPartyResponse != nil {
		entity.ThirdPartyResponse = *modelPayment.ThirdPartyResponse
	}
	if modelPayment.Subject != nil {
		entity.Subject = *modelPayment.Subject
	}
	if modelPayment.Description != nil {
		entity.Description = *modelPayment.Description
	}
	if modelPayment.NotifyURL != nil {
		entity.NotifyURL = *modelPayment.NotifyURL
	}
	if modelPayment.ReturnURL != nil {
		entity.ReturnURL = *modelPayment.ReturnURL
	}

	return entity
}

// logEntityToModel 将支付日志实体转换为GORM模型
func (r *PaymentRepository) logEntityToModel(entity *payment.PaymentLog) *model.PaymentLog {
	modelLog := &model.PaymentLog{
		ID:             entity.ID.String(),
		PaymentOrderID: entity.PaymentOrderID.String(),
		Action:         entity.Action,
		CreatedAt:      entity.CreatedAt,
	}

	if entity.RequestData != "" {
		modelLog.RequestData = &entity.RequestData
	}
	if entity.ResponseData != "" {
		modelLog.ResponseData = &entity.ResponseData
	}
	if entity.ErrorMessage != "" {
		modelLog.ErrorMessage = &entity.ErrorMessage
	}

	return modelLog
}

// logModelToEntity 将支付日志GORM模型转换为领域实体
func (r *PaymentRepository) logModelToEntity(modelLog *model.PaymentLog) *payment.PaymentLog {
	entity := &payment.PaymentLog{
		ID:             uuid.MustParse(modelLog.ID),
		PaymentOrderID: uuid.MustParse(modelLog.PaymentOrderID),
		Action:         modelLog.Action,
		CreatedAt:      modelLog.CreatedAt,
	}

	if modelLog.RequestData != nil {
		entity.RequestData = *modelLog.RequestData
	}
	if modelLog.ResponseData != nil {
		entity.ResponseData = *modelLog.ResponseData
	}
	if modelLog.ErrorMessage != nil {
		entity.ErrorMessage = *modelLog.ErrorMessage
	}

	return entity
}
