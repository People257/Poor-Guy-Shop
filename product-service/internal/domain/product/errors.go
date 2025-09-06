package product

import "errors"

// 商品相关错误
var (
	ErrProductNotFound         = errors.New("商品不存在")
	ErrProductNameRequired     = errors.New("商品名称不能为空")
	ErrProductNameTooLong      = errors.New("商品名称长度不能超过200个字符")
	ErrProductSlugRequired     = errors.New("商品slug不能为空")
	ErrProductSlugTooLong      = errors.New("商品slug长度不能超过200个字符")
	ErrProductSlugExists       = errors.New("商品slug已存在")
	ErrProductSlugInvalid      = errors.New("商品slug格式无效，只能包含字母、数字、连字符和下划线")
	ErrProductPriceInvalid     = errors.New("商品价格不能为负数")
	ErrProductSalePriceHigher  = errors.New("商品销售价格不能高于原价")
	ErrProductImageURLsInvalid = errors.New("商品图片URLs格式无效")
	ErrCategoryNotFound        = errors.New("商品分类不存在")
	ErrBrandNotFound           = errors.New("商品品牌不存在")
)

// SKU相关错误
var (
	ErrSKUNotFound          = errors.New("商品SKU不存在")
	ErrSKUCodeRequired      = errors.New("SKU编码不能为空")
	ErrSKUCodeTooLong       = errors.New("SKU编码长度不能超过100个字符")
	ErrSKUCodeExists        = errors.New("SKU编码已存在")
	ErrSKUNameRequired      = errors.New("SKU名称不能为空")
	ErrSKUNameTooLong       = errors.New("SKU名称长度不能超过200个字符")
	ErrSKUAttributesInvalid = errors.New("SKU属性格式无效")
)
