package product

import (
	"context"

	"github.com/shopspring/decimal"

	"github.com/people257/poor-guy-shop/product-service/internal/domain/product"
)

// Service 商品应用服务
type Service struct {
	domainService *product.DomainService
	repo          product.Repository
	skuRepo       product.SKURepository
}

// NewService 创建商品应用服务
func NewService(domainService *product.DomainService, repo product.Repository, skuRepo product.SKURepository) *Service {
	return &Service{
		domainService: domainService,
		repo:          repo,
		skuRepo:       skuRepo,
	}
}

// CreateProductDTO 创建商品DTO
type CreateProductDTO struct {
	Name             string          `json:"name"`
	Slug             string          `json:"slug"`
	Description      string          `json:"description"`
	ShortDescription string          `json:"short_description"`
	CategoryID       string          `json:"category_id"`
	BrandID          *string         `json:"brand_id"`
	MarketPrice      decimal.Decimal `json:"market_price"`
	SalePrice        decimal.Decimal `json:"sale_price"`
	CostPrice        decimal.Decimal `json:"cost_price"`
	MainImageURL     string          `json:"main_image_url"`
	ImageURLs        []string        `json:"image_urls"`
	VideoURL         string          `json:"video_url"`
	Tags             []string        `json:"tags"`
	Specifications   map[string]any  `json:"specifications"`
	IsFeatured       bool            `json:"is_featured"`
	IsVirtual        bool            `json:"is_virtual"`
	SEOTitle         string          `json:"seo_title"`
	SEODescription   string          `json:"seo_description"`
	SEOKeywords      string          `json:"seo_keywords"`
	SortOrder        int             `json:"sort_order"`
}

// UpdateProductDTO 更新商品DTO
type UpdateProductDTO struct {
	ID               string                `json:"id"`
	Name             string                `json:"name"`
	Slug             string                `json:"slug"`
	Description      string                `json:"description"`
	ShortDescription string                `json:"short_description"`
	CategoryID       string                `json:"category_id"`
	BrandID          *string               `json:"brand_id"`
	MarketPrice      decimal.Decimal       `json:"market_price"`
	SalePrice        decimal.Decimal       `json:"sale_price"`
	CostPrice        decimal.Decimal       `json:"cost_price"`
	MainImageURL     string                `json:"main_image_url"`
	ImageURLs        []string              `json:"image_urls"`
	VideoURL         string                `json:"video_url"`
	Tags             []string              `json:"tags"`
	Specifications   map[string]any        `json:"specifications"`
	Status           product.ProductStatus `json:"status"`
	IsFeatured       bool                  `json:"is_featured"`
	IsVirtual        bool                  `json:"is_virtual"`
	SEOTitle         string                `json:"seo_title"`
	SEODescription   string                `json:"seo_description"`
	SEOKeywords      string                `json:"seo_keywords"`
	SortOrder        int                   `json:"sort_order"`
}

// ListProductsDTO 商品列表查询DTO
type ListProductsDTO struct {
	Page       int                    `json:"page"`
	PageSize   int                    `json:"page_size"`
	CategoryID *string                `json:"category_id"`
	BrandID    *string                `json:"brand_id"`
	Status     *product.ProductStatus `json:"status"`
	IsFeatured *bool                  `json:"is_featured"`
	Keyword    string                 `json:"keyword"`
	PriceMin   *decimal.Decimal       `json:"price_min"`
	PriceMax   *decimal.Decimal       `json:"price_max"`
	SortBy     string                 `json:"sort_by"`
	SortOrder  string                 `json:"sort_order"`
}

// SearchProductsDTO 商品搜索DTO
type SearchProductsDTO struct {
	Keyword    string           `json:"keyword"`
	Page       int              `json:"page"`
	PageSize   int              `json:"page_size"`
	CategoryID *string          `json:"category_id"`
	BrandID    *string          `json:"brand_id"`
	PriceMin   *decimal.Decimal `json:"price_min"`
	PriceMax   *decimal.Decimal `json:"price_max"`
}

