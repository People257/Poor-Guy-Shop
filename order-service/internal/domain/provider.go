package domain

import (
	"github.com/google/wire"

	"github.com/people257/poor-guy-shop/order-service/internal/domain/cart"
	"github.com/people257/poor-guy-shop/order-service/internal/domain/order"
)

// ProviderSet 领域服务依赖注入提供者集合
var ProviderSet = wire.NewSet(
	order.NewDomainService,
	cart.NewDomainService,
)
