package cart

import (
	"context"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
)

// DomainService 购物车领域服务
type DomainService interface {
	// 添加商品到购物车
	AddToCart(ctx context.Context, cart *ShoppingCart) (*ShoppingCart, error)

	// 更新购物车商品数量
	UpdateQuantity(ctx context.Context, cart *ShoppingCart, quantity int32) (*ShoppingCart, error)

	// 更新购物车商品选中状态
	UpdateSelection(ctx context.Context, cart *ShoppingCart, selected bool) (*ShoppingCart, error)

	// 从购物车移除商品
	RemoveFromCart(ctx context.Context, cart *ShoppingCart) error

	// 清空购物车
	ClearCart(ctx context.Context, userID string) error

	// 批量更新选中状态
	BatchUpdateSelection(ctx context.Context, carts []*ShoppingCart, selected bool) error
}

// domainService 购物车领域服务实现
type domainService struct {
	cartRepo Repository
}

// NewDomainService 创建购物车领域服务
func NewDomainService(cartRepo Repository) DomainService {
	return &domainService{
		cartRepo: cartRepo,
	}
}

// AddToCart 添加商品到购物车
func (ds *domainService) AddToCart(ctx context.Context, cart *ShoppingCart) (*ShoppingCart, error) {
	// 验证商品数量
	if cart.Quantity <= 0 {
		return nil, fmt.Errorf("商品数量必须大于0")
	}

	// 验证商品价格
	if cart.Price.LessThanOrEqual(decimal.Zero) {
		return nil, fmt.Errorf("商品价格必须大于0")
	}

	// 设置默认值
	cart.Selected = true
	cart.CreatedAt = time.Now().Format("2006-01-02 15:04:05")
	cart.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	cart.Version = 1

	// 调用仓储创建购物车项
	if err := ds.cartRepo.Create(ctx, cart); err != nil {
		return nil, fmt.Errorf("添加商品到购物车失败: %w", err)
	}

	return cart, nil
}

// UpdateQuantity 更新购物车商品数量
func (ds *domainService) UpdateQuantity(ctx context.Context, cart *ShoppingCart, quantity int32) (*ShoppingCart, error) {
	// 验证数量
	if quantity <= 0 {
		return nil, fmt.Errorf("商品数量必须大于0")
	}

	// 更新数量和时间
	cart.Quantity = quantity
	cart.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	cart.Version++

	// 调用仓储更新
	if err := ds.cartRepo.Update(ctx, cart); err != nil {
		return nil, fmt.Errorf("更新购物车商品数量失败: %w", err)
	}

	return cart, nil
}

// UpdateSelection 更新购物车商品选中状态
func (ds *domainService) UpdateSelection(ctx context.Context, cart *ShoppingCart, selected bool) (*ShoppingCart, error) {
	// 更新选中状态和时间
	cart.Selected = selected
	cart.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
	cart.Version++

	// 调用仓储更新
	if err := ds.cartRepo.Update(ctx, cart); err != nil {
		return nil, fmt.Errorf("更新购物车商品选中状态失败: %w", err)
	}

	return cart, nil
}

// RemoveFromCart 从购物车移除商品
func (ds *domainService) RemoveFromCart(ctx context.Context, cart *ShoppingCart) error {
	// 调用仓储删除
	if err := ds.cartRepo.Delete(ctx, cart.ID); err != nil {
		return fmt.Errorf("从购物车移除商品失败: %w", err)
	}

	return nil
}

// ClearCart 清空购物车
func (ds *domainService) ClearCart(ctx context.Context, userID string) error {
	// 调用仓储删除用户所有购物车项
	if err := ds.cartRepo.DeleteByUserID(ctx, userID); err != nil {
		return fmt.Errorf("清空购物车失败: %w", err)
	}

	return nil
}

// BatchUpdateSelection 批量更新选中状态
func (ds *domainService) BatchUpdateSelection(ctx context.Context, carts []*ShoppingCart, selected bool) error {
	// 更新每个购物车项的选中状态
	for _, cart := range carts {
		cart.Selected = selected
		cart.UpdatedAt = time.Now().Format("2006-01-02 15:04:05")
		cart.Version++
	}

	// 批量更新
	if err := ds.cartRepo.BatchUpdate(ctx, carts); err != nil {
		return fmt.Errorf("批量更新购物车选中状态失败: %w", err)
	}

	return nil
}
