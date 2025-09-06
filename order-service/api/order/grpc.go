package order

import (
	"context"
	"strconv"
	"time"

	"github.com/people257/poor-guy-shop/common/auth"

	"github.com/shopspring/decimal"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/people257/poor-guy-shop/order-service/gen/proto/order/order"
	orderapp "github.com/people257/poor-guy-shop/order-service/internal/application/order"
	orderdomain "github.com/people257/poor-guy-shop/order-service/internal/domain/order"
)

// GrpcHandler 订单gRPC处理器
type GrpcHandler struct {
	pb.UnimplementedOrderServiceServer
	orderService *orderapp.Service
}

// NewGrpcHandler 创建订单gRPC处理器
func NewGrpcHandler(orderService *orderapp.Service) *GrpcHandler {
	return &GrpcHandler{
		orderService: orderService,
	}
}

// CreateOrder 创建订单
func (h *GrpcHandler) CreateOrder(ctx context.Context, req *pb.CreateOrderReq) (*pb.CreateOrderResp, error) {
	// 从认证上下文获取用户ID
	userID := auth.UserIDFromContext(ctx)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "用户未认证")
	}
	// 构建应用层请求
	var items []orderapp.CreateOrderItemRequest
	for _, item := range req.Items {
		price, err := decimal.NewFromString(item.Price)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "商品价格格式错误: %v", err)
		}

		items = append(items, orderapp.CreateOrderItemRequest{
			ProductID:   item.ProductId,
			SkuID:       item.SkuId,
			Quantity:    item.Quantity,
			Price:       price,
			ProductName: item.ProductName,
			SkuName:     item.SkuName,
		})
	}

	discountAmount := decimal.Zero
	if req.DiscountAmount != "" {
		var err error
		discountAmount, err = decimal.NewFromString(req.DiscountAmount)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "优惠金额格式错误: %v", err)
		}
	}

	shippingFee := decimal.Zero
	if req.ShippingFee != "" {
		var err error
		shippingFee, err = decimal.NewFromString(req.ShippingFee)
		if err != nil {
			return nil, status.Errorf(codes.InvalidArgument, "运费格式错误: %v", err)
		}
	}

	appReq := orderapp.CreateOrderRequest{
		UserID: userID, // 使用从上下文获取的用户ID
		Items:  items,
		Address: orderapp.CreateOrderAddressRequest{
			ReceiverName:  req.Address.ReceiverName,
			ReceiverPhone: req.Address.ReceiverPhone,
			Province:      req.Address.Province,
			City:          req.Address.City,
			District:      req.Address.District,
			DetailAddress: req.Address.DetailAddress,
			PostalCode:    req.Address.PostalCode,
		},
		PaymentMethod:  req.PaymentMethod,
		Remark:         req.Remark,
		DiscountAmount: discountAmount,
		ShippingFee:    shippingFee,
	}

	// 调用应用服务
	orderEntity, err := h.orderService.CreateOrder(ctx, appReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "创建订单失败: %v", err)
	}

	// 转换为响应
	return &pb.CreateOrderResp{
		Order: h.entityToProto(orderEntity),
	}, nil
}

// GetOrder 获取订单详情
func (h *GrpcHandler) GetOrder(ctx context.Context, req *pb.GetOrderReq) (*pb.GetOrderResp, error) {
	// 从认证上下文获取用户ID
	userID := auth.UserIDFromContext(ctx)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "用户未认证")
	}
	appReq := orderapp.GetOrderRequest{
		OrderID: req.OrderId,
		UserID:  userID, // 使用从上下文获取的用户ID
	}

	orderEntity, err := h.orderService.GetOrder(ctx, appReq)
	if err != nil {
		if err == orderdomain.ErrOrderNotFound {
			return nil, status.Errorf(codes.NotFound, "订单不存在")
		}
		return nil, status.Errorf(codes.Internal, "获取订单失败: %v", err)
	}

	return &pb.GetOrderResp{
		Order: h.entityToProto(orderEntity),
	}, nil
}

// ListOrders 获取订单列表
func (h *GrpcHandler) ListOrders(ctx context.Context, req *pb.ListOrdersReq) (*pb.ListOrdersResp, error) {
	// 从认证上下文获取用户ID
	userID := auth.UserIDFromContext(ctx)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "用户未认证")
	}

	appReq := orderapp.ListOrdersRequest{
		UserID:   userID,
		Status:   req.Status,
		Page:     req.Page,
		PageSize: req.PageSize,
	}

	result, err := h.orderService.ListOrders(ctx, appReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "获取订单列表失败: %v", err)
	}

	// 转换为proto对象
	var orders []*pb.Order
	for _, orderEntity := range result.Orders {
		orders = append(orders, h.entityToProto(orderEntity))
	}

	return &pb.ListOrdersResp{
		Orders:     orders,
		Total:      result.Total,
		Page:       result.Page,
		PageSize:   result.PageSize,
		TotalPages: result.TotalPages,
	}, nil
}

