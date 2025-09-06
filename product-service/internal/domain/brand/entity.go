package brand

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// Brand 品牌实体
type Brand struct {
	ID          string    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"type:varchar(100);not null"`
	Slug        string    `json:"slug" gorm:"type:varchar(100);uniqueIndex;not null"`
	Description string    `json:"description" gorm:"type:text"`
	LogoURL     string    `json:"logo_url" gorm:"type:varchar(500)"`
	WebsiteURL  string    `json:"website_url" gorm:"type:varchar(500)"`
	SortOrder   int       `json:"sort_order" gorm:"type:integer;not null;default:0"`
	IsActive    bool      `json:"is_active" gorm:"type:boolean;not null;default:true"`
	CreatedAt   time.Time `json:"created_at" gorm:"type:timestamp with time zone;not null;default:now()"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"type:timestamp with time zone;not null;default:now()"`
}

// TableName 指定表名
func (Brand) TableName() string {
	return "brands"
}

// NewBrand 创建新品牌
func NewBrand(name, slug, description, logoURL, websiteURL string, sortOrder int) (*Brand, error) {
	if err := validateBrandName(name); err != nil {
		return nil, err
	}

	if err := validateBrandSlug(slug); err != nil {
		return nil, err
	}

	brand := &Brand{
		ID:          uuid.New().String(),
		Name:        strings.TrimSpace(name),
		Slug:        strings.TrimSpace(slug),
		Description: strings.TrimSpace(description),
		LogoURL:     strings.TrimSpace(logoURL),
		WebsiteURL:  strings.TrimSpace(websiteURL),
		SortOrder:   sortOrder,
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return brand, nil
}

// Update 更新品牌信息
func (b *Brand) Update(name, slug, description, logoURL, websiteURL string, sortOrder int, isActive bool) error {
	if err := validateBrandName(name); err != nil {
		return err
	}

	if err := validateBrandSlug(slug); err != nil {
		return err
	}

	b.Name = strings.TrimSpace(name)
	b.Slug = strings.TrimSpace(slug)
	b.Description = strings.TrimSpace(description)
	b.LogoURL = strings.TrimSpace(logoURL)
	b.WebsiteURL = strings.TrimSpace(websiteURL)
	b.SortOrder = sortOrder
	b.IsActive = isActive
	b.UpdatedAt = time.Now()

	return nil
}

// Deactivate 停用品牌
func (b *Brand) Deactivate() {
	b.IsActive = false
	b.UpdatedAt = time.Now()
}

// Activate 激活品牌
func (b *Brand) Activate() {
	b.IsActive = true
	b.UpdatedAt = time.Now()
}

// validateBrandName 验证品牌名称
func validateBrandName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return ErrBrandNameRequired
	}
	if len(name) > 100 {
		return ErrBrandNameTooLong
	}
	return nil
}

// validateBrandSlug 验证品牌slug
func validateBrandSlug(slug string) error {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return ErrBrandSlugRequired
	}
	if len(slug) > 100 {
		return ErrBrandSlugTooLong
	}

	// 检查slug格式（只能包含字母、数字、连字符和下划线）
	for _, r := range slug {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') || r == '-' || r == '_') {
			return ErrBrandSlugInvalid
		}
	}

	return nil
}
