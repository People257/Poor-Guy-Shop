package address

import (
	"github.com/google/wire"
)

// ServiceProviderSet 地址应用服务提供者集合
var ServiceProviderSet = wire.NewSet(
	NewService,
)