// ProductDTO 商品DTO
type ProductDTO struct {
	ID               string                `json:"id"`
	Name             string                `json:"name"`
	Slug             string                `json:"slug"`
	Description      string                `json:"description"`
	ShortDescription string                `json:"short_description"`
	CategoryID       string                `json:"category_id"`
	BrandID          *string               `json:"brand_id"`
	MarketPrice      string                `json:"market_price"`
	SalePrice        string                `json:"sale_price"`
	CostPrice        string                `json:"cost_price"`
	MainImageURL     string                `json:"main_image_url"`
	ImageURLs        []string              `json:"image_urls"`
	VideoURL         string                `json:"video_url"`
	Tags             []string              `json:"tags"`
	Specifications   map[string]any        `json:"specifications"`
	Status           product.ProductStatus `json:"status"`
	IsFeatured       bool                  `json:"is_featured"`
	IsVirtual        bool                  `json:"is_virtual"`
	PublishAt        *string               `json:"publish_at"`
	SEOTitle         string                `json:"seo_title"`
	SEODescription   string                `json:"seo_description"`
	SEOKeywords      string                `json:"seo_keywords"`
	SortOrder        int                   `json:"sort_order"`
	CreatedAt        string                `json:"created_at"`
	UpdatedAt        string                `json:"updated_at"`
	SKUs             []*ProductSKUDTO      `json:"skus,omitempty"`
	CategoryName     string                `json:"category_name,omitempty"`
	BrandName        string                `json:"brand_name,omitempty"`
}

// ProductSKUDTO 商品SKU DTO
type ProductSKUDTO struct {
	ID               string            `json:"id"`
	ProductID        string            `json:"product_id"`
	SKUCode          string            `json:"sku_code"`
	Name             string            `json:"name"`
	MarketPrice      string            `json:"market_price"`
	SalePrice        string            `json:"sale_price"`
	CostPrice        string            `json:"cost_price"`
	StockQuantity    int               `json:"stock_quantity"`
	ReservedQuantity int               `json:"reserved_quantity"`
	SoldQuantity     int               `json:"sold_quantity"`
	Weight           string            `json:"weight"`
	Dimensions       map[string]any    `json:"dimensions"`
	ImageURL         string            `json:"image_url"`
	Attributes       map[string]string `json:"attributes"`
	Status           int               `json:"status"`
	CreatedAt        string            `json:"created_at"`
	UpdatedAt        string            `json:"updated_at"`
}

// ProductListResult 商品列表结果
type ProductListResult struct {
	Products []*ProductDTO `json:"products"`
	Total    int64         `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"page_size"`
}

// CreateProduct 创建商品
func (s *Service) CreateProduct(ctx context.Context, dto *CreateProductDTO) (*ProductDTO, error) {
	p, err := s.domainService.CreateProduct(
		ctx,
		dto.Name,
		dto.Slug,
		dto.Description,
		dto.ShortDescription,
		dto.CategoryID,
		dto.BrandID,
		dto.MarketPrice,
		dto.SalePrice,
		dto.CostPrice,
		dto.MainImageURL,
		dto.ImageURLs,
		dto.VideoURL,
		dto.Tags,
		dto.Specifications,
		dto.IsFeatured,
		dto.IsVirtual,
		dto.SEOTitle,
		dto.SEODescription,
		dto.SEOKeywords,
		dto.SortOrder,
	)
	if err != nil {
		return nil, err
	}

	return s.toProductDTO(p), nil
}

// UpdateProduct 更新商品
func (s *Service) UpdateProduct(ctx context.Context, dto *UpdateProductDTO) (*ProductDTO, error) {
	p, err := s.domainService.UpdateProduct(
		ctx,
		dto.ID,
		dto.Name,
		dto.Slug,
		dto.Description,
		dto.ShortDescription,
		dto.CategoryID,
		dto.BrandID,
		dto.MarketPrice,
		dto.SalePrice,
		dto.CostPrice,
		dto.MainImageURL,
		dto.ImageURLs,
		dto.VideoURL,
		dto.Tags,
		dto.Specifications,
		dto.Status,
		dto.IsFeatured,
		dto.IsVirtual,
		dto.SEOTitle,
		dto.SEODescription,
		dto.SEOKeywords,
		dto.SortOrder,
	)
	if err != nil {
		return nil, err
	}

	return s.toProductDTO(p), nil
}

// DeleteProduct 删除商品
func (s *Service) DeleteProduct(ctx context.Context, id string) error {
	return s.domainService.DeleteProduct(ctx, id)
}

// GetProduct 获取商品详情
func (s *Service) GetProduct(ctx context.Context, id string) (*ProductDTO, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, product.ErrProductNotFound
	}

	return s.toProductDTO(p), nil
}

