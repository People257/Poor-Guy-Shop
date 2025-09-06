package client

import "errors"

// 客户端错误定义
var (
	ErrServiceUnavailable     = errors.New("service unavailable")
	ErrInvalidRequest         = errors.New("invalid request")
	ErrInsufficientStock      = errors.New("insufficient stock")
	ErrPaymentFailed          = errors.New("payment failed")
	ErrInventoryReserveFailed = errors.New("inventory reserve failed")
	ErrInventoryConfirmFailed = errors.New("inventory confirm failed")
	ErrInventoryReleaseFailed = errors.New("inventory release failed")
)

// ClientError 客户端错误类型
type ClientError struct {
	Service string
	Op      string
	Err     error
}

func (e *ClientError) Error() string {
	return e.Service + "." + e.Op + ": " + e.Err.Error()
}

func (e *ClientError) Unwrap() error {
	return e.Err
}

// NewClientError 创建客户端错误
func NewClientError(service, op string, err error) *ClientError {
	return &ClientError{
		Service: service,
		Op:      op,
		Err:     err,
	}
}
