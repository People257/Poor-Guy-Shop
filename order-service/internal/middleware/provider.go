package middleware

import (
	"github.com/google/wire"
)

// MiddlewareProviderSet 中间件提供器集合
var MiddlewareProviderSet = wire.NewSet(
	NewAuthInterceptor,
)
