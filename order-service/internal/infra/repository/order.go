package repository

import (
	"context"
	"fmt"
	"time"

	"gorm.io/gorm"

	"github.com/people257/poor-guy-shop/order-service/gen/gen/model"
	"github.com/people257/poor-guy-shop/order-service/gen/gen/query"
	"github.com/people257/poor-guy-shop/order-service/internal/domain/order"
)

// orderRepository 订单仓储实现
type orderRepository struct {
	db    *gorm.DB
	query *query.Query
}

// NewOrderRepository 创建订单仓储
func NewOrderRepository(db *gorm.DB, q *query.Query) order.Repository {
	return &orderRepository{
		db:    db,
		query: q,
	}
}

// Create 创建订单
func (r *orderRepository) Create(ctx context.Context, orderEntity *order.Order, items []*order.OrderItem, address *order.OrderAddress) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// 1. 创建订单主记录
		orderModel := r.domainToModel(orderEntity)
		if err := tx.Create(orderModel).Error; err != nil {
			return fmt.Errorf("创建订单失败: %w", err)
		}

		// 设置订单ID到实体
		orderEntity.ID = orderModel.ID
		orderEntity.OrderNo = orderModel.OrderNo

		// 2. 创建订单商品项
		for _, item := range items {
			itemModel := r.itemDomainToModel(item)
			itemModel.OrderID = orderModel.ID
			if err := tx.Create(itemModel).Error; err != nil {
				return fmt.Errorf("创建订单商品项失败: %w", err)
			}
			item.ID = itemModel.ID
			item.OrderID = itemModel.OrderID
		}

		// 3. 创建订单地址
		addressModel := r.addressDomainToModel(address)
		addressModel.OrderID = orderModel.ID
		if err := tx.Create(addressModel).Error; err != nil {
			return fmt.Errorf("创建订单地址失败: %w", err)
		}
		address.ID = addressModel.ID
		address.OrderID = addressModel.OrderID

		// 4. 记录状态日志
		remark := "订单创建"
		statusLog := &model.OrderStatusLog{
			OrderID:  orderModel.ID,
			ToStatus: orderModel.Status,
			Remark:   &remark,
		}
		if err := tx.Create(statusLog).Error; err != nil {
			return fmt.Errorf("创建状态日志失败: %w", err)
		}

		return nil
	})
}

// GetByID 根据ID获取订单
func (r *orderRepository) GetByID(ctx context.Context, id string) (*order.Order, error) {
	orderModel, err := r.query.WithContext(ctx).Order.Where(r.query.Order.ID.Eq(id)).First()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, order.ErrOrderNotFound
		}
		return nil, fmt.Errorf("获取订单失败: %w", err)
	}

	return r.modelToDomain(orderModel), nil
}

// GetByOrderNo 根据订单号获取订单
func (r *orderRepository) GetByOrderNo(ctx context.Context, orderNo string) (*order.Order, error) {
	orderModel, err := r.query.WithContext(ctx).Order.Where(r.query.Order.OrderNo.Eq(orderNo)).First()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, order.ErrOrderNotFound
		}
		return nil, fmt.Errorf("获取订单失败: %w", err)
	}

	return r.modelToDomain(orderModel), nil
}

// ListByUserID 根据用户ID获取订单列表
func (r *orderRepository) ListByUserID(ctx context.Context, userID string, status int32, page, pageSize int32) ([]*order.Order, int64, error) {
	q := r.query.WithContext(ctx).Order.Where(r.query.Order.UserID.Eq(userID))

	if status > 0 {
		q = q.Where(r.query.Order.Status.Eq(status))
	}

	// 获取总数
	total, err := q.Count()
	if err != nil {
		return nil, 0, fmt.Errorf("获取订单总数失败: %w", err)
	}

	// 分页查询
	offset := (page - 1) * pageSize
	orderModels, err := q.Order(r.query.Order.CreatedAt.Desc()).Limit(int(pageSize)).Offset(int(offset)).Find()
	if err != nil {
		return nil, 0, fmt.Errorf("获取订单列表失败: %w", err)
	}

	// 转换为领域对象
	var orders []*order.Order
	for _, orderModel := range orderModels {
		orders = append(orders, r.modelToDomain(orderModel))
	}

	return orders, total, nil
}

// Update 更新订单
func (r *orderRepository) Update(ctx context.Context, orderEntity *order.Order) error {
	orderModel := r.domainToModel(orderEntity)

	_, err := r.query.WithContext(ctx).Order.Where(r.query.Order.ID.Eq(orderEntity.ID)).Updates(orderModel)
	if err != nil {
		return fmt.Errorf("更新订单失败: %w", err)
	}

	return nil
}

