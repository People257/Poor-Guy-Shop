package brand

import "errors"

// 品牌相关错误
var (
	ErrBrandNotFound     = errors.New("品牌不存在")
	ErrBrandNameRequired = errors.New("品牌名称不能为空")
	ErrBrandNameTooLong  = errors.New("品牌名称长度不能超过100个字符")
	ErrBrandSlugRequired = errors.New("品牌slug不能为空")
	ErrBrandSlugTooLong  = errors.New("品牌slug长度不能超过100个字符")
	ErrBrandSlugExists   = errors.New("品牌slug已存在")
	ErrBrandSlugInvalid  = errors.New("品牌slug格式无效，只能包含字母、数字、连字符和下划线")
	ErrBrandHasProducts  = errors.New("品牌下还有商品，无法删除")
)
