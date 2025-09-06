package reservation

import (
	"context"
	"time"

	"github.com/google/uuid"

	"github.com/people257/poor-guy-shop/inventory-service/internal/domain/inventory"
	"github.com/people257/poor-guy-shop/inventory-service/internal/domain/reservation"
)

// Service 预占应用服务
type Service struct {
	reservationDomain *reservation.DomainService
	reservationRepo   reservation.Repository
}

// NewService 创建预占应用服务
func NewService(reservationDomain *reservation.DomainService, reservationRepo reservation.Repository) *Service {
	return &Service{
		reservationDomain: reservationDomain,
		reservationRepo:   reservationRepo,
	}
}

// CreateReservation 创建预占记录
func (s *Service) CreateReservation(ctx context.Context, skuID, orderID uuid.UUID, quantity int32, expiresAt *time.Time) (*inventory.InventoryReservation, error) {
	return s.reservationDomain.CreateReservation(ctx, skuID, orderID, quantity, expiresAt)
}

// GetReservation 获取预占记录
func (s *Service) GetReservation(ctx context.Context, id uuid.UUID) (*inventory.InventoryReservation, error) {
	return s.reservationRepo.GetByID(ctx, id)
}

// GetReservationsByOrderID 根据订单ID获取预占记录
func (s *Service) GetReservationsByOrderID(ctx context.Context, orderID uuid.UUID) ([]*inventory.InventoryReservation, error) {
	return s.reservationDomain.GetReservationsByOrderID(ctx, orderID)
}

// GetReservationsBySkuID 根据SKU ID获取预占记录
func (s *Service) GetReservationsBySkuID(ctx context.Context, skuID uuid.UUID, page, pageSize int) ([]*inventory.InventoryReservation, int64, error) {
	offset := (page - 1) * pageSize
	return s.reservationDomain.GetReservationsBySkuID(ctx, skuID, offset, pageSize)
}

// ReleaseReservationsByOrderID 释放订单的所有预占
func (s *Service) ReleaseReservationsByOrderID(ctx context.Context, orderID uuid.UUID) error {
	return s.reservationDomain.ReleaseReservationsByOrderID(ctx, orderID)
}

// ConfirmReservationsByOrderID 确认订单的所有预占
func (s *Service) ConfirmReservationsByOrderID(ctx context.Context, orderID uuid.UUID) error {
	return s.reservationDomain.ConfirmReservationsByOrderID(ctx, orderID)
}

// CleanupExpiredReservations 清理过期的预占记录
func (s *Service) CleanupExpiredReservations(ctx context.Context, limit int) (int, error) {
	return s.reservationDomain.CleanupExpiredReservations(ctx, limit)
}

// GetReservationsByStatus 根据状态获取预占记录
func (s *Service) GetReservationsByStatus(ctx context.Context, status inventory.ReservationStatus, page, pageSize int) ([]*inventory.InventoryReservation, int64, error) {
	offset := (page - 1) * pageSize
	return s.reservationDomain.GetReservationsByStatus(ctx, status, offset, pageSize)
}

// GetExpiredReservations 获取过期的预占记录
func (s *Service) GetExpiredReservations(ctx context.Context, limit int) ([]*inventory.InventoryReservation, error) {
	return s.reservationRepo.GetExpiredReservations(ctx, limit)
}

// UpdateReservation 更新预占记录
func (s *Service) UpdateReservation(ctx context.Context, reservation *inventory.InventoryReservation) error {
	return s.reservationRepo.Update(ctx, reservation)
}

// DeleteReservation 删除预占记录
func (s *Service) DeleteReservation(ctx context.Context, id uuid.UUID) error {
	return s.reservationRepo.Delete(ctx, id)
}
