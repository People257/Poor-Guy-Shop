package product

import (
	"context"

	"github.com/shopspring/decimal"
)

// Repository 商品仓储接口
type Repository interface {
	// Create 创建商品
	Create(ctx context.Context, product *Product) error

	// Update 更新商品
	Update(ctx context.Context, product *Product) error

	// Delete 删除商品
	Delete(ctx context.Context, id string) error

	// GetByID 根据ID获取商品
	GetByID(ctx context.Context, id string) (*Product, error)

	// GetBySlug 根据slug获取商品
	GetBySlug(ctx context.Context, slug string) (*Product, error)

	// List 获取商品列表
	List(ctx context.Context, params ListParams) ([]*Product, int64, error)

	// Search 搜索商品
	Search(ctx context.Context, params SearchParams) ([]*Product, int64, error)

	// ExistsBySlug 检查slug是否存在
	ExistsBySlug(ctx context.Context, slug string, excludeID ...string) (bool, error)

	// GetByCategoryID 根据分类ID获取商品
	GetByCategoryID(ctx context.Context, categoryID string) ([]*Product, error)

	// GetByBrandID 根据品牌ID获取商品
	GetByBrandID(ctx context.Context, brandID string) ([]*Product, error)
}

// SKURepository SKU仓储接口
type SKURepository interface {
	// Create 创建SKU
	Create(ctx context.Context, sku *ProductSKU) error

	// Update 更新SKU
	Update(ctx context.Context, sku *ProductSKU) error

	// Delete 删除SKU
	Delete(ctx context.Context, id string) error

	// GetByID 根据ID获取SKU
	GetByID(ctx context.Context, id string) (*ProductSKU, error)

	// GetBySKUCode 根据SKU编码获取SKU
	GetBySKUCode(ctx context.Context, skuCode string) (*ProductSKU, error)

	// ListByProductID 根据商品ID获取SKU列表
	ListByProductID(ctx context.Context, productID string, isActive *bool) ([]*ProductSKU, error)

	// ExistsBySKUCode 检查SKU编码是否存在
	ExistsBySKUCode(ctx context.Context, skuCode string, excludeID ...string) (bool, error)

	// DeleteByProductID 根据商品ID删除所有SKU
	DeleteByProductID(ctx context.Context, productID string) error
}

// ListParams 商品列表查询参数
type ListParams struct {
	Page       int
	PageSize   int
	CategoryID *string
	BrandID    *string
	Status     *ProductStatus
	IsFeatured *bool
	Keyword    string
	PriceMin   *decimal.Decimal
	PriceMax   *decimal.Decimal
	SortBy     string // created_at, price, name, sort_order
	SortOrder  string // asc, desc
}

// SearchParams 商品搜索参数
type SearchParams struct {
	Keyword    string
	Page       int
	PageSize   int
	CategoryID *string
	BrandID    *string
	PriceMin   *decimal.Decimal
	PriceMax   *decimal.Decimal
}
