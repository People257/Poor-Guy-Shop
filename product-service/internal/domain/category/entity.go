package category

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/uuid"
)

// Category 分类实体
type Category struct {
	ID          string    `json:"id" gorm:"type:uuid;primary_key;default:gen_random_uuid()"`
	Name        string    `json:"name" gorm:"type:varchar(100);not null"`
	Slug        string    `json:"slug" gorm:"type:varchar(100);uniqueIndex;not null"`
	Description string    `json:"description" gorm:"type:text"`
	ParentID    *string   `json:"parent_id" gorm:"type:uuid;index"`
	Level       int       `json:"level" gorm:"type:integer;not null;default:1"`
	SortOrder   int       `json:"sort_order" gorm:"type:integer;not null;default:0"`
	IconURL     string    `json:"icon_url" gorm:"type:varchar(500)"`
	BannerURL   string    `json:"banner_url" gorm:"type:varchar(500)"`
	IsActive    bool      `json:"is_active" gorm:"type:boolean;not null;default:true"`
	CreatedAt   time.Time `json:"created_at" gorm:"type:timestamp with time zone;not null;default:now()"`
	UpdatedAt   time.Time `json:"updated_at" gorm:"type:timestamp with time zone;not null;default:now()"`

	// 关联关系
	Parent   *Category   `json:"parent,omitempty" gorm:"foreignKey:ParentID"`
	Children []*Category `json:"children,omitempty" gorm:"foreignKey:ParentID"`
}

// TableName 指定表名
func (Category) TableName() string {
	return "categories"
}

// NewCategory 创建新分类
func NewCategory(name, slug, description string, parentID *string, sortOrder int, iconURL, bannerURL string) (*Category, error) {
	if err := validateCategoryName(name); err != nil {
		return nil, err
	}

	if err := validateCategorySlug(slug); err != nil {
		return nil, err
	}

	category := &Category{
		ID:          uuid.New().String(),
		Name:        strings.TrimSpace(name),
		Slug:        strings.TrimSpace(slug),
		Description: strings.TrimSpace(description),
		ParentID:    parentID,
		Level:       1, // 默认为1级，实际会根据父分类计算
		SortOrder:   sortOrder,
		IconURL:     strings.TrimSpace(iconURL),
		BannerURL:   strings.TrimSpace(bannerURL),
		IsActive:    true,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	return category, nil
}

// Update 更新分类信息
func (c *Category) Update(name, slug, description string, parentID *string, sortOrder int, iconURL, bannerURL string, isActive bool) error {
	if err := validateCategoryName(name); err != nil {
		return err
	}

	if err := validateCategorySlug(slug); err != nil {
		return err
	}

	c.Name = strings.TrimSpace(name)
	c.Slug = strings.TrimSpace(slug)
	c.Description = strings.TrimSpace(description)
	c.ParentID = parentID
	c.SortOrder = sortOrder
	c.IconURL = strings.TrimSpace(iconURL)
	c.BannerURL = strings.TrimSpace(bannerURL)
	c.IsActive = isActive
	c.UpdatedAt = time.Now()

	return nil
}

// Deactivate 停用分类
func (c *Category) Deactivate() {
	c.IsActive = false
	c.UpdatedAt = time.Now()
}

// Activate 激活分类
func (c *Category) Activate() {
	c.IsActive = true
	c.UpdatedAt = time.Now()
}

// SetLevel 设置分类级别
func (c *Category) SetLevel(level int) {
	if level > 0 {
		c.Level = level
		c.UpdatedAt = time.Now()
	}
}

// HasParent 是否有父分类
func (c *Category) HasParent() bool {
	return c.ParentID != nil && *c.ParentID != ""
}

// IsRoot 是否为根分类
func (c *Category) IsRoot() bool {
	return !c.HasParent()
}

// CanBeDeleted 是否可以被删除
func (c *Category) CanBeDeleted() bool {
	// 如果有子分类，不能删除
	return len(c.Children) == 0
}

// validateCategoryName 验证分类名称
func validateCategoryName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return ErrCategoryNameRequired
	}
	if len(name) > 100 {
		return ErrCategoryNameTooLong
	}
	return nil
}

// validateCategorySlug 验证分类slug
func validateCategorySlug(slug string) error {
	slug = strings.TrimSpace(slug)
	if slug == "" {
		return ErrCategorySlugRequired
	}
	if len(slug) > 100 {
		return ErrCategorySlugTooLong
	}

	// 检查slug格式（只能包含字母、数字、连字符和下划线）
	for _, r := range slug {
		if !((r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') ||
			(r >= '0' && r <= '9') || r == '-' || r == '_') {
			return fmt.Errorf("slug只能包含字母、数字、连字符和下划线")
		}
	}

	return nil
}
