package inventory

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/people257/poor-guy-shop/inventory-service/internal/application/reservation"
	"github.com/people257/poor-guy-shop/inventory-service/internal/domain/inventory"
	"github.com/people257/poor-guy-shop/inventory-service/internal/infra/client"
	"github.com/people257/poor-guy-shop/inventory-service/internal/infra/client/product"
)

// BusinessService 库存业务服务（集成内部RPC调用）
type BusinessService struct {
	inventoryService *Service
	reservationApp   *reservation.Service
	clientManager    *client.Manager
}

// NewBusinessService 创建库存业务服务
func NewBusinessService(inventoryService *Service, reservationApp *reservation.Service, clientManager *client.Manager) *BusinessService {
	return &BusinessService{
		inventoryService: inventoryService,
		reservationApp:   reservationApp,
		clientManager:    clientManager,
	}
}

// ReserveInventoryWithValidation 预占库存（带商品验证）
func (s *BusinessService) ReserveInventoryWithValidation(ctx context.Context, orderID uuid.UUID, items []inventory.ReserveItem, expiresAt *time.Time) ([]*inventory.InventoryReservation, error) {
	// 1. 验证商品是否存在且可售
	skuIDs := make([]uuid.UUID, len(items))
	for i, item := range items {
		skuIDs[i] = item.SkuID
	}

	validProducts, err := s.clientManager.ProductClient.ValidateProducts(ctx, skuIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to validate products: %w", err)
	}

	// 检查是否有无效商品
	var invalidSkus []uuid.UUID
	for _, item := range items {
		if !validProducts[item.SkuID] {
			invalidSkus = append(invalidSkus, item.SkuID)
		}
	}

	if len(invalidSkus) > 0 {
		return nil, fmt.Errorf("invalid products found: %v", invalidSkus)
	}

	// 2. 检查库存可用性
	available, insufficientSkus, err := s.inventoryService.CheckInventoryAvailability(ctx, items)
	if err != nil {
		return nil, fmt.Errorf("failed to check inventory availability: %w", err)
	}

	if !available {
		// 通知订单服务库存不足
		if err := s.clientManager.OrderClient.NotifyInventoryReserved(ctx, orderID.String(), false,
			fmt.Sprintf("insufficient inventory for skus: %v", insufficientSkus)); err != nil {
			// 记录错误但不阻塞主流程
		}
		return nil, fmt.Errorf("insufficient inventory for skus: %v", insufficientSkus)
	}

	// 3. 执行库存预占
	reservations, err := s.inventoryService.ReserveInventory(ctx, orderID, items, expiresAt)
	if err != nil {
		// 通知订单服务预占失败
		if notifyErr := s.clientManager.OrderClient.NotifyInventoryReserved(ctx, orderID.String(), false, err.Error()); notifyErr != nil {
			// 记录错误但不阻塞主流程
		}
		return nil, fmt.Errorf("failed to reserve inventory: %w", err)
	}

	// 4. 通知订单服务预占成功
	if err := s.clientManager.OrderClient.NotifyInventoryReserved(ctx, orderID.String(), true, "inventory reserved successfully"); err != nil {
		// 记录错误但不回滚，因为库存已经预占成功
	}

	return reservations, nil
}

// ConfirmInventoryWithOrderValidation 确认扣减库存（带订单验证）
func (s *BusinessService) ConfirmInventoryWithOrderValidation(ctx context.Context, orderID uuid.UUID) error {
	// 1. 验证订单状态
	orderStatus, err := s.clientManager.OrderClient.GetOrderStatus(ctx, orderID.String())
	if err != nil {
		return fmt.Errorf("failed to get order status: %w", err)
	}

	// 检查订单状态是否允许确认扣减
	if orderStatus != "paid" && orderStatus != "confirmed" {
		return fmt.Errorf("order status %s does not allow inventory confirmation", orderStatus)
	}

	// 2. 获取预占记录
	reservations, err := s.reservationApp.GetReservationsByOrderID(ctx, orderID)
	if err != nil {
		return fmt.Errorf("failed to get reservations: %w", err)
	}

	if len(reservations) == 0 {
		return fmt.Errorf("no reservations found for order %s", orderID.String())
	}

	// 3. 确认扣减库存
	for _, reservation := range reservations {
		if reservation.Status == inventory.ReservationStatusReserved {
			if err := s.inventoryService.inventoryDomain.ConfirmReservation(ctx, reservation); err != nil {
				// 通知订单服务确认失败
				if notifyErr := s.clientManager.OrderClient.NotifyInventoryConfirmed(ctx, orderID.String(), false, err.Error()); notifyErr != nil {
					// 记录错误
				}
				return fmt.Errorf("failed to confirm reservation %s: %w", reservation.ID.String(), err)
			}
		}
	}

	// 4. 通知订单服务确认成功
	if err := s.clientManager.OrderClient.NotifyInventoryConfirmed(ctx, orderID.String(), true, "inventory confirmed successfully"); err != nil {
		// 记录错误但不回滚
	}

	return nil
}

