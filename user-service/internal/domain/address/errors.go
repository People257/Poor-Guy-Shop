package address

import "errors"

// 地址相关错误定义
var (
	// ErrAddressNotFound 地址不存在
	ErrAddressNotFound = errors.New("地址不存在")

	// ErrAddressLimitExceeded 地址数量超限
	ErrAddressLimitExceeded = errors.New("地址数量已达上限")

	// ErrInvalidReceiverName 收货人姓名无效
	ErrInvalidReceiverName = errors.New("收货人姓名无效")

	// ErrInvalidReceiverPhone 收货人电话无效
	ErrInvalidReceiverPhone = errors.New("收货人电话无效")

	// ErrInvalidAddress 地址信息无效
	ErrInvalidAddress = errors.New("地址信息无效")

	// ErrInvalidPostalCode 邮政编码无效
	ErrInvalidPostalCode = errors.New("邮政编码无效")

	// ErrUnauthorizedAccess 无权限访问
	ErrUnauthorizedAccess = errors.New("无权限访问此地址")

	// ErrCannotDeleteDefaultAddress 不能删除默认地址
	ErrCannotDeleteDefaultAddress = errors.New("不能删除默认地址，请先设置其他地址为默认")
)
