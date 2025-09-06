package category

import "context"

// Repository 分类仓储接口
type Repository interface {
	// Create 创建分类
	Create(ctx context.Context, category *Category) error

	// Update 更新分类
	Update(ctx context.Context, category *Category) error

	// Delete 删除分类
	Delete(ctx context.Context, id string) error

	// GetByID 根据ID获取分类
	GetByID(ctx context.Context, id string) (*Category, error)

	// GetBySlug 根据slug获取分类
	GetBySlug(ctx context.Context, slug string) (*Category, error)

	// List 获取分类列表
	List(ctx context.Context, params ListParams) ([]*Category, int64, error)

	// GetTree 获取分类树
	GetTree(ctx context.Context, activeOnly bool) ([]*Category, error)

	// GetByParentID 根据父分类ID获取子分类
	GetByParentID(ctx context.Context, parentID string) ([]*Category, error)

	// CountChildren 统计子分类数量
	CountChildren(ctx context.Context, id string) (int64, error)

	// CountProducts 统计分类下商品数量
	CountProducts(ctx context.Context, id string) (int64, error)

	// ExistsBySlug 检查slug是否存在
	ExistsBySlug(ctx context.Context, slug string, excludeID ...string) (bool, error)

	// GetMaxLevel 获取最大层级
	GetMaxLevel(ctx context.Context) (int, error)

	// UpdateLevel 更新分类层级
	UpdateLevel(ctx context.Context, id string, level int) error

	// BatchUpdateLevel 批量更新分类层级
	BatchUpdateLevel(ctx context.Context, updates []LevelUpdate) error
}

// ListParams 列表查询参数
type ListParams struct {
	Page      int
	PageSize  int
	ParentID  *string
	Level     *int
	IsActive  *bool
	Keyword   string
	SortBy    string // name, sort_order, created_at
	SortOrder string // asc, desc
}

// LevelUpdate 层级更新参数
type LevelUpdate struct {
	ID    string
	Level int
}