// UpdateOrderStatus 更新订单状态
func (h *GrpcHandler) UpdateOrderStatus(ctx context.Context, req *pb.UpdateOrderStatusReq) (*pb.UpdateOrderStatusResp, error) {
	// 从认证上下文获取用户ID
	userID := auth.UserIDFromContext(ctx)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "用户未认证")
	}

	appReq := orderapp.UpdateOrderStatusRequest{
		OrderID: req.OrderId,
		UserID:  userID,
		Status:  req.Status,
		Reason:  req.Reason,
	}

	err := h.orderService.UpdateOrderStatus(ctx, appReq)
	if err != nil {
		if err == orderdomain.ErrOrderNotFound {
			return nil, status.Errorf(codes.NotFound, "订单不存在")
		}
		return nil, status.Errorf(codes.Internal, "更新订单状态失败: %v", err)
	}

	return &pb.UpdateOrderStatusResp{
		Success: true,
	}, nil
}

// CancelOrder 取消订单
func (h *GrpcHandler) CancelOrder(ctx context.Context, req *pb.CancelOrderReq) (*pb.CancelOrderResp, error) {
	// 从认证上下文获取用户ID
	userID := auth.UserIDFromContext(ctx)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "用户未认证")
	}

	appReq := orderapp.CancelOrderRequest{
		OrderID: req.OrderId,
		UserID:  userID,
		Reason:  req.Reason,
	}

	err := h.orderService.CancelOrder(ctx, appReq)
	if err != nil {
		if err == orderdomain.ErrOrderNotFound {
			return nil, status.Errorf(codes.NotFound, "订单不存在")
		}
		return nil, status.Errorf(codes.Internal, "取消订单失败: %v", err)
	}

	return &pb.CancelOrderResp{
		Success: true,
	}, nil
}

// PayOrder 支付订单
func (h *GrpcHandler) PayOrder(ctx context.Context, req *pb.PayOrderReq) (*pb.PayOrderResp, error) {
	// 从认证上下文获取用户ID
	userID := auth.UserIDFromContext(ctx)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "用户未认证")
	}

	appReq := orderapp.PayOrderRequest{
		OrderID:       req.OrderId,
		UserID:        userID,
		PaymentMethod: req.PaymentMethod,
	}

	err := h.orderService.PayOrder(ctx, appReq)
	if err != nil {
		if err == orderdomain.ErrOrderNotFound {
			return nil, status.Errorf(codes.NotFound, "订单不存在")
		}
		return nil, status.Errorf(codes.Internal, "支付订单失败: %v", err)
	}

	return &pb.PayOrderResp{
		Success: true,
	}, nil
}

// entityToProto 将领域实体转换为proto对象
func (h *GrpcHandler) entityToProto(orderEntity *orderdomain.Order) *pb.Order {
	pbOrder := &pb.Order{
		Id:            orderEntity.ID,
		OrderNo:       orderEntity.OrderNo,
		UserId:        orderEntity.UserID,
		Status:        pb.OrderStatus(orderEntity.Status),
		TotalAmount:   orderEntity.TotalAmount.String(),
		ActualAmount:  orderEntity.ActualAmount.String(),
		PaymentMethod: pb.PaymentMethod(pb.PaymentMethod_value[orderEntity.PaymentMethod]),
		PaymentStatus: orderEntity.PaymentStatus,
		Remark:        orderEntity.Remark,
	}

	if !orderEntity.DiscountAmount.IsZero() {
		pbOrder.DiscountAmount = orderEntity.DiscountAmount.String()
	}
	if !orderEntity.ShippingFee.IsZero() {
		pbOrder.ShippingFee = orderEntity.ShippingFee.String()
	}
	if orderEntity.CancelReason != "" {
		pbOrder.CancelReason = orderEntity.CancelReason
	}

	// 时间字段转换
	if orderEntity.PaymentTime != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", orderEntity.PaymentTime); err == nil {
			pbOrder.PaymentTime = timestamppb.New(t)
		}
	}
	if orderEntity.DeliveryTime != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", orderEntity.DeliveryTime); err == nil {
			pbOrder.DeliveryTime = timestamppb.New(t)
		}
	}
	if orderEntity.ReceiveTime != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", orderEntity.ReceiveTime); err == nil {
			pbOrder.ReceiveTime = timestamppb.New(t)
		}
	}
	if orderEntity.CancelTime != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", orderEntity.CancelTime); err == nil {
			pbOrder.CancelTime = timestamppb.New(t)
		}
	}
	if orderEntity.CreatedAt != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", orderEntity.CreatedAt); err == nil {
			pbOrder.CreatedAt = timestamppb.New(t)
		}
	}
	if orderEntity.UpdatedAt != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", orderEntity.UpdatedAt); err == nil {
			pbOrder.UpdatedAt = timestamppb.New(t)
		}
	}

	return pbOrder
}

// parseTime 解析时间字符串
func (h *GrpcHandler) parseTime(timeStr string) *timestamppb.Timestamp {
	if timeStr == "" {
		return nil
	}

	t, err := time.Parse("2006-01-02 15:04:05", timeStr)
	if err != nil {
		return nil
	}

	return timestamppb.New(t)
}

// parseDecimal 解析decimal字符串
func (h *GrpcHandler) parseDecimal(decimalStr string) (decimal.Decimal, error) {
	if decimalStr == "" {
		return decimal.Zero, nil
	}
	return decimal.NewFromString(decimalStr)
}

// parseInt32 解析int32字符串
func (h *GrpcHandler) parseInt32(intStr string) (int32, error) {
	if intStr == "" {
		return 0, nil
	}
	val, err := strconv.ParseInt(intStr, 10, 32)
	if err != nil {
		return 0, err
	}
	return int32(val), nil
}
