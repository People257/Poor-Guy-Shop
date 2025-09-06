package inventory

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/people257/poor-guy-shop/inventory-service/internal/application/reservation"
)

// Scheduler 定时任务调度器
type Scheduler struct {
	businessService *BusinessService
	eventHandler    *EventHandler
	reservationApp  *reservation.Service

	stopCh chan struct{}
}

// NewScheduler 创建定时任务调度器
func NewScheduler(businessService *BusinessService, eventHandler *EventHandler, reservationApp *reservation.Service) *Scheduler {
	return &Scheduler{
		businessService: businessService,
		eventHandler:    eventHandler,
		reservationApp:  reservationApp,
		stopCh:          make(chan struct{}),
	}
}

// Start 启动定时任务
func (s *Scheduler) Start(ctx context.Context) {
	// 清理过期预占记录任务 - 每5分钟执行一次
	go s.runCleanupExpiredReservations(ctx)

	// 库存告警检查任务 - 每10分钟执行一次
	go s.runInventoryAlertCheck(ctx)

	// 健康检查任务 - 每1分钟执行一次
	go s.runHealthCheck(ctx)
}

// Stop 停止定时任务
func (s *Scheduler) Stop() {
	close(s.stopCh)
}

// runCleanupExpiredReservations 运行清理过期预占记录任务
func (s *Scheduler) runCleanupExpiredReservations(ctx context.Context) {
	ticker := time.NewTicker(5 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-s.stopCh:
			return
		case <-ticker.C:
			if err := s.cleanupExpiredReservations(ctx); err != nil {
				log.Printf("Failed to cleanup expired reservations: %v", err)
			}
		}
	}
}

// runInventoryAlertCheck 运行库存告警检查任务
func (s *Scheduler) runInventoryAlertCheck(ctx context.Context) {
	ticker := time.NewTicker(10 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-s.stopCh:
			return
		case <-ticker.C:
			if err := s.checkInventoryAlerts(ctx); err != nil {
				log.Printf("Failed to check inventory alerts: %v", err)
			}
		}
	}
}

// runHealthCheck 运行健康检查任务
func (s *Scheduler) runHealthCheck(ctx context.Context) {
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-s.stopCh:
			return
		case <-ticker.C:
			s.performHealthCheck(ctx)
		}
	}
}

// cleanupExpiredReservations 清理过期预占记录
func (s *Scheduler) cleanupExpiredReservations(ctx context.Context) error {
	const batchSize = 100

	cleanedCount, err := s.reservationApp.CleanupExpiredReservations(ctx, batchSize)
	if err != nil {
		return fmt.Errorf("failed to cleanup expired reservations: %w", err)
	}

	if cleanedCount > 0 {
		log.Printf("Cleaned up %d expired reservations", cleanedCount)
	}

	return nil
}

// checkInventoryAlerts 检查库存告警
func (s *Scheduler) checkInventoryAlerts(ctx context.Context) error {
	const pageSize = 50

	// 检查库存不足的商品
	lowStockInventories, _, err := s.businessService.inventoryService.ListLowStockInventory(ctx, 1, pageSize)
	if err != nil {
		return fmt.Errorf("failed to get low stock inventories: %w", err)
	}

	for _, inv := range lowStockInventories {
		if err := s.eventHandler.HandleInventoryAlert(ctx, inv.SkuID, "low_stock", inv.AvailableQuantity); err != nil {
			log.Printf("Failed to handle low stock alert for SKU %s: %v", inv.SkuID.String(), err)
		}
	}

	// 检查售罄的商品
	outOfStockInventories, _, err := s.businessService.inventoryService.ListOutOfStockInventory(ctx, 1, pageSize)
	if err != nil {
		return fmt.Errorf("failed to get out of stock inventories: %w", err)
	}

	for _, inv := range outOfStockInventories {
		if err := s.eventHandler.HandleInventoryAlert(ctx, inv.SkuID, "out_of_stock", inv.AvailableQuantity); err != nil {
			log.Printf("Failed to handle out of stock alert for SKU %s: %v", inv.SkuID.String(), err)
		}
	}

	return nil
}

// performHealthCheck 执行健康检查
func (s *Scheduler) performHealthCheck(ctx context.Context) {
	if s.businessService == nil {
		log.Printf("Health check failed: business service is nil")
		return
	}

	healthResults := s.businessService.HealthCheck(ctx)
	unhealthyServices := make([]string, 0)

	for service, err := range healthResults {
		if err != nil {
			unhealthyServices = append(unhealthyServices, service)
		}
	}

	if len(unhealthyServices) > 0 {
		log.Printf("Health check found unhealthy services: %v", unhealthyServices)
	}
}
