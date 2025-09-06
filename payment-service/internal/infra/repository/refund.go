package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/people257/poor-guy-shop/payment-service/gen/gen/model"
	"github.com/people257/poor-guy-shop/payment-service/gen/gen/query"
	"github.com/people257/poor-guy-shop/payment-service/internal/domain/refund"
	"gorm.io/gorm"
)

// RefundRepository 退款仓储实现
type RefundRepository struct {
	db *gorm.DB
	q  *query.Query
}

// NewRefundRepository 创建退款仓储
func NewRefundRepository(db *gorm.DB, q *query.Query) refund.Repository {
	return &RefundRepository{
		db: db,
		q:  q,
	}
}

// Create 创建退款记录
func (r *RefundRepository) Create(ctx context.Context, refundEntity *refund.Refund) error {
	modelRefund := r.entityToModel(refundEntity)

	if err := r.q.Refund.WithContext(ctx).Create(modelRefund); err != nil {
		return fmt.Errorf("failed to create refund: %w", err)
	}

	refundEntity.ID = uuid.MustParse(modelRefund.ID)
	return nil
}

// GetByID 根据ID获取退款记录
func (r *RefundRepository) GetByID(ctx context.Context, id uuid.UUID) (*refund.Refund, error) {
	modelRefund, err := r.q.Refund.WithContext(ctx).Where(r.q.Refund.ID.Eq(id.String())).First()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("refund not found")
		}
		return nil, fmt.Errorf("failed to get refund: %w", err)
	}

	return r.modelToEntity(modelRefund), nil
}

// GetByPaymentOrderID 根据支付订单ID获取退款记录
func (r *RefundRepository) GetByPaymentOrderID(ctx context.Context, paymentOrderID uuid.UUID) ([]*refund.Refund, error) {
	models, err := r.q.Refund.WithContext(ctx).
		Where(r.q.Refund.PaymentOrderID.Eq(paymentOrderID.String())).
		Order(r.q.Refund.CreatedAt.Desc()).
		Find()
	if err != nil {
		return nil, fmt.Errorf("failed to get refunds: %w", err)
	}

	entities := make([]*refund.Refund, len(models))
	for i, modelRefund := range models {
		entities[i] = r.modelToEntity(modelRefund)
	}

	return entities, nil
}

// Update 更新退款记录
func (r *RefundRepository) Update(ctx context.Context, refundEntity *refund.Refund) error {
	updates := map[string]interface{}{
		"status":                string(refundEntity.Status),
		"third_party_refund_id": refundEntity.ThirdPartyRefundID,
		"third_party_response":  refundEntity.ThirdPartyResponse,
		"updated_at":            refundEntity.UpdatedAt,
		"processed_at":          refundEntity.ProcessedAt,
	}

	_, err := r.q.Refund.WithContext(ctx).Where(r.q.Refund.ID.Eq(refundEntity.ID.String())).Updates(updates)
	if err != nil {
		return fmt.Errorf("failed to update refund: %w", err)
	}

	return nil
}

// List 分页查询退款记录
func (r *RefundRepository) List(ctx context.Context, page, pageSize int) ([]*refund.Refund, int64, error) {
	offset := (page - 1) * pageSize

	// 查询总数
	count, err := r.q.Refund.WithContext(ctx).Count()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to count refunds: %w", err)
	}

	// 查询列表
	models, err := r.q.Refund.WithContext(ctx).
		Order(r.q.Refund.CreatedAt.Desc()).
		Limit(pageSize).
		Offset(offset).
		Find()
	if err != nil {
		return nil, 0, fmt.Errorf("failed to list refunds: %w", err)
	}

	entities := make([]*refund.Refund, len(models))
	for i, modelRefund := range models {
		entities[i] = r.modelToEntity(modelRefund)
	}

	return entities, count, nil
}

// entityToModel 将领域实体转换为GORM模型
func (r *RefundRepository) entityToModel(entity *refund.Refund) *model.Refund {
	modelRefund := &model.Refund{
		ID:             entity.ID.String(),
		PaymentOrderID: entity.PaymentOrderID.String(),
		Amount:         entity.Amount,
		Status:         string(entity.Status),
		CreatedAt:      entity.CreatedAt,
		UpdatedAt:      entity.UpdatedAt,
		ProcessedAt:    entity.ProcessedAt,
	}

	if entity.Reason != "" {
		modelRefund.Reason = &entity.Reason
	}
	if entity.ThirdPartyRefundID != "" {
		modelRefund.ThirdPartyRefundID = &entity.ThirdPartyRefundID
	}
	if entity.ThirdPartyResponse != "" {
		modelRefund.ThirdPartyResponse = &entity.ThirdPartyResponse
	}

	return modelRefund
}

// modelToEntity 将GORM模型转换为领域实体
func (r *RefundRepository) modelToEntity(modelRefund *model.Refund) *refund.Refund {
	entity := &refund.Refund{
		ID:             uuid.MustParse(modelRefund.ID),
		PaymentOrderID: uuid.MustParse(modelRefund.PaymentOrderID),
		Amount:         modelRefund.Amount,
		Status:         refund.RefundStatus(modelRefund.Status),
		CreatedAt:      modelRefund.CreatedAt,
		UpdatedAt:      modelRefund.UpdatedAt,
		ProcessedAt:    modelRefund.ProcessedAt,
	}

	if modelRefund.Reason != nil {
		entity.Reason = *modelRefund.Reason
	}
	if modelRefund.ThirdPartyRefundID != nil {
		entity.ThirdPartyRefundID = *modelRefund.ThirdPartyRefundID
	}
	if modelRefund.ThirdPartyResponse != nil {
		entity.ThirdPartyResponse = *modelRefund.ThirdPartyResponse
	}

	return entity
}
