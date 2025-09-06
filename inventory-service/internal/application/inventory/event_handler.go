package inventory

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/people257/poor-guy-shop/inventory-service/internal/domain/inventory"
)

// OrderEvent 订单事件
type OrderEvent struct {
	Type      string      `json:"type"`
	OrderID   uuid.UUID   `json:"order_id"`
	UserID    uuid.UUID   `json:"user_id"`
	Items     []OrderItem `json:"items"`
	Status    string      `json:"status"`
	Timestamp time.Time   `json:"timestamp"`
}

// OrderItem 订单商品项
type OrderItem struct {
	SkuID    uuid.UUID `json:"sku_id"`
	Quantity int32     `json:"quantity"`
	Price    float64   `json:"price"`
}

// EventHandler 库存事件处理器
type EventHandler struct {
	businessService *BusinessService
}

// NewEventHandler 创建事件处理器
func NewEventHandler(businessService *BusinessService) *EventHandler {
	return &EventHandler{
		businessService: businessService,
	}
}

// HandleOrderCreated 处理订单创建事件
func (h *EventHandler) HandleOrderCreated(ctx context.Context, eventData []byte) error {
	var event OrderEvent
	if err := json.Unmarshal(eventData, &event); err != nil {
		return fmt.Errorf("failed to unmarshal order created event: %w", err)
	}

	// 转换为库存预占项
	items := make([]inventory.ReserveItem, len(event.Items))
	for i, item := range event.Items {
		items[i] = inventory.ReserveItem{
			SkuID:    item.SkuID,
			Quantity: item.Quantity,
		}
	}

	// 设置30分钟过期时间
	expiresAt := time.Now().Add(30 * time.Minute)

	// 预占库存
	_, err := h.businessService.ReserveInventoryWithValidation(ctx, event.OrderID, items, &expiresAt)
	if err != nil {
		return fmt.Errorf("failed to reserve inventory for order %s: %w", event.OrderID.String(), err)
	}

	return nil
}

// HandleOrderPaid 处理订单支付成功事件
func (h *EventHandler) HandleOrderPaid(ctx context.Context, eventData []byte) error {
	var event OrderEvent
	if err := json.Unmarshal(eventData, &event); err != nil {
		return fmt.Errorf("failed to unmarshal order paid event: %w", err)
	}

	// 确认扣减库存
	if err := h.businessService.ConfirmInventoryWithOrderValidation(ctx, event.OrderID); err != nil {
		return fmt.Errorf("failed to confirm inventory for order %s: %w", event.OrderID.String(), err)
	}

	return nil
}

// HandleOrderCancelled 处理订单取消事件
func (h *EventHandler) HandleOrderCancelled(ctx context.Context, eventData []byte) error {
	var event OrderEvent
	if err := json.Unmarshal(eventData, &event); err != nil {
		return fmt.Errorf("failed to unmarshal order cancelled event: %w", err)
	}

	// 释放库存
	if err := h.businessService.ReleaseInventoryWithOrderValidation(ctx, event.OrderID); err != nil {
		return fmt.Errorf("failed to release inventory for order %s: %w", event.OrderID.String(), err)
	}

	return nil
}

// HandleOrderExpired 处理订单过期事件
func (h *EventHandler) HandleOrderExpired(ctx context.Context, eventData []byte) error {
	var event OrderEvent
	if err := json.Unmarshal(eventData, &event); err != nil {
		return fmt.Errorf("failed to unmarshal order expired event: %w", err)
	}

	// 释放库存
	if err := h.businessService.ReleaseInventoryWithOrderValidation(ctx, event.OrderID); err != nil {
		return fmt.Errorf("failed to release inventory for order %s: %w", event.OrderID.String(), err)
	}

	return nil
}

// HandleProductUpdated 处理商品更新事件
func (h *EventHandler) HandleProductUpdated(ctx context.Context, eventData []byte) error {
	// 商品更新时可能需要同步库存信息
	// 这里可以实现相关逻辑
	return nil
}

// HandleInventoryAlert 处理库存告警事件
func (h *EventHandler) HandleInventoryAlert(ctx context.Context, skuID uuid.UUID, alertType string, currentQuantity int32) error {
	// 这里可以实现告警通知逻辑
	// 例如发送邮件、短信、推送通知等

	// 获取商品信息
	inventoryWithProduct, err := h.businessService.GetInventoryWithProductInfo(ctx, skuID)
	if err != nil {
		return fmt.Errorf("failed to get inventory with product info: %w", err)
	}

	// 构建告警消息
	message := fmt.Sprintf("库存告警: 商品 %s (SKU: %s) %s, 当前库存: %d",
		inventoryWithProduct.ProductInfo.Name,
		inventoryWithProduct.ProductInfo.SKU,
		getAlertTypeMessage(alertType),
		currentQuantity)

	// TODO: 发送告警通知
	// - 发送邮件给管理员
	// - 发送消息到监控系统
	// - 记录告警日志

	fmt.Printf("Inventory Alert: %s\n", message)
	return nil
}

// getAlertTypeMessage 获取告警类型消息
func getAlertTypeMessage(alertType string) string {
	switch alertType {
	case "low_stock":
		return "库存不足"
	case "out_of_stock":
		return "库存售罄"
	default:
		return "库存异常"
	}
}

// CleanupExpiredReservations 清理过期预占记录（定时任务）
func (h *EventHandler) CleanupExpiredReservations(ctx context.Context) error {
	// TODO: 实现定时清理过期预占记录的逻辑
	// 这个方法可以被定时任务调用
	return nil
}
