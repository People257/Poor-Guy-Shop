package category

import "errors"

// 分类相关错误
var (
	ErrCategoryNotFound       = errors.New("分类不存在")
	ErrCategoryNameRequired   = errors.New("分类名称不能为空")
	ErrCategoryNameTooLong    = errors.New("分类名称长度不能超过100个字符")
	ErrCategorySlugRequired   = errors.New("分类slug不能为空")
	ErrCategorySlugTooLong    = errors.New("分类slug长度不能超过100个字符")
	ErrCategorySlugExists     = errors.New("分类slug已存在")
	ErrCategoryHasChildren    = errors.New("分类下还有子分类，无法删除")
	ErrCategoryHasProducts    = errors.New("分类下还有商品，无法删除")
	ErrCategoryCircularRef    = errors.New("分类不能设置自己或子分类为父分类")
	ErrCategoryLevelTooDeep   = errors.New("分类层级过深，最多支持5级")
	ErrParentCategoryNotFound = errors.New("父分类不存在")
	ErrParentCategoryInactive = errors.New("父分类未激活")
)
