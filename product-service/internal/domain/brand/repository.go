package brand

import "context"

// Repository 品牌仓储接口
type Repository interface {
	// Create 创建品牌
	Create(ctx context.Context, brand *Brand) error

	// Update 更新品牌
	Update(ctx context.Context, brand *Brand) error

	// Delete 删除品牌
	Delete(ctx context.Context, id string) error

	// GetByID 根据ID获取品牌
	GetByID(ctx context.Context, id string) (*Brand, error)

	// GetBySlug 根据slug获取品牌
	GetBySlug(ctx context.Context, slug string) (*Brand, error)

	// List 获取品牌列表
	List(ctx context.Context, params ListParams) ([]*Brand, int64, error)

	// CountProducts 统计品牌下商品数量
	CountProducts(ctx context.Context, id string) (int64, error)

	// ExistsBySlug 检查slug是否存在
	ExistsBySlug(ctx context.Context, slug string, excludeID ...string) (bool, error)
}

// ListParams 列表查询参数
type ListParams struct {
	Page      int
	PageSize  int
	IsActive  *bool
	Keyword   string
	SortBy    string // name, sort_order, created_at
	SortOrder string // asc, desc
}
