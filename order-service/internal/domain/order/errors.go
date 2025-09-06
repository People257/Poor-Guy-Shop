package order

import "errors"

// 订单领域错误定义
var (
	ErrOrderNotFound        = errors.New("order not found")
	ErrOrderCannotCancel    = errors.New("order cannot be cancelled")
	ErrOrderCannotPay       = errors.New("order cannot be paid")
	ErrOrderCannotShip      = errors.New("order cannot be shipped")
	ErrOrderCannotConfirm   = errors.New("order cannot be confirmed")
	ErrInvalidOrderStatus   = errors.New("invalid order status")
	ErrInvalidPaymentMethod = errors.New("invalid payment method")
	ErrOrderItemNotFound    = errors.New("order item not found")
	ErrOrderAddressNotFound = errors.New("order address not found")
	ErrInsufficientStock    = errors.New("insufficient stock")
	ErrInvalidQuantity      = errors.New("invalid quantity")
	ErrOrderAlreadyPaid     = errors.New("order already paid")
	ErrOrderExpired         = errors.New("order expired")
)
