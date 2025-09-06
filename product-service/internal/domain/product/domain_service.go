package product

import (
	"context"
	"fmt"

	"github.com/shopspring/decimal"
)

// DomainService 商品领域服务
type DomainService struct {
	repo    Repository
	skuRepo SKURepository
}

// NewDomainService 创建商品领域服务
func NewDomainService(repo Repository, skuRepo SKURepository) *DomainService {
	return &DomainService{
		repo:    repo,
		skuRepo: skuRepo,
	}
}

// CreateProduct 创建商品
func (s *DomainService) CreateProduct(
	ctx context.Context,
	name, slug, description, shortDescription, categoryID string,
	brandID *string,
	marketPrice, salePrice, costPrice decimal.Decimal,
	mainImageURL string,
	imageURLs []string,
	videoURL string,
	tags []string,
	specifications map[string]any,
	isFeatured, isVirtual bool,
	seoTitle, seoDescription, seoKeywords string,
	sortOrder int,
) (*Product, error) {
	// 检查slug是否已存在
	exists, err := s.repo.ExistsBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("检查slug是否存在失败: %w", err)
	}
	if exists {
		return nil, ErrProductSlugExists
	}

	// 创建商品实体
	product, err := NewProduct(
		name, slug, description, shortDescription, categoryID, brandID,
		marketPrice, salePrice, costPrice,
		mainImageURL, imageURLs, videoURL,
		tags, specifications,
		isFeatured, isVirtual,
		seoTitle, seoDescription, seoKeywords,
		sortOrder,
	)
	if err != nil {
		return nil, err
	}

	// 保存商品
	if err := s.repo.Create(ctx, product); err != nil {
		return nil, fmt.Errorf("创建商品失败: %w", err)
	}

	return product, nil
}

// UpdateProduct 更新商品
func (s *DomainService) UpdateProduct(
	ctx context.Context,
	id, name, slug, description, shortDescription, categoryID string,
	brandID *string,
	marketPrice, salePrice, costPrice decimal.Decimal,
	mainImageURL string,
	imageURLs []string,
	videoURL string,
	tags []string,
	specifications map[string]any,
	status ProductStatus,
	isFeatured, isVirtual bool,
	seoTitle, seoDescription, seoKeywords string,
	sortOrder int,
) (*Product, error) {
	// 获取现有商品
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取商品失败: %w", err)
	}
	if product == nil {
		return nil, ErrProductNotFound
	}

	// 检查slug是否已存在（排除当前商品）
	if slug != product.Slug {
		exists, err := s.repo.ExistsBySlug(ctx, slug, id)
		if err != nil {
			return nil, fmt.Errorf("检查slug是否存在失败: %w", err)
		}
		if exists {
			return nil, ErrProductSlugExists
		}
	}

	// 更新商品信息
	if err := product.Update(
		name, slug, description, shortDescription, categoryID, brandID,
		marketPrice, salePrice, costPrice,
		mainImageURL, imageURLs, videoURL,
		tags, specifications,
		status, isFeatured, isVirtual,
		seoTitle, seoDescription, seoKeywords,
		sortOrder,
	); err != nil {
		return nil, err
	}

	// 保存更新
	if err := s.repo.Update(ctx, product); err != nil {
		return nil, fmt.Errorf("更新商品失败: %w", err)
	}

	return product, nil
}

// DeleteProduct 删除商品
func (s *DomainService) DeleteProduct(ctx context.Context, id string) error {
	// 获取商品
	product, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("获取商品失败: %w", err)
	}
	if product == nil {
		return ErrProductNotFound
	}

	// 软删除商品
	product.Delete()

	// 保存更新
	if err := s.repo.Update(ctx, product); err != nil {
		return fmt.Errorf("删除商品失败: %w", err)
	}

	return nil
}

