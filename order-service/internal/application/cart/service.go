package cart

import (
	"context"
	"fmt"
	"time"

	"github.com/shopspring/decimal"

	"github.com/people257/poor-guy-shop/order-service/internal/domain/cart"
)

// Service 购物车应用服务
type Service struct {
	cartRepo cart.Repository
	cartDS   cart.DomainService
}

// NewService 创建购物车应用服务
func NewService(cartRepo cart.Repository, cartDS cart.DomainService) *Service {
	return &Service{
		cartRepo: cartRepo,
		cartDS:   cartDS,
	}
}

// AddToCartRequest 添加到购物车请求
type AddToCartRequest struct {
	UserID    string          `json:"user_id"`
	ProductID string          `json:"product_id"`
	SkuID     string          `json:"sku_id"`
	Quantity  int32           `json:"quantity"`
	Price     decimal.Decimal `json:"price"`
}

// AddToCart 添加商品到购物车
func (s *Service) AddToCart(ctx context.Context, req AddToCartRequest) (*cart.ShoppingCart, error) {
	// 检查购物车中是否已存在该商品
	existingItem, err := s.cartRepo.GetByUserAndProduct(ctx, req.UserID, req.ProductID, req.SkuID)
	if err != nil && err != cart.ErrCartItemNotFound {
		return nil, fmt.Errorf("检查购物车商品失败: %w", err)
	}

	if existingItem != nil {
		// 如果已存在，更新数量
		return s.cartDS.UpdateQuantity(ctx, existingItem, existingItem.Quantity+req.Quantity)
	}

	// 创建新的购物车项
	cartItem := &cart.ShoppingCart{
		UserID:    req.UserID,
		ProductID: req.ProductID,
		SkuID:     req.SkuID,
		Quantity:  req.Quantity,
		Price:     req.Price,
		Selected:  true,
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
		UpdatedAt: time.Now().Format("2006-01-02 15:04:05"),
	}

	return s.cartDS.AddToCart(ctx, cartItem)
}

// UpdateQuantityRequest 更新购物车商品数量请求
type UpdateQuantityRequest struct {
	CartID   string `json:"cart_id"`
	UserID   string `json:"user_id"`
	Quantity int32  `json:"quantity"`
}

// UpdateQuantity 更新购物车商品数量
func (s *Service) UpdateQuantity(ctx context.Context, req UpdateQuantityRequest) (*cart.ShoppingCart, error) {
	// 获取购物车项
	cartItem, err := s.cartRepo.GetByID(ctx, req.CartID)
	if err != nil {
		return nil, fmt.Errorf("获取购物车商品失败: %w", err)
	}

	// 检查是否属于该用户
	if cartItem.UserID != req.UserID {
		return nil, cart.ErrCartItemNotFound
	}

	// 使用领域服务更新数量
	return s.cartDS.UpdateQuantity(ctx, cartItem, req.Quantity)
}

// UpdateSelectionRequest 更新购物车商品选中状态请求
type UpdateSelectionRequest struct {
	CartID   string `json:"cart_id"`
	UserID   string `json:"user_id"`
	Selected bool   `json:"selected"`
}

// UpdateSelection 更新购物车商品选中状态
func (s *Service) UpdateSelection(ctx context.Context, req UpdateSelectionRequest) (*cart.ShoppingCart, error) {
	// 获取购物车项
	cartItem, err := s.cartRepo.GetByID(ctx, req.CartID)
	if err != nil {
		return nil, fmt.Errorf("获取购物车商品失败: %w", err)
	}

	// 检查是否属于该用户
	if cartItem.UserID != req.UserID {
		return nil, cart.ErrCartItemNotFound
	}

	// 使用领域服务更新选中状态
	return s.cartDS.UpdateSelection(ctx, cartItem, req.Selected)
}

// RemoveFromCartRequest 从购物车移除商品请求
type RemoveFromCartRequest struct {
	CartID string `json:"cart_id"`
	UserID string `json:"user_id"`
}