// Delete 删除订单（软删除）
func (r *orderRepository) Delete(ctx context.Context, id string) error {
	_, err := r.query.WithContext(ctx).Order.Where(r.query.Order.ID.Eq(id)).Delete()
	if err != nil {
		return fmt.Errorf("删除订单失败: %w", err)
	}

	return nil
}

// GetOrderItems 获取订单商品项
func (r *orderRepository) GetOrderItems(ctx context.Context, orderID string) ([]*order.OrderItem, error) {
	itemModels, err := r.query.WithContext(ctx).OrderItem.Where(r.query.OrderItem.OrderID.Eq(orderID)).Find()
	if err != nil {
		return nil, fmt.Errorf("获取订单商品项失败: %w", err)
	}

	var items []*order.OrderItem
	for _, itemModel := range itemModels {
		items = append(items, r.itemModelToDomain(itemModel))
	}

	return items, nil
}

// GetOrderAddress 获取订单地址
func (r *orderRepository) GetOrderAddress(ctx context.Context, orderID string) (*order.OrderAddress, error) {
	addressModel, err := r.query.WithContext(ctx).OrderAddress.Where(r.query.OrderAddress.OrderID.Eq(orderID)).First()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, order.ErrOrderAddressNotFound
		}
		return nil, fmt.Errorf("获取订单地址失败: %w", err)
	}

	return r.addressModelToDomain(addressModel), nil
}

// CreateStatusLog 创建状态日志
func (r *orderRepository) CreateStatusLog(ctx context.Context, orderID string, status int32, remark string) error {
	remarkPtr := remark
	statusLog := &model.OrderStatusLog{
		OrderID:  orderID,
		ToStatus: status,
		Remark:   &remarkPtr,
	}

	if err := r.db.WithContext(ctx).Create(statusLog).Error; err != nil {
		return fmt.Errorf("创建状态日志失败: %w", err)
	}

	return nil
}

// GetStatusLogs 获取订单状态日志
func (r *orderRepository) GetStatusLogs(ctx context.Context, orderID string) ([]*order.OrderStatusLog, error) {
	logModels, err := r.query.WithContext(ctx).OrderStatusLog.Where(r.query.OrderStatusLog.OrderID.Eq(orderID)).Order(r.query.OrderStatusLog.CreatedAt.Desc()).Find()
	if err != nil {
		return nil, fmt.Errorf("获取状态日志失败: %w", err)
	}

	var logs []*order.OrderStatusLog
	for _, logModel := range logModels {
		logs = append(logs, r.statusLogModelToDomain(logModel))
	}

	return logs, nil
}

// 领域对象转换为数据模型
func (r *orderRepository) domainToModel(orderEntity *order.Order) *model.Order {
	orderModel := &model.Order{
		ID:            orderEntity.ID,
		OrderNo:       orderEntity.OrderNo,
		UserID:        orderEntity.UserID,
		Status:        orderEntity.Status,
		TotalAmount:   orderEntity.TotalAmount,
		ActualAmount:  orderEntity.ActualAmount,
		PaymentMethod: &orderEntity.PaymentMethod,
		PaymentStatus: &orderEntity.PaymentStatus,
		Remark:        &orderEntity.Remark,
		Version:       orderEntity.Version,
	}

	if !orderEntity.DiscountAmount.IsZero() {
		orderModel.DiscountAmount = &orderEntity.DiscountAmount
	}
	if !orderEntity.ShippingFee.IsZero() {
		orderModel.ShippingFee = &orderEntity.ShippingFee
	}

	// 处理时间字段
	if orderEntity.PaymentTime != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", orderEntity.PaymentTime); err == nil {
			orderModel.PaymentTime = &t
		}
	}
	if orderEntity.DeliveryTime != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", orderEntity.DeliveryTime); err == nil {
			orderModel.DeliveryTime = &t
		}
	}
	if orderEntity.ReceiveTime != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", orderEntity.ReceiveTime); err == nil {
			orderModel.ReceiveTime = &t
		}
	}
	if orderEntity.CancelTime != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", orderEntity.CancelTime); err == nil {
			orderModel.CancelTime = &t
		}
	}

	if orderEntity.CancelReason != "" {
		orderModel.CancelReason = &orderEntity.CancelReason
	}

	return orderModel
}

