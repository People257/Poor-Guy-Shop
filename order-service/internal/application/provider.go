package application

import (
	"github.com/google/wire"

	"github.com/people257/poor-guy-shop/order-service/internal/application/cart"
	"github.com/people257/poor-guy-shop/order-service/internal/application/order"
)

// ProviderSet 应用服务提供者集合
var ProviderSet = wire.NewSet(
	order.NewService,
	cart.NewService,
)