// ListProducts 获取商品列表
func (s *Service) ListProducts(ctx context.Context, dto *ListProductsDTO) (*ProductListResult, error) {
	params := product.ListParams{
		Page:       dto.Page,
		PageSize:   dto.PageSize,
		CategoryID: dto.CategoryID,
		BrandID:    dto.BrandID,
		Status:     dto.Status,
		IsFeatured: dto.IsFeatured,
		Keyword:    dto.Keyword,
		PriceMin:   dto.PriceMin,
		PriceMax:   dto.PriceMax,
		SortBy:     dto.SortBy,
		SortOrder:  dto.SortOrder,
	}

	products, total, err := s.repo.List(ctx, params)
	if err != nil {
		return nil, err
	}

	productDTOs := make([]*ProductDTO, len(products))
	for i, p := range products {
		productDTOs[i] = s.toProductDTO(p)
	}

	return &ProductListResult{
		Products: productDTOs,
		Total:    total,
		Page:     dto.Page,
		PageSize: dto.PageSize,
	}, nil
}

// SearchProducts 搜索商品
func (s *Service) SearchProducts(ctx context.Context, dto *SearchProductsDTO) (*ProductListResult, error) {
	params := product.SearchParams{
		Keyword:    dto.Keyword,
		Page:       dto.Page,
		PageSize:   dto.PageSize,
		CategoryID: dto.CategoryID,
		BrandID:    dto.BrandID,
		PriceMin:   dto.PriceMin,
		PriceMax:   dto.PriceMax,
	}

	products, total, err := s.repo.Search(ctx, params)
	if err != nil {
		return nil, err
	}

	productDTOs := make([]*ProductDTO, len(products))
	for i, p := range products {
		productDTOs[i] = s.toProductDTO(p)
	}

	return &ProductListResult{
		Products: productDTOs,
		Total:    total,
		Page:     dto.Page,
		PageSize: dto.PageSize,
	}, nil
}

// toProductDTO 转换为商品DTO
func (s *Service) toProductDTO(p *product.Product) *ProductDTO {
	dto := &ProductDTO{
		ID:               p.ID,
		Name:             p.Name,
		Slug:             p.Slug,
		Description:      p.Description,
		ShortDescription: p.ShortDescription,
		CategoryID:       p.CategoryID,
		BrandID:          p.BrandID,
		MarketPrice:      p.MarketPrice.String(),
		SalePrice:        p.SalePrice.String(),
		CostPrice:        p.CostPrice.String(),
		MainImageURL:     p.MainImageURL,
		ImageURLs:        p.GetImageURLs(),
		VideoURL:         p.VideoURL,
		Tags:             p.GetTags(),
		Specifications:   p.GetSpecifications(),
		Status:           p.Status,
		IsFeatured:       p.IsFeatured,
		IsVirtual:        p.IsVirtual,
		SEOTitle:         p.SEOTitle,
		SEODescription:   p.SEODescription,
		SEOKeywords:      p.SEOKeywords,
		SortOrder:        p.SortOrder,
		CreatedAt:        p.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:        p.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}

	if p.PublishAt != nil {
		publishAt := p.PublishAt.Format("2006-01-02T15:04:05Z07:00")
		dto.PublishAt = &publishAt
	}

	// 转换SKUs
	if len(p.SKUs) > 0 {
		dto.SKUs = make([]*ProductSKUDTO, len(p.SKUs))
		for i, sku := range p.SKUs {
			dto.SKUs[i] = s.toProductSKUDTO(sku)
		}
	}

	return dto
}

// toProductSKUDTO 转换为SKU DTO
func (s *Service) toProductSKUDTO(sku *product.ProductSKU) *ProductSKUDTO {
	return &ProductSKUDTO{
		ID:               sku.ID,
		ProductID:        sku.ProductID,
		SKUCode:          sku.SKUCode,
		Name:             sku.Name,
		MarketPrice:      sku.MarketPrice.String(),
		SalePrice:        sku.SalePrice.String(),
		CostPrice:        sku.CostPrice.String(),
		StockQuantity:    sku.StockQuantity,
		ReservedQuantity: sku.ReservedQuantity,
		SoldQuantity:     sku.SoldQuantity,
		Weight:           sku.Weight.String(),
		Dimensions:       sku.GetDimensions(),
		ImageURL:         sku.ImageURL,
		Attributes:       sku.GetAttributes(),
		Status:           int(sku.Status),
		CreatedAt:        sku.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
		UpdatedAt:        sku.UpdatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}

