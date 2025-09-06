package repository

import (
	"context"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"github.com/people257/poor-guy-shop/product-service/internal/domain/brand"
)

// BrandRepository 品牌仓储实现
type BrandRepository struct {
	db *gorm.DB
}

// NewBrandRepository 创建品牌仓储
func NewBrandRepository(db *gorm.DB) brand.Repository {
	return &BrandRepository{
		db: db,
	}
}

// Create 创建品牌
func (r *BrandRepository) Create(ctx context.Context, b *brand.Brand) error {
	return r.db.WithContext(ctx).Create(b).Error
}

// Update 更新品牌
func (r *BrandRepository) Update(ctx context.Context, b *brand.Brand) error {
	return r.db.WithContext(ctx).Save(b).Error
}

// Delete 删除品牌
func (r *BrandRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&brand.Brand{}, "id = ?", id).Error
}

// GetByID 根据ID获取品牌
func (r *BrandRepository) GetByID(ctx context.Context, id string) (*brand.Brand, error) {
	var b brand.Brand
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&b).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &b, nil
}

// GetBySlug 根据slug获取品牌
func (r *BrandRepository) GetBySlug(ctx context.Context, slug string) (*brand.Brand, error) {
	var b brand.Brand
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&b).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &b, nil
}

// List 获取品牌列表
func (r *BrandRepository) List(ctx context.Context, params brand.ListParams) ([]*brand.Brand, int64, error) {
	db := r.db.WithContext(ctx)

	// 构建查询条件
	if params.IsActive != nil {
		db = db.Where("is_active = ?", *params.IsActive)
	}

	if params.Keyword != "" {
		keyword := "%" + strings.ToLower(params.Keyword) + "%"
		db = db.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?", keyword, keyword)
	}

	// 统计总数
	var total int64
	if err := db.Model(&brand.Brand{}).Count(&total).Error; err != nil {
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

	var brands []*brand.Brand
	err := db.Find(&brands).Error
	if err != nil {
		return nil, 0, err
	}

	return brands, total, nil
}

// CountProducts 统计品牌下商品数量
func (r *BrandRepository) CountProducts(ctx context.Context, id string) (int64, error) {
	var count int64
	// 这里需要查询products表，暂时返回0
	// TODO: 实现商品计数逻辑
	return count, nil
}

// ExistsBySlug 检查slug是否存在
func (r *BrandRepository) ExistsBySlug(ctx context.Context, slug string, excludeID ...string) (bool, error) {
	db := r.db.WithContext(ctx).Model(&brand.Brand{}).Where("slug = ?", slug)

	if len(excludeID) > 0 && excludeID[0] != "" {
		db = db.Where("id != ?", excludeID[0])
	}

	var count int64
	err := db.Count(&count).Error
	return count > 0, err
}
