package repository

import (
	"context"
	"fmt"
	"strings"

	"gorm.io/gorm"

	"github.com/people257/poor-guy-shop/product-service/internal/domain/category"
)

// CategoryRepository 分类仓储实现
type CategoryRepository struct {
	db *gorm.DB
}

// NewCategoryRepository 创建分类仓储
func NewCategoryRepository(db *gorm.DB) category.Repository {
	return &CategoryRepository{
		db: db,
	}
}

// Create 创建分类
func (r *CategoryRepository) Create(ctx context.Context, cat *category.Category) error {
	return r.db.WithContext(ctx).Create(cat).Error
}

// Update 更新分类
func (r *CategoryRepository) Update(ctx context.Context, cat *category.Category) error {
	return r.db.WithContext(ctx).Save(cat).Error
}

// Delete 删除分类
func (r *CategoryRepository) Delete(ctx context.Context, id string) error {
	return r.db.WithContext(ctx).Delete(&category.Category{}, "id = ?", id).Error
}

// GetByID 根据ID获取分类
func (r *CategoryRepository) GetByID(ctx context.Context, id string) (*category.Category, error) {
	var cat category.Category
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&cat).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &cat, nil
}

// GetBySlug 根据slug获取分类
func (r *CategoryRepository) GetBySlug(ctx context.Context, slug string) (*category.Category, error) {
	var cat category.Category
	err := r.db.WithContext(ctx).Where("slug = ?", slug).First(&cat).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, err
	}
	return &cat, nil
}

// List 获取分类列表
func (r *CategoryRepository) List(ctx context.Context, params category.ListParams) ([]*category.Category, int64, error) {
	db := r.db.WithContext(ctx)

	// 构建查询条件
	if params.ParentID != nil {
		if *params.ParentID == "" {
			db = db.Where("parent_id IS NULL")
		} else {
			db = db.Where("parent_id = ?", *params.ParentID)
		}
	}

	if params.Level != nil {
		db = db.Where("level = ?", *params.Level)
	}

	if params.IsActive != nil {
		db = db.Where("is_active = ?", *params.IsActive)
	}

	if params.Keyword != "" {
		keyword := "%" + strings.ToLower(params.Keyword) + "%"
		db = db.Where("LOWER(name) LIKE ? OR LOWER(description) LIKE ?", keyword, keyword)
	}

	// 统计总数
	var total int64
	if err := db.Model(&category.Category{}).Count(&total).Error; err != nil {
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

	var categories []*category.Category
	err := db.Find(&categories).Error
	if err != nil {
		return nil, 0, err
	}

	return categories, total, nil
}

// GetTree 获取分类树
func (r *CategoryRepository) GetTree(ctx context.Context, activeOnly bool) ([]*category.Category, error) {
	db := r.db.WithContext(ctx)

	if activeOnly {
		db = db.Where("is_active = ?", true)
	}

	var categories []*category.Category
	err := db.Order("level ASC, sort_order ASC, created_at ASC").Find(&categories).Error
	if err != nil {
		return nil, err
	}

	// 构建树形结构
	return r.buildTree(categories, nil), nil
}

// GetByParentID 根据父分类ID获取子分类
func (r *CategoryRepository) GetByParentID(ctx context.Context, parentID string) ([]*category.Category, error) {
	var categories []*category.Category
	err := r.db.WithContext(ctx).Where("parent_id = ?", parentID).Order("sort_order ASC, created_at ASC").Find(&categories).Error
	return categories, err
}

// CountChildren 统计子分类数量
func (r *CategoryRepository) CountChildren(ctx context.Context, id string) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).Model(&category.Category{}).Where("parent_id = ?", id).Count(&count).Error
	return count, err
}

// CountProducts 统计分类下商品数量
func (r *CategoryRepository) CountProducts(ctx context.Context, id string) (int64, error) {
	var count int64
	// 这里需要查询products表，暂时返回0
	// TODO: 实现商品计数逻辑
	return count, nil
}

// ExistsBySlug 检查slug是否存在
func (r *CategoryRepository) ExistsBySlug(ctx context.Context, slug string, excludeID ...string) (bool, error) {
	db := r.db.WithContext(ctx).Model(&category.Category{}).Where("slug = ?", slug)

	if len(excludeID) > 0 && excludeID[0] != "" {
		db = db.Where("id != ?", excludeID[0])
	}

	var count int64
	err := db.Count(&count).Error
	return count > 0, err
}

// GetMaxLevel 获取最大层级
func (r *CategoryRepository) GetMaxLevel(ctx context.Context) (int, error) {
	var maxLevel int
	err := r.db.WithContext(ctx).Model(&category.Category{}).Select("COALESCE(MAX(level), 0)").Scan(&maxLevel).Error
	return maxLevel, err
}

// UpdateLevel 更新分类层级
func (r *CategoryRepository) UpdateLevel(ctx context.Context, id string, level int) error {
	return r.db.WithContext(ctx).Model(&category.Category{}).Where("id = ?", id).Update("level", level).Error
}

// BatchUpdateLevel 批量更新分类层级
func (r *CategoryRepository) BatchUpdateLevel(ctx context.Context, updates []category.LevelUpdate) error {
	return r.db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		for _, update := range updates {
			if err := tx.Model(&category.Category{}).Where("id = ?", update.ID).Update("level", update.Level).Error; err != nil {
				return err
			}
		}
		return nil
	})
}

// buildTree 构建树形结构
func (r *CategoryRepository) buildTree(categories []*category.Category, parentID *string) []*category.Category {
	var result []*category.Category

	for _, cat := range categories {
		// 检查是否为当前层级的节点
		if (parentID == nil && cat.ParentID == nil) ||
			(parentID != nil && cat.ParentID != nil && *cat.ParentID == *parentID) {

			// 递归构建子节点
			cat.Children = r.buildTree(categories, &cat.ID)
			result = append(result, cat)
		}
	}

	return result
}
