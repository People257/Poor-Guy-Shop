package inventory

import "errors"

var (
	// ErrInventoryNotFound 库存记录不存在
	ErrInventoryNotFound = errors.New("inventory not found")

	// ErrInsufficientInventory 库存不足
	ErrInsufficientInventory = errors.New("insufficient inventory")

	// ErrInsufficientReservedInventory 预占库存不足
	ErrInsufficientReservedInventory = errors.New("insufficient reserved inventory")

	// ErrReservationNotFound 预占记录不存在
	ErrReservationNotFound = errors.New("reservation not found")

	// ErrReservationExpired 预占已过期
	ErrReservationExpired = errors.New("reservation expired")

	// ErrReservationAlreadyConfirmed 预占已确认
	ErrReservationAlreadyConfirmed = errors.New("reservation already confirmed")

	// ErrReservationAlreadyReleased 预占已释放
	ErrReservationAlreadyReleased = errors.New("reservation already released")

	// ErrInvalidQuantity 无效的数量
	ErrInvalidQuantity = errors.New("invalid quantity")

	// ErrInvalidSkuID 无效的SKU ID
	ErrInvalidSkuID = errors.New("invalid sku id")

	// ErrInvalidOrderID 无效的订单ID
	ErrInvalidOrderID = errors.New("invalid order id")
)
