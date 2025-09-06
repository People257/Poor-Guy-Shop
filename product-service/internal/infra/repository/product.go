package repository

import (
	"context"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"github.com/people257/poor-guy-shop/product-service/internal/domain/product"
)

// ProductRepository 商品仓储实现
type ProductRepository struct {
	db *gorm.DB
}

// NewProductRepository 创建商品仓储
func NewProductRepository(db *gorm.DB) product.Repository {
	return &ProductRepository{
		db: db,
	}
}

// Create 创建商品
func (r *ProductRepository) Create(ctx context.Context, p *product.Product) error {
	return r.db.WithContext(ctx).Create(p).Error
}

// Update 更新商品
func (r *ProductRepository) Update(ctx context.Context, p *product.Product) error {
	return r.db.WithContext(ctx).Save(p).Error
}

// Delete 删除商品
func (r *ProductRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&product.Product{}, "id = ?", id).Error
}

// GetByID 根据ID获取商品
func (r *ProductRepository) GetByID(ctx context.Context, id string) (*product.Product, error) {
	var p product.Product
	err := r.db.WithContext(ctx).Preload("SKUs").Where("id = ?", id).First(&p).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

// GetBySlug 根据slug获取商品
func (r *ProductRepository) GetBySlug(ctx context.Context, slug string) (*product.Product, error) {
	var p product.Product
	err := r.db.WithContext(ctx).Preload("SKUs").Where("slug = ?", slug).First(&p).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &p, nil
}

// List 获取商品列表
func (r *ProductRepository) List(ctx context.Context, params product.ListParams) ([]*product.Product, int64, error) {
	db := r.db.WithContext(ctx)

	// 构建查询条件
	if params.CategoryID != nil {
		db = db.Where("category_id = ?", *params.CategoryID)
	}

	if params.BrandID != nil {
		db = db.Where("brand_id = ?", *params.BrandID)
	}

	if params.Status != nil {
		db = db.Where("status = ?", *params.Status)
	}

	if params.IsFeatured != nil {
		db = db.Where("is_featured = ?", *params.IsFeatured)
	}

	if params.Keyword != "" {
		keyword := "%" + strings.ToLower(params.Keyword) + "%"
		db = db.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ? OR LOWER(short_description) LIKE ?", keyword, keyword, keyword)
	}

	if params.PriceMin != nil {
		db = db.Where("sale_price >= ?", *params.PriceMin)
	}

	if params.PriceMax != nil {
		db = db.Where("sale_price <= ?", *params.PriceMax)
	}

	// 统计总数
	var total int64
	if err := db.Model(&product.Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 排序
	orderBy := "sort_order ASC, created_at DESC"
	if params.SortBy != "" {
		order := "ASC"
		if params.SortOrder == "desc" {
			order = "DESC"
		}
		switch params.SortBy {
		case "name":
			orderBy = fmt.Sprintf("name %s", order)
		case "price":
			orderBy = fmt.Sprintf("sale_price %s", order)
		case "sort_order":
			orderBy = fmt.Sprintf("sort_order %s", order)
		case "created_at":
			orderBy = fmt.Sprintf("created_at %s", order)
		}
	}
	db = db.Order(orderBy)

	// 分页
	if params.Page > 0 && params.PageSize > 0 {
		offset := (params.Page - 1) * params.PageSize
		db = db.Offset(offset).Limit(params.PageSize)
	}

	var products []*product.Product
	err := db.Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// Search 搜索商品
func (r *ProductRepository) Search(ctx context.Context, params product.SearchParams) ([]*product.Product, int64, error) {
	db := r.db.WithContext(ctx)

	// 构建搜索条件
	if params.Keyword != "" {
		keyword := "%" + strings.ToLower(params.Keyword) + "%"
		db = db.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ? OR LOWER(short_description) LIKE ?", keyword, keyword, keyword)
	}

	if params.CategoryID != nil {
		db = db.Where("category_id = ?", *params.CategoryID)
	}

	if params.BrandID != nil {
		db = db.Where("brand_id = ?", *params.BrandID)
	}

	if params.PriceMin != nil {
		db = db.Where("sale_price >= ?", *params.PriceMin)
	}

	if params.PriceMax != nil {
		db = db.Where("sale_price <= ?", *params.PriceMax)
	}

	// 只搜索上架商品
	db = db.Where("status = ?", product.ProductStatusActive)

	// 统计总数
	var total int64
	if err := db.Model(&product.Product{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 排序（搜索结果按相关性排序，这里简化为按创建时间倒序）
	db = db.Order("created_at DESC")

	// 分页
	if params.Page > 0 && params.PageSize > 0 {
		offset := (params.Page - 1) * params.PageSize
		db = db.Offset(offset).Limit(params.PageSize)
	}

	var products []*product.Product
	err := db.Find(&products).Error
	if err != nil {
		return nil, 0, err
	}

	return products, total, nil
}

// ExistsBySlug 检查slug是否存在
func (r *ProductRepository) ExistsBySlug(ctx context.Context, slug string, excludeID ...string) (bool, error) {
	db := r.db.WithContext(ctx).Model(&product.Product{}).Where("slug = ?", slug)

	if len(excludeID) > 0 && excludeID[0] != "" {
		db = db.Where("id != ?", excludeID[0])
	}

	var count int64
	err := db.Count(&count).Error
	return count > 0, err
}

// GetByCategoryID 根据分类ID获取商品
func (r *ProductRepository) GetByCategoryID(ctx context.Context, categoryID string) ([]*product.Product, error) {
	var products []*product.Product
	err := r.db.WithContext(ctx).Where("category_id = ?", categoryID).Find(&products).Error
	return products, err
}

// GetByBrandID 根据品牌ID获取商品
func (r *ProductRepository) GetByBrandID(ctx context.Context, brandID string) ([]*product.Product, error) {
	var products []*product.Product
	err := r.db.WithContext(ctx).Where("brand_id = ?", brandID).Find(&products).Error
	return products, err
}

// ProductSKURepository SKU仓储实现
type ProductSKURepository struct {
	db *gorm.DB
}

// NewProductSKURepository 创建SKU仓储
func NewProductSKURepository(db *gorm.DB) product.SKURepository {
	return &ProductSKURepository{
		db: db,
	}
}

// Create 创建SKU
func (r *ProductSKURepository) Create(ctx context.Context, sku *product.ProductSKU) error {
	return r.db.WithContext(ctx).Create(sku).Error
}

// Update 更新SKU
func (r *ProductSKURepository) Update(ctx context.Context, sku *product.ProductSKU) error {
	return r.db.WithContext(ctx).Save(sku).Error
}

// Delete 删除SKU
func (r *ProductSKURepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&product.ProductSKU{}, "id = ?", id).Error
}

// GetByID 根据ID获取SKU
func (r *ProductSKURepository) GetByID(ctx context.Context, id string) (*product.ProductSKU, error) {
	var sku product.ProductSKU
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&sku).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &sku, nil
}

// GetBySKUCode 根据SKU编码获取SKU
func (r *ProductSKURepository) GetBySKUCode(ctx context.Context, skuCode string) (*product.ProductSKU, error) {
	var sku product.ProductSKU
	err := r.db.WithContext(ctx).Where("sku_code = ?", skuCode).First(&sku).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &sku, nil
}

// ListByProductID 根据商品ID获取SKU列表
func (r *ProductSKURepository) ListByProductID(ctx context.Context, productID string, isActive *bool) ([]*product.ProductSKU, error) {
	db := r.db.WithContext(ctx).Where("product_id = ?", productID)

	if isActive != nil {
		db = db.Where("is_active = ?", *isActive)
	}

	var skus []*product.ProductSKU
	err := db.Order("sort_order ASC, created_at ASC").Find(&skus).Error
	return skus, err
}

// ExistsBySKUCode 检查SKU编码是否存在
func (r *ProductSKURepository) ExistsBySKUCode(ctx context.Context, skuCode string, excludeID ...string) (bool, error) {
	db := r.db.WithContext(ctx).Model(&product.ProductSKU{}).Where("sku_code = ?", skuCode)

	if len(excludeID) > 0 && excludeID[0] != "" {
		db = db.Where("id != ?", excludeID[0])
	}

	var count int64
	err := db.Count(&count).Error
	return count > 0, err
}

// DeleteByProductID 根据商品ID删除所有SKU
func (r *ProductSKURepository) DeleteByProductID(ctx context.Context, productID string) error {
	return r.db.WithContext(ctx).Where("product_id = ?", productID).Delete(&product.ProductSKU{}).Error
}
