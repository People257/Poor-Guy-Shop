package address

import (
	"github.com/google/wire"
)

// HandlerProviderSet 地址API处理器提供者集合
var HandlerProviderSet = wire.NewSet(
	NewAddressServer,
)