// CreateProductSKU 创建商品SKU
func (s *DomainService) CreateProductSKU(
	ctx context.Context,
	productID, skuCode, name string,
	marketPrice, salePrice, costPrice decimal.Decimal,
	weight decimal.Decimal,
	dimensions map[string]any,
	imageURL string,
	attributes map[string]string,
	sortOrder int,
) (*ProductSKU, error) {
	// 检查商品是否存在
	product, err := s.repo.GetByID(ctx, productID)
	if err != nil {
		return nil, fmt.Errorf("获取商品失败: %w", err)
	}
	if product == nil {
		return nil, ErrProductNotFound
	}

	// 检查SKU编码是否已存在
	exists, err := s.skuRepo.ExistsBySKUCode(ctx, skuCode)
	if err != nil {
		return nil, fmt.Errorf("检查SKU编码是否存在失败: %w", err)
	}
	if exists {
		return nil, ErrSKUCodeExists
	}

	// 创建SKU实体
	sku, err := NewProductSKU(
		productID, skuCode, name,
		marketPrice, salePrice, costPrice,
		0, // 初始库存为0
		weight, dimensions, imageURL,
		attributes, sortOrder,
	)
	if err != nil {
		return nil, err
	}

	// 保存SKU
	if err := s.skuRepo.Create(ctx, sku); err != nil {
		return nil, fmt.Errorf("创建SKU失败: %w", err)
	}

	return sku, nil
}

// UpdateProductSKU 更新商品SKU
func (s *DomainService) UpdateProductSKU(
	ctx context.Context,
	id, skuCode, name string,
	marketPrice, salePrice, costPrice decimal.Decimal,
	weight decimal.Decimal,
	dimensions map[string]any,
	imageURL string,
	attributes map[string]string,
	isActive bool,
	sortOrder int,
) (*ProductSKU, error) {
	// 获取现有SKU
	sku, err := s.skuRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("获取SKU失败: %w", err)
	}
	if sku == nil {
		return nil, ErrSKUNotFound
	}

	// 检查SKU编码是否已存在（排除当前SKU）
	if skuCode != sku.SKUCode {
		exists, err := s.skuRepo.ExistsBySKUCode(ctx, skuCode, id)
		if err != nil {
			return nil, fmt.Errorf("检查SKU编码是否存在失败: %w", err)
		}
		if exists {
			return nil, ErrSKUCodeExists
		}
	}

	// 更新SKU信息
	if err := sku.Update(
		skuCode, name,
		marketPrice, salePrice, costPrice,
		sku.StockQuantity, // 保持原有库存
		weight, dimensions, imageURL,
		attributes, isActive, sortOrder,
	); err != nil {
		return nil, err
	}

	// 保存更新
	if err := s.skuRepo.Update(ctx, sku); err != nil {
		return nil, fmt.Errorf("更新SKU失败: %w", err)
	}

	return sku, nil
}

// DeleteProductSKU 删除商品SKU
func (s *DomainService) DeleteProductSKU(ctx context.Context, id string) error {
	// 获取SKU
	sku, err := s.skuRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("获取SKU失败: %w", err)
	}
	if sku == nil {
		return ErrSKUNotFound
	}

	// 删除SKU
	if err := s.skuRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("删除SKU失败: %w", err)
	}

	return nil
}

// UpdateStock 更新库存
func (s *DomainService) UpdateStock(ctx context.Context, skuID string, quantity int, operation string) error {
	// 获取SKU
	sku, err := s.skuRepo.GetByID(ctx, skuID)
	if err != nil {
		return fmt.Errorf("获取SKU失败: %w", err)
	}
	if sku == nil {
		return ErrSKUNotFound
	}

	// 更新库存
	switch operation {
	case "add":
		sku.AddStock(quantity)
	case "subtract":
		if err := sku.SubtractStock(quantity); err != nil {
			return err
		}
	default:
		return fmt.Errorf("无效的库存操作: %s", operation)
	}

	// 保存更新
	if err := s.skuRepo.Update(ctx, sku); err != nil {
		return fmt.Errorf("更新库存失败: %w", err)
	}

	return nil
}

