package reservation

import (
	"context"
	"time"

	"github.com/google/uuid"
	"github.com/people257/poor-guy-shop/inventory-service/internal/domain/inventory"
)

// DomainService 预占领域服务
type DomainService struct {
	reservationRepo Repository
	inventoryDomain *inventory.DomainService
}

// NewDomainService 创建预占领域服务
func NewDomainService(reservationRepo Repository, inventoryDomain *inventory.DomainService) *DomainService {
	return &DomainService{
		reservationRepo: reservationRepo,
		inventoryDomain: inventoryDomain,
	}
}

// CreateReservation 创建预占记录
func (s *DomainService) CreateReservation(ctx context.Context, skuID, orderID uuid.UUID, quantity int32, expiresAt *time.Time) (*inventory.InventoryReservation, error) {
	if quantity <= 0 {
		return nil, inventory.ErrInvalidQuantity
	}

	// 创建预占记录
	reservation := inventory.NewInventoryReservation(skuID, orderID, quantity, expiresAt)

	// 保存预占记录
	if err := s.reservationRepo.Create(ctx, reservation); err != nil {
		return nil, err
	}

	return reservation, nil
}

// GetReservationsByOrderID 根据订单ID获取预占记录
func (s *DomainService) GetReservationsByOrderID(ctx context.Context, orderID uuid.UUID) ([]*inventory.InventoryReservation, error) {
	return s.reservationRepo.GetByOrderID(ctx, orderID)
}

// ReleaseReservationsByOrderID 释放订单的所有预占
func (s *DomainService) ReleaseReservationsByOrderID(ctx context.Context, orderID uuid.UUID) error {
	reservations, err := s.reservationRepo.GetByOrderID(ctx, orderID)
	if err != nil {
		return err
	}

	for _, reservation := range reservations {
		if reservation.Status == inventory.ReservationStatusReserved {
			// 释放库存
			if err := s.inventoryDomain.ReleaseReservation(ctx, reservation); err != nil {
				// TODO: 记录错误日志，但继续处理其他预占
				continue
			}

			// 更新预占记录状态
			if err := s.reservationRepo.UpdateWithVersion(ctx, reservation, reservation.Version); err != nil {
				// TODO: 记录错误日志
			}
		}
	}

	return nil
}

// ConfirmReservationsByOrderID 确认订单的所有预占
func (s *DomainService) ConfirmReservationsByOrderID(ctx context.Context, orderID uuid.UUID) error {
	reservations, err := s.reservationRepo.GetByOrderID(ctx, orderID)
	if err != nil {
		return err
	}

	for _, reservation := range reservations {
		if reservation.Status == inventory.ReservationStatusReserved {
			// 确认预占
			if err := s.inventoryDomain.ConfirmReservation(ctx, reservation); err != nil {
				// TODO: 记录错误日志，但继续处理其他预占
				continue
			}

			// 更新预占记录状态
			if err := s.reservationRepo.UpdateWithVersion(ctx, reservation, reservation.Version); err != nil {
				// TODO: 记录错误日志
			}
		}
	}

	return nil
}

// CleanupExpiredReservations 清理过期的预占记录
func (s *DomainService) CleanupExpiredReservations(ctx context.Context, limit int) (int, error) {
	expiredReservations, err := s.reservationRepo.GetExpiredReservations(ctx, limit)
	if err != nil {
		return 0, err
	}

	cleanedCount := 0
	for _, reservation := range expiredReservations {
		if reservation.Status == inventory.ReservationStatusReserved {
			// 释放库存
			if err := s.inventoryDomain.ReleaseReservation(ctx, reservation); err != nil {
				// TODO: 记录错误日志
				continue
			}

			// 标记为过期
			reservation.MarkExpired()

			// 更新预占记录状态
			if err := s.reservationRepo.UpdateWithVersion(ctx, reservation, reservation.Version-1); err != nil {
				// TODO: 记录错误日志
				continue
			}

			cleanedCount++
		}
	}

	return cleanedCount, nil
}

// GetReservationsBySkuID 根据SKU ID获取预占记录
func (s *DomainService) GetReservationsBySkuID(ctx context.Context, skuID uuid.UUID, offset, limit int) ([]*inventory.InventoryReservation, int64, error) {
	return s.reservationRepo.GetBySkuID(ctx, skuID, offset, limit)
}

// GetReservationsByStatus 根据状态获取预占记录
func (s *DomainService) GetReservationsByStatus(ctx context.Context, status inventory.ReservationStatus, offset, limit int) ([]*inventory.InventoryReservation, int64, error) {
	return s.reservationRepo.GetByStatus(ctx, status, offset, limit)
}