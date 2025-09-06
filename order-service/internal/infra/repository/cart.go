package repository

import (
	"context"
	"fmt"

	"gorm.io/gorm"

	"github.com/people257/poor-guy-shop/order-service/gen/gen/model"
	"github.com/people257/poor-guy-shop/order-service/gen/gen/query"
	"github.com/people257/poor-guy-shop/order-service/internal/domain/cart"
)

// cartRepository 购物车仓储实现
type cartRepository struct {
	db    *gorm.DB
	query *query.Query
}

// NewCartRepository 创建购物车仓储
func NewCartRepository(db *gorm.DB, q *query.Query) cart.Repository {
	return &cartRepository{
		db:    db,
		query: q,
	}
}

// Create 创建购物车项
func (r *cartRepository) Create(ctx context.Context, cartEntity *cart.ShoppingCart) error {
	cartModel := r.domainToModel(cartEntity)

	if err := r.db.WithContext(ctx).Create(cartModel).Error; err != nil {
		return fmt.Errorf("创建购物车项失败: %w", err)
	}

	// 设置生成的ID
	cartEntity.ID = cartModel.ID

	return nil
}

// GetByID 根据ID获取购物车项
func (r *cartRepository) GetByID(ctx context.Context, id string) (*cart.ShoppingCart, error) {
	cartModel, err := r.query.WithContext(ctx).ShoppingCart.Where(r.query.ShoppingCart.ID.Eq(id)).First()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, cart.ErrCartItemNotFound
		}
		return nil, fmt.Errorf("获取购物车项失败: %w", err)
	}

	return r.modelToDomain(cartModel), nil
}

// GetByUserID 根据用户ID获取购物车项列表
func (r *cartRepository) GetByUserID(ctx context.Context, userID string) ([]*cart.ShoppingCart, error) {
	cartModels, err := r.query.WithContext(ctx).ShoppingCart.Where(r.query.ShoppingCart.UserID.Eq(userID)).Order(r.query.ShoppingCart.CreatedAt.Desc()).Find()
	if err != nil {
		return nil, fmt.Errorf("获取购物车列表失败: %w", err)
	}

	var cartItems []*cart.ShoppingCart
	for _, cartModel := range cartModels {
		cartItems = append(cartItems, r.modelToDomain(cartModel))
	}

	return cartItems, nil
}

// GetByUserAndProduct 根据用户ID和商品ID获取购物车项
func (r *cartRepository) GetByUserAndProduct(ctx context.Context, userID, productID, skuID string) (*cart.ShoppingCart, error) {
	q := r.query.WithContext(ctx).ShoppingCart.Where(
		r.query.ShoppingCart.UserID.Eq(userID),
		r.query.ShoppingCart.ProductID.Eq(productID),
	)

	if skuID != "" {
		q = q.Where(r.query.ShoppingCart.SkuID.Eq(skuID))
	} else {
		q = q.Where(r.query.ShoppingCart.SkuID.IsNull())
	}

	cartModel, err := q.First()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, cart.ErrCartItemNotFound
		}
		return nil, fmt.Errorf("获取购物车项失败: %w", err)
	}

	return r.modelToDomain(cartModel), nil
}

// Update 更新购物车项
func (r *cartRepository) Update(ctx context.Context, cartEntity *cart.ShoppingCart) error {
	cartModel := r.domainToModel(cartEntity)

	_, err := r.query.WithContext(ctx).ShoppingCart.Where(r.query.ShoppingCart.ID.Eq(cartEntity.ID)).Updates(cartModel)
	if err != nil {
		return fmt.Errorf("更新购物车项失败: %w", err)
	}

	return nil
}

// Delete 删除购物车项
func (r *cartRepository) Delete(ctx context.Context, id string) error {
	_, err := r.query.WithContext(ctx).ShoppingCart.Where(r.query.ShoppingCart.ID.Eq(id)).Delete()
	if err != nil {
		return fmt.Errorf("删除购物车项失败: %w", err)
	}

	return nil
}

// DeleteByUserID 删除用户的所有购物车项
func (r *cartRepository) DeleteByUserID(ctx context.Context, userID string) error {
	_, err := r.query.WithContext(ctx).ShoppingCart.Where(r.query.ShoppingCart.UserID.Eq(userID)).Delete()
	if err != nil {
		return fmt.Errorf("清空购物车失败: %w", err)
	}

	return nil
}

// BatchUpdate 批量更新购物车项
func (r *cartRepository) BatchUpdate(ctx context.Context, cartItems []*cart.ShoppingCart) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, cartItem := range cartItems {
			cartModel := r.domainToModel(cartItem)
			if err := tx.Where("id = ?", cartItem.ID).Updates(cartModel).Error; err != nil {
				return fmt.Errorf("批量更新购物车项失败: %w", err)
			}
		}
		return nil
	})
}

// GetSelectedItems 获取用户购物车中选中的商品
func (r *cartRepository) GetSelectedItems(ctx context.Context, userID string) ([]*cart.ShoppingCart, error) {
	cartModels, err := r.query.WithContext(ctx).ShoppingCart.Where(
		r.query.ShoppingCart.UserID.Eq(userID),
		r.query.ShoppingCart.Selected.Is(true),
	).Order(r.query.ShoppingCart.CreatedAt.Desc()).Find()

	if err != nil {
		return nil, fmt.Errorf("获取选中购物车项失败: %w", err)
	}

	var cartItems []*cart.ShoppingCart
	for _, cartModel := range cartModels {
		cartItems = append(cartItems, r.modelToDomain(cartModel))
	}

	return cartItems, nil
}

// CountByUserID 统计用户购物车商品数量
func (r *cartRepository) CountByUserID(ctx context.Context, userID string) (int64, error) {
	count, err := r.query.WithContext(ctx).ShoppingCart.Where(r.query.ShoppingCart.UserID.Eq(userID)).Count()
	if err != nil {
		return 0, fmt.Errorf("统计购物车商品数量失败: %w", err)
	}

	return count, nil
}

// 领域对象转换为数据模型
func (r *cartRepository) domainToModel(cartEntity *cart.ShoppingCart) *model.ShoppingCart {
	cartModel := &model.ShoppingCart{
		ID:        cartEntity.ID,
		UserID:    cartEntity.UserID,
		ProductID: cartEntity.ProductID,
		Quantity:  cartEntity.Quantity,
		Price:     cartEntity.Price,
		Selected:  cartEntity.Selected,
		Version:   cartEntity.Version,
	}

	if cartEntity.SkuID != "" {
		cartModel.SkuID = &cartEntity.SkuID
	}

	return cartModel
}

// 数据模型转换为领域对象
func (r *cartRepository) modelToDomain(cartModel *model.ShoppingCart) *cart.ShoppingCart {
	cartEntity := &cart.ShoppingCart{
		ID:        cartModel.ID,
		UserID:    cartModel.UserID,
		ProductID: cartModel.ProductID,
		Quantity:  cartModel.Quantity,
		Price:     cartModel.Price,
		Selected:  cartModel.Selected,
		Version:   cartModel.Version,
		CreatedAt: cartModel.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdatedAt: cartModel.UpdatedAt.Format("2006-01-02 15:04:05"),
	}

	if cartModel.SkuID != nil {
		cartEntity.SkuID = *cartModel.SkuID
	}

	return cartEntity
}
