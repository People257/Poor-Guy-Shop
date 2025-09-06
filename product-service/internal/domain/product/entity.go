package product

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
)

// ProductStatus 商品状态
type ProductStatus int

const (
	ProductStatusDraft ProductStatus = iota + 1
	ProductStatusActive
	ProductStatusInactive
	ProductStatusDeleted
)

// Product 商品实体
type Product struct {
	ID               string  `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name             string  `json:"name" gorm:"type:varchar(200);not null"`
	Slug             string  `json:"slug" gorm:"type:varchar(200);uniqueIndex;not null"`
	Description      string  `json:"description" gorm:"type:text"`
	ShortDescription string  `json:"short_description" gorm:"type:varchar(500)"`
	CategoryID       string  `json:"category_id" gorm:"type:uuid;not null;index"`
	BrandID          *string `json:"brand_id" gorm:"type:uuid;index"`

	// 价格信息
	MarketPrice decimal.Decimal `json:"market_price" gorm:"type:decimal(10,2)"`
	SalePrice   decimal.Decimal `json:"sale_price" gorm:"type:decimal(10,2);not null"`
	CostPrice   decimal.Decimal `json:"cost_price" gorm:"type:decimal(10,2)"`

	// 媒体资源
	MainImageURL string          `json:"main_image_url" gorm:"type:text"`
	ImageURLs    json.RawMessage `json:"image_urls" gorm:"type:jsonb"`
	VideoURL     string          `json:"video_url" gorm:"type:text"`

	// 商品属性
	Tags           json.RawMessage `json:"tags" gorm:"type:jsonb"`
	Specifications json.RawMessage `json:"specifications" gorm:"type:jsonb"`

	// 状态管理
	Status     ProductStatus `json:"status" gorm:"type:integer;not null;default:1"`
	IsFeatured bool          `json:"is_featured" gorm:"type:boolean;not null;default:false"`
	IsVirtual  bool          `json:"is_virtual" gorm:"type:boolean;not null;default:false"`

	// 发布时间管理
	PublishAt *time.Time `json:"publish_at" gorm:"type:timestamp with time zone"`

	// SEO信息
	SEOTitle       string `json:"seo_title" gorm:"type:varchar(200)"`
	SEODescription string `json:"seo_description" gorm:"type:text"`
	SEOKeywords    string `json:"seo_keywords" gorm:"type:varchar(500)"`

	// 排序权重
	SortOrder int `json:"sort_order" gorm:"type:integer;not null;default:0"`

	CreatedAt time.Time  `json:"created_at" gorm:"type:timestamp with time zone;not null;default:now()"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"type:timestamp with time zone;not null;default:now()"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"type:timestamp with time zone"`

	// 关联关系
	SKUs []*ProductSKU `json:"skus,omitempty" gorm:"foreignKey:ProductID"`
}

// TableName 指定表名
func (Product) TableName() string {
	return "products"
}

// ProductSKU 商品SKU实体
type ProductSKU struct {
	ID        string `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	ProductID string `json:"product_id" gorm:"type:uuid;not null;index"`
	SKUCode   string `json:"sku_code" gorm:"type:varchar(100);uniqueIndex;not null"`
	Name      string `json:"name" gorm:"type:varchar(200)"`

	// 价格信息
	MarketPrice decimal.Decimal `json:"market_price" gorm:"type:decimal(10,2)"`
	SalePrice   decimal.Decimal `json:"sale_price" gorm:"type:decimal(10,2);not null"`
	CostPrice   decimal.Decimal `json:"cost_price" gorm:"type:decimal(10,2)"`

	// 库存信息
	StockQuantity    int `json:"stock_quantity" gorm:"type:integer;not null;default:0"`
	ReservedQuantity int `json:"reserved_quantity" gorm:"type:integer;not null;default:0"`
	SoldQuantity     int `json:"sold_quantity" gorm:"type:integer;not null;default:0"`

	// SKU属性 (JSON格式存储规格信息)
	Attributes json.RawMessage `json:"attributes" gorm:"type:jsonb"`

	// 物理属性
	Weight     decimal.Decimal `json:"weight" gorm:"type:decimal(8,3)"`
	Dimensions json.RawMessage `json:"dimensions" gorm:"type:jsonb"`

	// 媒体资源
	ImageURL string `json:"image_url" gorm:"type:text"`

	// 状态管理
	Status int `json:"status" gorm:"type:integer;not null;default:1"`

	CreatedAt time.Time  `json:"created_at" gorm:"type:timestamp with time zone;not null;default:now()"`
	UpdatedAt time.Time  `json:"updated_at" gorm:"type:timestamp with time zone;not null;default:now()"`
	DeletedAt *time.Time `json:"deleted_at" gorm:"type:timestamp with time zone"`
}

