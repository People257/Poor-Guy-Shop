package cart

import (
	"context"
	"time"

	"github.com/people257/poor-guy-shop/common/auth"
	"github.com/shopspring/decimal"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	pb "github.com/people257/poor-guy-shop/order-service/gen/proto/order/cart"
	cartapp "github.com/people257/poor-guy-shop/order-service/internal/application/cart"
	cartdomain "github.com/people257/poor-guy-shop/order-service/internal/domain/cart"
)

// GrpcHandler 购物车gRPC处理器
type GrpcHandler struct {
	pb.UnimplementedCartServiceServer
	cartService *cartapp.Service
}

// NewGrpcHandler 创建购物车gRPC处理器
func NewGrpcHandler(cartService *cartapp.Service) *GrpcHandler {
	return &GrpcHandler{
		cartService: cartService,
	}
}

// AddCartItem 添加商品到购物车
func (h *GrpcHandler) AddCartItem(ctx context.Context, req *pb.AddCartItemReq) (*pb.AddCartItemResp, error) {
	// 从认证上下文获取用户ID
	userID := auth.UserIDFromContext(ctx)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "用户未认证")
	}

	// TODO: 从product-service获取商品信息和价格
	// 这里暂时使用固定价格，实际应该调用product-service获取实时价格
	price := decimal.NewFromFloat(99.99)

	appReq := cartapp.AddToCartRequest{
		UserID:    userID,
		ProductID: req.ProductId,
		SkuID:     req.SkuId,
		Quantity:  req.Quantity,
		Price:     price,
	}

	cartItem, err := h.cartService.AddToCart(ctx, appReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "添加商品到购物车失败: %v", err)
	}

	return &pb.AddCartItemResp{
		Item: h.entityToProto(cartItem),
	}, nil
}

// UpdateCartItem 更新购物车商品
func (h *GrpcHandler) UpdateCartItem(ctx context.Context, req *pb.UpdateCartItemReq) (*pb.UpdateCartItemResp, error) {
	// 从认证上下文获取用户ID
	userID := auth.UserIDFromContext(ctx)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "用户未认证")
	}

	appReq := cartapp.UpdateQuantityRequest{
		CartID:   req.ItemId,
		UserID:   userID,
		Quantity: req.Quantity,
	}

	cartItem, err := h.cartService.UpdateQuantity(ctx, appReq)
	if err != nil {
		if err == cartdomain.ErrCartItemNotFound {
			return nil, status.Errorf(codes.NotFound, "购物车商品不存在")
		}
		return nil, status.Errorf(codes.Internal, "更新购物车商品失败: %v", err)
	}

	return &pb.UpdateCartItemResp{
		Item: h.entityToProto(cartItem),
	}, nil
}

// RemoveCartItem 删除购物车商品
func (h *GrpcHandler) RemoveCartItem(ctx context.Context, req *pb.RemoveCartItemReq) (*pb.RemoveCartItemResp, error) {
	// 从认证上下文获取用户ID
	userID := auth.UserIDFromContext(ctx)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "用户未认证")
	}

	appReq := cartapp.RemoveFromCartRequest{
		CartID: req.ItemId,
		UserID: userID,
	}

	err := h.cartService.RemoveFromCart(ctx, appReq)
	if err != nil {
		if err == cartdomain.ErrCartItemNotFound {
			return nil, status.Errorf(codes.NotFound, "购物车商品不存在")
		}
		return nil, status.Errorf(codes.Internal, "删除购物车商品失败: %v", err)
	}

	return &pb.RemoveCartItemResp{
		Success: true,
	}, nil
}

// GetCart 获取购物车
func (h *GrpcHandler) GetCart(ctx context.Context, req *pb.GetCartReq) (*pb.GetCartResp, error) {
	// 从认证上下文获取用户ID
	userID := auth.UserIDFromContext(ctx)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "用户未认证")
	}

	appReq := cartapp.GetCartRequest{
		UserID: userID,
	}

	result, err := h.cartService.GetCart(ctx, appReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "获取购物车失败: %v", err)
	}

	// 转换为proto对象
	var items []*pb.CartItem
	for _, item := range result.Items {
		items = append(items, h.entityToProto(item))
	}

	// 计算汇总信息
	summary := &pb.CartSummary{
		TotalItems:     int32(len(result.Items)),
		SelectedItems:  result.TotalCount,
		TotalAmount:    h.calculateTotalAmount(result.Items).String(),
		SelectedAmount: result.TotalAmount.String(),
	}

	return &pb.GetCartResp{
		Items:   items,
		Summary: summary,
	}, nil
}

