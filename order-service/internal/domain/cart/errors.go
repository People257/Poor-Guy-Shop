package cart

import "errors"

// 购物车领域错误定义
var (
	ErrCartItemNotFound   = errors.New("cart item not found")
	ErrInvalidQuantity    = errors.New("invalid quantity")
	ErrProductNotFound    = errors.New("product not found")
	ErrSkuNotFound        = errors.New("sku not found")
	ErrInsufficientStock  = errors.New("insufficient stock")
	ErrProductUnavailable = errors.New("product unavailable")
	ErrCartEmpty          = errors.New("cart is empty")
	ErrNoSelectedItems    = errors.New("no selected items")
	ErrInvalidPrice       = errors.New("invalid price")
	ErrDuplicateCartItem  = errors.New("duplicate cart item")
)