// RemoveFromCart 从购物车移除商品
func (s *Service) RemoveFromCart(ctx context.Context, req RemoveFromCartRequest) error {
	// 获取购物车项
	cartItem, err := s.cartRepo.GetByID(ctx, req.CartID)
	if err != nil {
		return fmt.Errorf("获取购物车商品失败: %w", err)
	}

	// 检查是否属于该用户
	if cartItem.UserID != req.UserID {
		return cart.ErrCartItemNotFound
	}

	// 使用领域服务移除商品
	return s.cartDS.RemoveFromCart(ctx, cartItem)
}

// GetCartRequest 获取购物车请求
type GetCartRequest struct {
	UserID string `json:"user_id"`
}

// GetCartResponse 获取购物车响应
type GetCartResponse struct {
	Items       []*cart.ShoppingCart `json:"items"`
	TotalAmount decimal.Decimal      `json:"total_amount"`
	TotalCount  int32                `json:"total_count"`
}

// GetCart 获取用户购物车
func (s *Service) GetCart(ctx context.Context, req GetCartRequest) (*GetCartResponse, error) {
	// 获取用户购物车商品
	items, err := s.cartRepo.GetByUserID(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("获取购物车失败: %w", err)
	}

	// 计算总金额和总数量
	var totalAmount decimal.Decimal
	var totalCount int32

	for _, item := range items {
		if item.Selected {
			subtotal := item.Price.Mul(decimal.NewFromInt32(item.Quantity))
			totalAmount = totalAmount.Add(subtotal)
			totalCount += item.Quantity
		}
	}

	return &GetCartResponse{
		Items:       items,
		TotalAmount: totalAmount,
		TotalCount:  totalCount,
	}, nil
}

// ClearCartRequest 清空购物车请求
type ClearCartRequest struct {
	UserID string `json:"user_id"`
}

// ClearCart 清空用户购物车
func (s *Service) ClearCart(ctx context.Context, req ClearCartRequest) error {
	return s.cartDS.ClearCart(ctx, req.UserID)
}

// BatchUpdateSelectionRequest 批量更新选中状态请求
type BatchUpdateSelectionRequest struct {
	UserID   string   `json:"user_id"`
	CartIDs  []string `json:"cart_ids"`
	Selected bool     `json:"selected"`
}

// BatchUpdateSelection 批量更新购物车商品选中状态
func (s *Service) BatchUpdateSelection(ctx context.Context, req BatchUpdateSelectionRequest) error {
	// 获取用户购物车商品
	items, err := s.cartRepo.GetByUserID(ctx, req.UserID)
	if err != nil {
		return fmt.Errorf("获取购物车失败: %w", err)
	}

	// 过滤出需要更新的商品
	var targetItems []*cart.ShoppingCart
	for _, item := range items {
		for _, cartID := range req.CartIDs {
			if item.ID == cartID {
				targetItems = append(targetItems, item)
				break
			}
		}
	}

	// 使用领域服务批量更新
	return s.cartDS.BatchUpdateSelection(ctx, targetItems, req.Selected)
}

// GetSelectedItemsRequest 获取选中商品请求
type GetSelectedItemsRequest struct {
	UserID string `json:"user_id"`
}

// GetSelectedItems 获取用户购物车中选中的商品
func (s *Service) GetSelectedItems(ctx context.Context, req GetSelectedItemsRequest) ([]*cart.ShoppingCart, error) {
	// 获取用户购物车商品
	items, err := s.cartRepo.GetByUserID(ctx, req.UserID)
	if err != nil {
		return nil, fmt.Errorf("获取购物车失败: %w", err)
	}

	// 过滤出选中的商品
	var selectedItems []*cart.ShoppingCart
	for _, item := range items {
		if item.Selected {
			selectedItems = append(selectedItems, item)
		}
	}

	return selectedItems, nil
}
