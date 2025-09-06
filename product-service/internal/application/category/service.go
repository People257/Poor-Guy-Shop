package category

import (
	"context"

	"github.com/people257/poor-guy-shop/product-service/internal/domain/category"
)

// Service 分类应用服务
type Service struct {
	domainService *category.DomainService
	repo          category.Repository
}

// NewService 创建分类应用服务
func NewService(domainService *category.DomainService, repo category.Repository) *Service {
	return &Service{
		domainService: domainService,
		repo:          repo,
	}
}

// CreateCategoryDTO 创建分类DTO
type CreateCategoryDTO struct {
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	Description string  `json:"description"`
	ParentID    *string `json:"parent_id"`
	SortOrder   int     `json:"sort_order"`
	IconURL     string  `json:"icon_url"`
	BannerURL   string  `json:"banner_url"`
}

// UpdateCategoryDTO 更新分类DTO
type UpdateCategoryDTO struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Slug        string  `json:"slug"`
	Description string  `json:"description"`
	ParentID    *string `json:"parent_id"`
	SortOrder   int     `json:"sort_order"`
	IconURL     string  `json:"icon_url"`
	BannerURL   string  `json:"banner_url"`
	IsActive    bool    `json:"is_active"`
}

// ListCategoriesDTO 分类列表查询DTO
type ListCategoriesDTO struct {
	Page      int     `json:"page"`
	PageSize  int     `json:"page_size"`
	ParentID  *string `json:"parent_id"`
	Level     *int    `json:"level"`
	IsActive  *bool   `json:"is_active"`
	Keyword   string  `json:"keyword"`
	SortBy    string  `json:"sort_by"`
	SortOrder string  `json:"sort_order"`
}

// CategoryDTO 分类DTO
type CategoryDTO struct {
	ID          string         `json:"id"`
	Name        string         `json:"name"`
	Slug        string         `json:"slug"`
	Description string         `json:"description"`
	ParentID    *string        `json:"parent_id"`
	Level       int            `json:"level"`
	SortOrder   int            `json:"sort_order"`
	IconURL     string         `json:"icon_url"`
	BannerURL   string         `json:"banner_url"`
	IsActive    bool           `json:"is_active"`
	CreatedAt   string         `json:"created_at"`
	UpdatedAt   string         `json:"updated_at"`
	Children    []*CategoryDTO `json:"children,omitempty"`
}

// CategoryListResult 分类列表结果
type CategoryListResult struct {
	Categories []*CategoryDTO `json:"categories"`
	Total      int64          `json:"total"`
	Page       int            `json:"page"`
	PageSize   int            `json:"page_size"`
}

// CreateCategory 创建分类
func (s *Service) CreateCategory(ctx context.Context, dto *CreateCategoryDTO) (*CategoryDTO, error) {
	cat, err := s.domainService.CreateCategory(
		ctx,
		dto.Name,
		dto.Slug,
		dto.Description,
		dto.ParentID,
		dto.SortOrder,
		dto.IconURL,
		dto.BannerURL,
	)
	if err != nil {
		return nil, err
	}

	return s.toCategoryDTO(cat), nil
}

// UpdateCategory 更新分类
func (s *Service) UpdateCategory(ctx context.Context, dto *UpdateCategoryDTO) (*CategoryDTO, error) {
	cat, err := s.domainService.UpdateCategory(
		ctx,
		dto.ID,
		dto.Name,
		dto.Slug,
		dto.Description,
		dto.ParentID,
		dto.SortOrder,
		dto.IconURL,
		dto.BannerURL,
		dto.IsActive,
	)
	if err != nil {
		return nil, err
	}

	return s.toCategoryDTO(cat), nil
}

// DeleteCategory 删除分类
func (s *Service) DeleteCategory(ctx context.Context, id string) error {
	return s.domainService.DeleteCategory(ctx, id)
}

// GetCategory 获取分类详情
func (s *Service) GetCategory(ctx context.Context, id string) (*CategoryDTO, error) {
	cat, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if cat == nil {
		return nil, category.ErrCategoryNotFound
	}

	return s.toCategoryDTO(cat), nil
}

// ListCategories 获取分类列表
func (s *Service) ListCategories(ctx context.Context, dto *ListCategoriesDTO) (*CategoryListResult, error) {
	params := category.ListParams{
		Page:      dto.Page,
		PageSize:  dto.PageSize,
		ParentID:  dto.ParentID,
		Level:     dto.Level,
		IsActive:  dto.IsActive,
		Keyword:   dto.Keyword,
		SortBy:    dto.SortBy,
		SortOrder: dto.SortOrder,
	}

	categories, total, err := s.repo.List(ctx, params)
	if err != nil {
		return nil, err
	}

	categoryDTOs := make([]*CategoryDTO, len(categories))
	for i, cat := range categories {
		categoryDTOs[i] = s.toCategoryDTO(cat)
	}

	return &CategoryListResult{
		Categories: categoryDTOs,
		Total:      total,
		Page:       dto.Page,
		PageSize:   dto.PageSize,
	}, nil
}

// GetCategoryTree 获取分类树
func (s *Service) GetCategoryTree(ctx context.Context, activeOnly bool) ([]*CategoryDTO, error) {
	categories, err := s.repo.GetTree(ctx, activeOnly)
	if err != nil {
		return nil, err
	}

	return s.toCategoryTreeDTOs(categories), nil
}

// toCategoryDTO 转换为分类DTO
func (s *Service) toCategoryDTO(cat *category.Category) *CategoryDTO {
	return &CategoryDTO{
		ID:          cat.ID,
		Name:        cat.Name,
		Slug:        cat.Slug,
		Description: cat.Description,
		ParentID:    cat.ParentID,
		Level:       cat.Level,
		SortOrder:   cat.SortOrder,
		IconURL:     cat.IconURL,
		BannerURL:   cat.BannerURL,
		IsActive:    cat.IsActive,
		CreatedAt:   cat.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:   cat.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

// toCategoryTreeDTOs 转换为分类树DTOs
func (s *Service) toCategoryTreeDTOs(categories []*category.Category) []*CategoryDTO {
	dtos := make([]*CategoryDTO, len(categories))
	for i, cat := range categories {
		dto := s.toCategoryDTO(cat)
		if len(cat.Children) > 0 {
			dto.Children = s.toCategoryTreeDTOs(cat.Children)
		}
		dtos[i] = dto
	}
	return dtos
}
