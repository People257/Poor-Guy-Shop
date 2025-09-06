package brand

import (
	"context"
	"fmt"
)

// DomainService 品牌领域服务
type DomainService struct {
	repo Repository
}

// NewDomainService 创建品牌领域服务
func NewDomainService(repo Repository) *DomainService {
	return &DomainService{
		repo: repo,
	}
}

// CreateBrand 创建品牌
func (s *DomainService) CreateBrand(ctx context.Context, name, slug, description, logoURL, websiteURL string, sortOrder int) (*Brand, error) {
	// 检查slug是否已存在
	exists, err := s.repo.ExistsBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("检查slug是否存在失败: %w", err)
	}
	if exists {
		return nil, ErrBrandSlugExists
	}

	// 创建品牌实体
	brand, err := NewBrand(name, slug, description, logoURL, websiteURL, sortOrder)
	if err != nil {
		return nil, err
	}

	// 保存品牌
	if err := s.repo.Create(ctx, brand); err != nil {
		return nil, fmt.Errorf("创建品牌失败: %w", err)
	}

	return brand, nil
}

// UpdateBrand 更新品牌
func (s *DomainService) UpdateBrand(ctx context.Context, id, name, slug, description, logoURL, websiteURL string, sortOrder int, isActive bool) (*Brand, error) {
	// 获取现有品牌
	brand, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取品牌失败: %w", err)
	}
	if brand == nil {
		return nil, ErrBrandNotFound
	}

	// 检查slug是否已存在（排除当前品牌）
	if slug != brand.Slug {
		exists, err := s.repo.ExistsBySlug(ctx, slug, id)
		if err != nil {
			return nil, fmt.Errorf("检查slug是否存在失败: %w", err)
		}
		if exists {
			return nil, ErrBrandSlugExists
		}
	}

	// 更新品牌信息
	if err := brand.Update(name, slug, description, logoURL, websiteURL, sortOrder, isActive); err != nil {
		return nil, err
	}

	// 保存更新
	if err := s.repo.Update(ctx, brand); err != nil {
		return nil, fmt.Errorf("更新品牌失败: %w", err)
	}

	return brand, nil
}

// DeleteBrand 删除品牌
func (s *DomainService) DeleteBrand(ctx context.Context, id string) error {
	// 获取品牌
	brand, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("获取品牌失败: %w", err)
	}
	if brand == nil {
		return ErrBrandNotFound
	}

	// 检查是否有商品
	productCount, err := s.repo.CountProducts(ctx, id)
	if err != nil {
		return fmt.Errorf("统计品牌商品失败: %w", err)
	}
	if productCount > 0 {
		return ErrBrandHasProducts
	}

	// 删除品牌
	if err := s.repo.Delete(ctx, id); err != nil {
		return fmt.Errorf("删除品牌失败: %w", err)
	}

	return nil
}