// TableName 指定表名
func (ProductSKU) TableName() string {
	return "product_skus"
}

// NewProduct 创建新商品
func NewProduct(
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
	if err := validateProductName(name); err != nil {
		return nil, err
	}

	if err := validateProductSlug(slug); err != nil {
		return nil, err
	}

	if err := validateProductPrice(marketPrice, salePrice, costPrice); err != nil {
		return nil, err
	}

	// 序列化图片URLs
	imageURLsJSON, err := json.Marshal(imageURLs)
	if err != nil {
		return nil, ErrProductImageURLsInvalid
	}

	// 序列化标签
	tagsJSON, err := json.Marshal(tags)
	if err != nil {
		return nil, ErrProductImageURLsInvalid // 复用错误类型
	}

	// 序列化规格参数
	specificationsJSON, err := json.Marshal(specifications)
	if err != nil {
		return nil, ErrProductImageURLsInvalid // 复用错误类型
	}

	product := &Product{
		ID:               uuid.New().String(),
		Name:             strings.TrimSpace(name),
		Slug:             strings.TrimSpace(slug),
		Description:      strings.TrimSpace(description),
		ShortDescription: strings.TrimSpace(shortDescription),
		CategoryID:       categoryID,
		BrandID:          brandID,
		MarketPrice:      marketPrice,
		SalePrice:        salePrice,
		CostPrice:        costPrice,
		MainImageURL:     strings.TrimSpace(mainImageURL),
		ImageURLs:        imageURLsJSON,
		VideoURL:         strings.TrimSpace(videoURL),
		Tags:             tagsJSON,
		Specifications:   specificationsJSON,
		Status:           ProductStatusDraft,
		IsFeatured:       isFeatured,
		IsVirtual:        isVirtual,
		SEOTitle:         strings.TrimSpace(seoTitle),
		SEODescription:   strings.TrimSpace(seoDescription),
		SEOKeywords:      strings.TrimSpace(seoKeywords),
		SortOrder:        sortOrder,
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	return product, nil
}

// Update 更新商品信息
func (p *Product) Update(
	name, slug, description, shortDescription, categoryID string,
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
) error {
	if err := validateProductName(name); err != nil {
		return err
	}

	if err := validateProductSlug(slug); err != nil {
		return err
	}

	if err := validateProductPrice(marketPrice, salePrice, costPrice); err != nil {
		return err
	}

	// 序列化图片URLs
	imageURLsJSON, err := json.Marshal(imageURLs)
	if err != nil {
		return ErrProductImageURLsInvalid
	}

	// 序列化标签
	tagsJSON, err := json.Marshal(tags)
	if err != nil {
		return ErrProductImageURLsInvalid
	}

	// 序列化规格参数
	specificationsJSON, err := json.Marshal(specifications)
	if err != nil {
		return ErrProductImageURLsInvalid
	}

	p.Name = strings.TrimSpace(name)
	p.Slug = strings.TrimSpace(slug)
	p.Description = strings.TrimSpace(description)
	p.ShortDescription = strings.TrimSpace(shortDescription)
	p.CategoryID = categoryID
	p.BrandID = brandID
	p.MarketPrice = marketPrice
	p.SalePrice = salePrice
	p.CostPrice = costPrice
	p.MainImageURL = strings.TrimSpace(mainImageURL)
	p.ImageURLs = imageURLsJSON
	p.VideoURL = strings.TrimSpace(videoURL)
	p.Tags = tagsJSON
	p.Specifications = specificationsJSON
	p.Status = status
	p.IsFeatured = isFeatured
	p.IsVirtual = isVirtual
	p.SEOTitle = strings.TrimSpace(seoTitle)
	p.SEODescription = strings.TrimSpace(seoDescription)
	p.SEOKeywords = strings.TrimSpace(seoKeywords)
	p.SortOrder = sortOrder
	p.UpdatedAt = time.Now()

	return nil
}

// Activate 上架商品
func (p *Product) Activate() {
	p.Status = ProductStatusActive
	p.UpdatedAt = time.Now()
}

// Deactivate 下架商品
func (p *Product) Deactivate() {
	p.Status = ProductStatusInactive
	p.UpdatedAt = time.Now()
}

// Delete 删除商品
func (p *Product) Delete() {
	now := time.Now()
	p.DeletedAt = &now
	p.UpdatedAt = now
}

// IsActive 是否为激活状态
func (p *Product) IsActive() bool {
	return p.Status == ProductStatusActive
}

// GetImageURLs 获取图片URLs
func (p *Product) GetImageURLs() []string {
	var urls []string
	if len(p.ImageURLs) > 0 {
		json.Unmarshal(p.ImageURLs, &urls)
	}
	return urls
}

// GetTags 获取标签
func (p *Product) GetTags() []string {
	var tags []string
	if len(p.Tags) > 0 {
		json.Unmarshal(p.Tags, &tags)
	}
	return tags
}

// GetSpecifications 获取规格参数
func (p *Product) GetSpecifications() map[string]any {
	var specs map[string]any
	if len(p.Specifications) > 0 {
		json.Unmarshal(p.Specifications, &specs)
	}
	return specs
}

// NewProductSKU 创建新商品SKU
func NewProductSKU(
	productID, skuCode, name string,
	marketPrice, salePrice, costPrice decimal.Decimal,
	stockQuantity int,
	weight decimal.Decimal,
	dimensions map[string]any,
	imageURL string,
	attributes map[string]string,
	sortOrder int,
) (*ProductSKU, error) {
	if err := validateSKUCode(skuCode); err != nil {
		return nil, err
	}

	if err := validateSKUName(name); err != nil {
		return nil, err
	}

	if err := validateProductPrice(marketPrice, salePrice, costPrice); err != nil {
		return nil, err
	}

	// 序列化属性
	attributesJSON, err := json.Marshal(attributes)
	if err != nil {
		return nil, ErrSKUAttributesInvalid
	}

	// 序列化尺寸
	dimensionsJSON, err := json.Marshal(dimensions)
	if err != nil {
		return nil, ErrSKUAttributesInvalid
	}

	sku := &ProductSKU{
		ID:               uuid.New().String(),
		ProductID:        productID,
		SKUCode:          strings.TrimSpace(skuCode),
		Name:             strings.TrimSpace(name),
		MarketPrice:      marketPrice,
		SalePrice:        salePrice,
		CostPrice:        costPrice,
		StockQuantity:    stockQuantity,
		ReservedQuantity: 0,
		SoldQuantity:     0,
		Weight:           weight,
		Dimensions:       dimensionsJSON,
		ImageURL:         strings.TrimSpace(imageURL),
		Attributes:       attributesJSON,
		Status:           1, // 正常状态
		CreatedAt:        time.Now(),
		UpdatedAt:        time.Now(),
	}

	return sku, nil
}

// Update 更新SKU信息
func (s *ProductSKU) Update(
	skuCode, name string,
	marketPrice, salePrice, costPrice decimal.Decimal,
	stockQuantity int,
	weight decimal.Decimal,
	dimensions map[string]any,
	imageURL string,
	attributes map[string]string,
	isActive bool,
	sortOrder int,
) error {
	if err := validateSKUCode(skuCode); err != nil {
		return err
	}

	if err := validateSKUName(name); err != nil {
		return err
	}

	if err := validateProductPrice(marketPrice, salePrice, costPrice); err != nil {
		return err
	}

	// 序列化属性
	attributesJSON, err := json.Marshal(attributes)
	if err != nil {
		return ErrSKUAttributesInvalid
	}

	// 序列化尺寸
	dimensionsJSON, err := json.Marshal(dimensions)
	if err != nil {
		return ErrSKUAttributesInvalid
	}

	s.SKUCode = strings.TrimSpace(skuCode)
	s.Name = strings.TrimSpace(name)
	s.MarketPrice = marketPrice
	s.SalePrice = salePrice
	s.CostPrice = costPrice
	s.StockQuantity = stockQuantity
	s.Weight = weight
	s.Dimensions = dimensionsJSON
	s.ImageURL = strings.TrimSpace(imageURL)
	s.Attributes = attributesJSON
	s.Status = 1
	if !isActive {
		s.Status = 2 // 停售状态
	}
	s.UpdatedAt = time.Now()

	return nil
}

// GetAttributes 获取属性
func (s *ProductSKU) GetAttributes() map[string]string {
	var attrs map[string]string
	if len(s.Attributes) > 0 {
		json.Unmarshal(s.Attributes, &attrs)
	}
	return attrs
}

// GetDimensions 获取尺寸
func (s *ProductSKU) GetDimensions() map[string]any {
	var dims map[string]any
	if len(s.Dimensions) > 0 {
		json.Unmarshal(s.Dimensions, &dims)
	}
	return dims
}

// AddStock 增加库存
func (s *ProductSKU) AddStock(quantity int) {
	s.StockQuantity += quantity
	s.UpdatedAt = time.Now()
}

// SubtractStock 减少库存
func (s *ProductSKU) SubtractStock(quantity int) error {
	if s.StockQuantity < quantity {
		return fmt.Errorf("库存不足，当前库存: %d，需要: %d", s.StockQuantity, quantity)
	}
	s.StockQuantity -= quantity
	s.SoldQuantity += quantity
	s.UpdatedAt = time.Now()
	return nil
}

// ReserveStock 预留库存
func (s *ProductSKU) ReserveStock(quantity int) error {
	availableStock := s.StockQuantity - s.ReservedQuantity
	if availableStock < quantity {
		return fmt.Errorf("可用库存不足，可用库存: %d，需要: %d", availableStock, quantity)
	}
	s.ReservedQuantity += quantity
	s.UpdatedAt = time.Now()
	return nil
}

// ReleaseReservedStock 释放预留库存
func (s *ProductSKU) ReleaseReservedStock(quantity int) {
	if s.ReservedQuantity >= quantity {
		s.ReservedQuantity -= quantity
		s.UpdatedAt = time.Now()
	}
}

// 验证函数
func validateProductName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return ErrProductNameRequired
	}
	if len(name) > 200 {
		return ErrProductNameTooLong
	}
	return nil
}