// 数据模型转换为领域对象
func (r *orderRepository) modelToDomain(orderModel *model.Order) *order.Order {
	orderEntity := &order.Order{
		ID:           orderModel.ID,
		OrderNo:      orderModel.OrderNo,
		UserID:       orderModel.UserID,
		Status:       orderModel.Status,
		TotalAmount:  orderModel.TotalAmount,
		ActualAmount: orderModel.ActualAmount,
		Version:      orderModel.Version,
		CreatedAt:    orderModel.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:    orderModel.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if orderModel.DiscountAmount != nil {
		orderEntity.DiscountAmount = *orderModel.DiscountAmount
	}
	if orderModel.ShippingFee != nil {
		orderEntity.ShippingFee = *orderModel.ShippingFee
	}
	if orderModel.PaymentMethod != nil {
		orderEntity.PaymentMethod = *orderModel.PaymentMethod
	}
	if orderModel.PaymentStatus != nil {
		orderEntity.PaymentStatus = *orderModel.PaymentStatus
	}
	if orderModel.Remark != nil {
		orderEntity.Remark = *orderModel.Remark
	}

	// 处理时间字段
	if orderModel.PaymentTime != nil {
		orderEntity.PaymentTime = orderModel.PaymentTime.Format("2006-01-02 15:04:05")
	}
	if orderModel.DeliveryTime != nil {
		orderEntity.DeliveryTime = orderModel.DeliveryTime.Format("2006-01-02 15:04:05")
	}
	if orderModel.ReceiveTime != nil {
		orderEntity.ReceiveTime = orderModel.ReceiveTime.Format("2006-01-02 15:04:05")
	}
	if orderModel.CancelTime != nil {
		orderEntity.CancelTime = orderModel.CancelTime.Format("2006-01-02 15:04:05")
	}
	if orderModel.CancelReason != nil {
		orderEntity.CancelReason = *orderModel.CancelReason
	}

	return orderEntity
}

// 订单商品项领域对象转换为数据模型
func (r *orderRepository) itemDomainToModel(item *order.OrderItem) *model.OrderItem {
	return &model.OrderItem{
		ID:          item.ID,
		OrderID:     item.OrderID,
		ProductID:   item.ProductID,
		SkuID:       &item.SkuID,
		Quantity:    item.Quantity,
		Price:       item.Price,
		ProductName: item.ProductName,
		SkuName:     &item.SkuName,
	}
}

// 订单商品项数据模型转换为领域对象
func (r *orderRepository) itemModelToDomain(itemModel *model.OrderItem) *order.OrderItem {
	item := &order.OrderItem{
		ID:          itemModel.ID,
		OrderID:     itemModel.OrderID,
		ProductID:   itemModel.ProductID,
		Quantity:    itemModel.Quantity,
		Price:       itemModel.Price,
		ProductName: itemModel.ProductName,
		CreatedAt:   itemModel.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:   itemModel.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if itemModel.SkuID != nil {
		item.SkuID = *itemModel.SkuID
	}
	if itemModel.SkuName != nil {
		item.SkuName = *itemModel.SkuName
	}

	return item
}

// 订单地址领域对象转换为数据模型
func (r *orderRepository) addressDomainToModel(address *order.OrderAddress) *model.OrderAddress {
	return &model.OrderAddress{
		ID:            address.ID,
		OrderID:       address.OrderID,
		ReceiverName:  address.ReceiverName,
		ReceiverPhone: address.ReceiverPhone,
		Province:      address.Province,
		City:          address.City,
		District:      address.District,
		Address:       address.DetailAddress,
		PostalCode:    &address.PostalCode,
	}
}

// 订单地址数据模型转换为领域对象
func (r *orderRepository) addressModelToDomain(addressModel *model.OrderAddress) *order.OrderAddress {
	address := &order.OrderAddress{
		ID:            addressModel.ID,
		OrderID:       addressModel.OrderID,
		ReceiverName:  addressModel.ReceiverName,
		ReceiverPhone: addressModel.ReceiverPhone,
		Province:      addressModel.Province,
		City:          addressModel.City,
		District:      addressModel.District,
		DetailAddress: addressModel.Address,
		CreatedAt:     addressModel.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt:     addressModel.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if addressModel.PostalCode != nil {
		address.PostalCode = *addressModel.PostalCode
	}

	return address
}

// 状态日志数据模型转换为领域对象
func (r *orderRepository) statusLogModelToDomain(logModel *model.OrderStatusLog) *order.OrderStatusLog {
	log := &order.OrderStatusLog{
		ID:        logModel.ID,
		OrderID:   logModel.OrderID,
		Status:    logModel.ToStatus,
		CreatedAt: logModel.CreatedAt.Format("2006-01-02 15:04:05"),
	}

	if logModel.Remark != nil {
		log.Remark = *logModel.Remark
	}

	return log
}
