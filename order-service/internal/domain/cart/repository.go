package cart

import (
	"context"
)

// Repository 购物车仓储接口
type Repository interface {
	// 创建购物车项
	Create(ctx context.Context, cart *ShoppingCart) error

	// 根据ID获取购物车项
	GetByID(ctx context.Context, id string) (*ShoppingCart, error)

	// 根据用户ID获取购物车项列表
	GetByUserID(ctx context.Context, userID string) ([]*ShoppingCart, error)

	// 根据用户ID和商品ID获取购物车项
	GetByUserAndProduct(ctx context.Context, userID, productID, skuID string) (*ShoppingCart, error)

	// 更新购物车项
	Update(ctx context.Context, cart *ShoppingCart) error

	// 删除购物车项
	Delete(ctx context.Context, id string) error

	// 删除用户的所有购物车项
	DeleteByUserID(ctx context.Context, userID string) error

	// 批量更新购物车项
	BatchUpdate(ctx context.Context, carts []*ShoppingCart) error

	// 获取用户购物车中选中的商品
	GetSelectedItems(ctx context.Context, userID string) ([]*ShoppingCart, error)

	// 统计用户购物车商品数量
	CountByUserID(ctx context.Context, userID string) (int64, error)
}
