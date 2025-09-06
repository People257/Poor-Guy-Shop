package brand

import (
	"context"

	"github.com/people257/poor-guy-shop/product-service/internal/domain/brand"
)

// Service 品牌应用服务
type Service struct {
	domainService *brand.DomainService
	repo          brand.Repository
}

// NewService 创建品牌应用服务
func NewService(domainService *brand.DomainService, repo brand.Repository) *Service {
	return &Service{
		domainService: domainService,
		repo:          repo,
	}
}

// CreateBrandDTO 创建品牌DTO
type CreateBrandDTO struct {
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	LogoURL     string `json:"logo_url"`
	WebsiteURL  string `json:"website_url"`
	SortOrder   int    `json:"sort_order"`
}

// UpdateBrandDTO 更新品牌DTO
type UpdateBrandDTO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	LogoURL     string `json:"logo_url"`
	WebsiteURL  string `json:"website_url"`
	SortOrder   int    `json:"sort_order"`
	IsActive    bool   `json:"is_active"`
}

// ListBrandsDTO 品牌列表查询DTO
type ListBrandsDTO struct {
	Page      int    `json:"page"`
	PageSize  int    `json:"page_size"`
	IsActive  *bool  `json:"is_active"`
	Keyword   string `json:"keyword"`
	SortBy    string `json:"sort_by"`
	SortOrder string `json:"sort_order"`
}

// BrandDTO 品牌DTO
type BrandDTO struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Slug        string `json:"slug"`
	Description string `json:"description"`
	LogoURL     string `json:"logo_url"`
	WebsiteURL  string `json:"website_url"`
	SortOrder   int    `json:"sort_order"`
	IsActive    bool   `json:"is_active"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// BrandListResult 品牌列表结果
type BrandListResult struct {
	Brands   []*BrandDTO `json:"brands"`
	Total    int64       `json:"total"`
	Page     int         `json:"page"`
	PageSize int         `json:"page_size"`
}

// CreateBrand 创建品牌
func (s *Service) CreateBrand(ctx context.Context, dto *CreateBrandDTO) (*BrandDTO, error) {
	b, err := s.domainService.CreateBrand(
		ctx,
		dto.Name,
		dto.Slug,
		dto.Description,
		dto.LogoURL,
		dto.WebsiteURL,
		dto.SortOrder,
	)
	if err != nil {
		return nil, err
	}

	return s.toBrandDTO(b), nil
}

// UpdateBrand 更新品牌
func (s *Service) UpdateBrand(ctx context.Context, dto *UpdateBrandDTO) (*BrandDTO, error) {
	b, err := s.domainService.UpdateBrand(
		ctx,
		dto.ID,
		dto.Name,
		dto.Slug,
		dto.Description,
		dto.LogoURL,
		dto.WebsiteURL,
		dto.SortOrder,
		dto.IsActive,
	)
	if err != nil {
		return nil, err
	}

	return s.toBrandDTO(b), nil
}

// DeleteBrand 删除品牌
func (s *Service) DeleteBrand(ctx context.Context, id string) error {
	return s.domainService.DeleteBrand(ctx, id)
}

// GetBrand 获取品牌详情
func (s *Service) GetBrand(ctx context.Context, id string) (*BrandDTO, error) {
	b, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if b == nil {
		return nil, brand.ErrBrandNotFound
	}

	return s.toBrandDTO(b), nil
}

// ListBrands 获取品牌列表
func (s *Service) ListBrands(ctx context.Context, dto *ListBrandsDTO) (*BrandListResult, error) {
	params := brand.ListParams{
		Page:      dto.Page,
		PageSize:  dto.PageSize,
		IsActive:  dto.IsActive,
		Keyword:   dto.Keyword,
		SortBy:    dto.SortBy,
		SortOrder: dto.SortOrder,
	}

	brands, total, err := s.repo.List(ctx, params)
	if err != nil {
		return nil, err
	}

	brandDTOs := make([]*BrandDTO, len(brands))
	for i, b := range brands {
		brandDTOs[i] = s.toBrandDTO(b)
	}

	return &BrandListResult{
		Brands:   brandDTOs,
		Total:    total,
		Page:     dto.Page,
		PageSize: dto.PageSize,
	}, nil
}

// toBrandDTO 转换为品牌DTO
func (s *Service) toBrandDTO(b *brand.Brand) *BrandDTO {
	return &BrandDTO{
		ID:          b.ID,
		Name:        b.Name,
		Slug:        b.Slug,
		Description: b.Description,
		LogoURL:     b.LogoURL,
		WebsiteURL:  b.WebsiteURL,
		SortOrder:   b.SortOrder,
		IsActive:    b.IsActive,
		CreatedAt:   b.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   b.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