// ClearCart 清空购物车
func (h *GrpcHandler) ClearCart(ctx context.Context, req *pb.ClearCartReq) (*pb.ClearCartResp, error) {
	// 从认证上下文获取用户ID
	userID := auth.UserIDFromContext(ctx)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "用户未认证")
	}

	appReq := cartapp.ClearCartRequest{
		UserID: userID,
	}

	err := h.cartService.ClearCart(ctx, appReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "清空购物车失败: %v", err)
	}

	return &pb.ClearCartResp{
		Success: true,
	}, nil
}

// SelectCartItems 选择购物车商品
func (h *GrpcHandler) SelectCartItems(ctx context.Context, req *pb.SelectCartItemsReq) (*pb.SelectCartItemsResp, error) {
	// 从认证上下文获取用户ID
	userID := auth.UserIDFromContext(ctx)
	if userID == "" {
		return nil, status.Error(codes.Unauthenticated, "用户未认证")
	}

	appReq := cartapp.BatchUpdateSelectionRequest{
		UserID:   userID,
		CartIDs:  req.ItemIds,
		Selected: req.Selected,
	}

	err := h.cartService.BatchUpdateSelection(ctx, appReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "选择购物车商品失败: %v", err)
	}

	// 重新获取购物车计算汇总
	cartReq := cartapp.GetCartRequest{UserID: userID}
	result, err := h.cartService.GetCart(ctx, cartReq)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "获取购物车汇总失败: %v", err)
	}

	summary := &pb.CartSummary{
		TotalItems:     int32(len(result.Items)),
		SelectedItems:  result.TotalCount,
		TotalAmount:    h.calculateTotalAmount(result.Items).String(),
		SelectedAmount: result.TotalAmount.String(),
	}

	return &pb.SelectCartItemsResp{
		Success: true,
		Summary: summary,
	}, nil
}

// entityToProto 将领域实体转换为proto对象
func (h *GrpcHandler) entityToProto(cartItem *cartdomain.ShoppingCart) *pb.CartItem {
	pbItem := &pb.CartItem{
		Id:        cartItem.ID,
		UserId:    cartItem.UserID,
		ProductId: cartItem.ProductID,
		SkuId:     cartItem.SkuID,
		Price:     cartItem.Price.String(),
		Quantity:  cartItem.Quantity,
		Selected:  cartItem.Selected,
		Available: true, // TODO: 从product-service获取库存状态
	}

	// 计算小计金额
	totalAmount := cartItem.Price.Mul(decimal.NewFromInt32(cartItem.Quantity))
	pbItem.TotalAmount = totalAmount.String()

	// TODO: 从product-service获取商品信息
	pbItem.ProductName = "商品名称"  // 实际应该从product-service获取
	pbItem.ProductImage = "商品图片" // 实际应该从product-service获取
	pbItem.SkuName = "SKU名称"     // 实际应该从product-service获取

	// 时间字段转换
	if cartItem.CreatedAt != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", cartItem.CreatedAt); err == nil {
			pbItem.CreatedAt = timestamppb.New(t)
		}
	}
	if cartItem.UpdatedAt != "" {
		if t, err := time.Parse("2006-01-02 15:04:05", cartItem.UpdatedAt); err == nil {
			pbItem.UpdatedAt = timestamppb.New(t)
		}
	}

	return pbItem
}

// calculateTotalAmount 计算总金额
func (h *GrpcHandler) calculateTotalAmount(items []*cartdomain.ShoppingCart) decimal.Decimal {
	var total decimal.Decimal
	for _, item := range items {
		subtotal := item.Price.Mul(decimal.NewFromInt32(item.Quantity))
		total = total.Add(subtotal)
	}
	return total
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
