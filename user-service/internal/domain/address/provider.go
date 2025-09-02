package address

import (
	"github.com/google/wire"
)

// DomainServiceProviderSet 地址领域服务提供者集合
var DomainServiceProviderSet = wire.NewSet(
	NewDomainService,
)