// ReleaseInventoryWithOrderValidation 释放库存（带订单验证）
func (s *BusinessService) ReleaseInventoryWithOrderValidation(ctx context.Context, orderID uuid.UUID) error {
	// 1. 验证订单状态
	orderStatus, err := s.clientManager.OrderClient.GetOrderStatus(ctx, orderID.String())
	if err != nil {
		return fmt.Errorf("failed to get order status: %w", err)
	}

	// 检查订单状态是否允许释放库存
	if orderStatus == "confirmed" || orderStatus == "shipped" {
		return fmt.Errorf("order status %s does not allow inventory release", orderStatus)
	}

	// 2. 释放库存
	if err := s.reservationApp.ReleaseReservationsByOrderID(ctx, orderID); err != nil {
		return fmt.Errorf("failed to release inventory: %w", err)
	}

	return nil
}

// GetInventoryWithProductInfo 获取库存信息（包含商品信息）
func (s *BusinessService) GetInventoryWithProductInfo(ctx context.Context, skuID uuid.UUID) (*InventoryWithProductInfo, error) {
	// 1. 获取库存信息
	inv, err := s.inventoryService.GetInventory(ctx, skuID)
	if err != nil {
		return nil, fmt.Errorf("failed to get inventory: %w", err)
	}

	// 2. 获取商品信息
	productInfo, err := s.clientManager.ProductClient.GetProduct(ctx, skuID)
	if err != nil {
		return nil, fmt.Errorf("failed to get product info: %w", err)
	}

	return &InventoryWithProductInfo{
		Inventory:   inv,
		ProductInfo: productInfo,
	}, nil
}

// BatchGetInventoryWithProductInfo 批量获取库存信息（包含商品信息）
func (s *BusinessService) BatchGetInventoryWithProductInfo(ctx context.Context, skuIDs []uuid.UUID) ([]*InventoryWithProductInfo, error) {
	// 1. 批量获取库存信息
	inventories, err := s.inventoryService.BatchGetInventory(ctx, skuIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to batch get inventory: %w", err)
	}

	// 2. 批量获取商品信息
	products, err := s.clientManager.ProductClient.BatchGetProducts(ctx, skuIDs)
	if err != nil {
		return nil, fmt.Errorf("failed to batch get products: %w", err)
	}

	// 3. 构建结果映射
	productMap := make(map[uuid.UUID]*product.ProductInfo)
	for _, prod := range products {
		productMap[prod.ID] = prod
	}

	// 4. 组合结果
	results := make([]*InventoryWithProductInfo, 0, len(inventories))
	for _, inv := range inventories {
		if productInfo, exists := productMap[inv.SkuID]; exists {
			results = append(results, &InventoryWithProductInfo{
				Inventory:   inv,
				ProductInfo: productInfo,
			})
		}
	}

	return results, nil
}

// InventoryWithProductInfo 库存和商品信息组合
type InventoryWithProductInfo struct {
	Inventory   *inventory.Inventory
	ProductInfo *product.ProductInfo
}

// HealthCheck 健康检查
func (s *BusinessService) HealthCheck(ctx context.Context) map[string]error {
	results := make(map[string]error)

	// 检查内部服务
	if s.inventoryService == nil {
		results["inventory_service"] = fmt.Errorf("inventory service is nil")
	}

	// 检查客户端连接
	if s.clientManager != nil {
		clientResults := s.clientManager.HealthCheck(ctx)
		for service, err := range clientResults {
			results["client_"+service] = err
		}
	}

	return results
}
