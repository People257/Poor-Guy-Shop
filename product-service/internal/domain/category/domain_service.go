package category

import (
	"context"
	"fmt"
)

// DomainService 分类领域服务
type DomainService struct {
	repo Repository
}

// NewDomainService 创建分类领域服务
func NewDomainService(repo Repository) *DomainService {
	return &DomainService{
		repo: repo,
	}
}

// CreateCategory 创建分类
func (s *DomainService) CreateCategory(ctx context.Context, name, slug, description string, parentID *string, sortOrder int, iconURL, bannerURL string) (*Category, error) {
	// 检查slug是否已存在
	exists, err := s.repo.ExistsBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("检查slug是否存在失败: %w", err)
	}
	if exists {
		return nil, ErrCategorySlugExists
	}

	// 创建分类实体
	category, err := NewCategory(name, slug, description, parentID, sortOrder, iconURL, bannerURL)
	if err != nil {
		return nil, err
	}

	// 处理父分类逻辑
	if parentID != nil && *parentID != "" {
		parent, err := s.repo.GetByID(ctx, *parentID)
		if err != nil {
			return nil, fmt.Errorf("获取父分类失败: %w", err)
		}
		if parent == nil {
			return nil, ErrParentCategoryNotFound
		}
		if !parent.IsActive {
			return nil, ErrParentCategoryInactive
		}

		// 计算分类级别
		level := parent.Level + 1
		if level > 5 { // 最多支持5级分类
			return nil, ErrCategoryLevelTooDeep
		}
		category.SetLevel(level)
	}

	// 保存分类
	if err := s.repo.Create(ctx, category); err != nil {
		return nil, fmt.Errorf("创建分类失败: %w", err)
	}

	return category, nil
}

// UpdateCategory 更新分类
func (s *DomainService) UpdateCategory(ctx context.Context, id string, name, slug, description string, parentID *string, sortOrder int, iconURL, bannerURL string, isActive bool) (*Category, error) {
	// 获取现有分类
	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取分类失败: %w", err)
	}
	if category == nil {
		return nil, ErrCategoryNotFound
	}

	// 检查slug是否已存在（排除当前分类）
	if slug != category.Slug {
		exists, err := s.repo.ExistsBySlug(ctx, slug, id)
		if err != nil {
			return nil, fmt.Errorf("检查slug是否存在失败: %w", err)
		}
		if exists {
			return nil, ErrCategorySlugExists
		}
	}

	// 处理父分类变更
	if parentID != nil && *parentID != "" {
		// 检查是否设置自己为父分类
		if *parentID == id {
			return nil, ErrCategoryCircularRef
		}

		// 检查是否设置子分类为父分类
		if err := s.checkCircularReference(ctx, id, *parentID); err != nil {
			return nil, err
		}

		// 获取新父分类
		parent, err := s.repo.GetByID(ctx, *parentID)
		if err != nil {
			return nil, fmt.Errorf("获取父分类失败: %w", err)
		}
		if parent == nil {
			return nil, ErrParentCategoryNotFound
		}
		if !parent.IsActive {
			return nil, ErrParentCategoryInactive
		}

		// 计算新级别
		newLevel := parent.Level + 1
		if newLevel > 5 {
			return nil, ErrCategoryLevelTooDeep
		}

		// 如果父分类发生变化，需要更新整个子树的级别
		if category.ParentID == nil || *category.ParentID != *parentID {
			if err := s.updateCategoryTree(ctx, id, newLevel); err != nil {
				return nil, fmt.Errorf("更新分类树级别失败: %w", err)
			}
		}
	} else {
		// 设置为根分类
		if category.ParentID != nil {
			if err := s.updateCategoryTree(ctx, id, 1); err != nil {
				return nil, fmt.Errorf("更新分类树级别失败: %w", err)
			}
		}
	}

	// 更新分类信息
	if err := category.Update(name, slug, description, parentID, sortOrder, iconURL, bannerURL, isActive); err != nil {
		return nil, err
	}

	// 保存更新
	if err := s.repo.Update(ctx, category); err != nil {
		return nil, fmt.Errorf("更新分类失败: %w", err)
	}

	return category, nil
}

// DeleteCategory 删除分类
func (s *DomainService) DeleteCategory(ctx context.Context, id string) error {
	// 获取分类
	category, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("获取分类失败: %w", err)
	}
	if category == nil {
		return ErrCategoryNotFound
	}

	// 检查是否有子分类
	childCount, err := s.repo.CountChildren(ctx, id)
	if err != nil {
		return fmt.Errorf("统计子分类失败: %w", err)
	}
	if childCount > 0 {
		return ErrCategoryHasChildren
	}

	// 检查是否有商品
	productCount, err := s.repo.CountProducts(ctx, id)
	if err != nil {
		return fmt.Errorf("统计分类商品失败: %w", err)
	}
	if productCount > 0 {
		return ErrCategoryHasProducts
	}

	// 删除分类
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("删除分类失败: %w", err)
	}

	return nil
}

// checkCircularReference 检查循环引用
func (s *DomainService) checkCircularReference(ctx context.Context, categoryID, parentID string) error {
	// 获取父分类的所有祖先
	current := parentID
	visited := make(map[string]bool)

	for current != "" {
		if visited[current] {
			return ErrCategoryCircularRef
		}
		if current == categoryID {
			return ErrCategoryCircularRef
		}

		visited[current] = true

		category, err := s.repo.GetByID(ctx, current)
		if err != nil {
			return fmt.Errorf("获取分类失败: %w", err)
		}
		if category == nil || category.ParentID == nil {
			break
		}

		current = *category.ParentID
	}

	return nil
}

// updateCategoryTree 更新分类树的级别
func (s *DomainService) updateCategoryTree(ctx context.Context, categoryID string, newLevel int) error {
	var updates []LevelUpdate

	// 收集需要更新的分类
	if err := s.collectLevelUpdates(ctx, categoryID, newLevel, &updates); err != nil {
		return err
	}

	// 批量更新
	if len(updates) > 0 {
		if err := s.repo.BatchUpdateLevel(ctx, updates); err != nil {
			return err
		}
	}

	return nil
}

// collectLevelUpdates 收集需要更新级别的分类
func (s *DomainService) collectLevelUpdates(ctx context.Context, categoryID string, level int, updates *[]LevelUpdate) error {
	*updates = append(*updates, LevelUpdate{
		ID:    categoryID,
		Level: level,
	})

	// 获取子分类
	children, err := s.repo.GetByParentID(ctx, categoryID)
	if err != nil {
		return err
	}

	// 递归处理子分类
	for _, child := range children {
		if err := s.collectLevelUpdates(ctx, child.ID, level+1, updates); err != nil {
			return err
		}
	}

	return nil
}