func validateProductSlug(slug string) error {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return ErrProductSlugRequired
	}
	if len(slug) > 200 {
		return ErrProductSlugTooLong
	}

	// 检查slug格式
	for _, r := range slug {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') || r == '-' || r == '_') {
			return ErrProductSlugInvalid
		}
	}

	return nil
}

func validateProductPrice(marketPrice, salePrice, costPrice decimal.Decimal) error {
	if salePrice.LessThan(decimal.Zero) {
		return ErrProductPriceInvalid
	}
	if !marketPrice.IsZero() && marketPrice.LessThan(decimal.Zero) {
		return ErrProductPriceInvalid
	}
	if !costPrice.IsZero() && costPrice.LessThan(decimal.Zero) {
		return ErrProductPriceInvalid
	}
	if !marketPrice.IsZero() && salePrice.GreaterThan(marketPrice) {
		return ErrProductSalePriceHigher
	}
	return nil
}

func validateSKUCode(skuCode string) error {
	skuCode = strings.TrimSpace(skuCode)
	if skuCode == "" {
		return ErrSKUCodeRequired
	}
	if len(skuCode) > 100 {
		return ErrSKUCodeTooLong
	}
	return nil
}

func validateSKUName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return ErrSKUNameRequired
	}
	if len(name) > 200 {
		return ErrSKUNameTooLong
	}
	return nil
}
