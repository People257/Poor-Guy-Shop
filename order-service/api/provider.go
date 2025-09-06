package api

import (
	"github.com/google/wire"

	"github.com/people257/poor-guy-shop/order-service/api/cart"
	"github.com/people257/poor-guy-shop/order-service/api/order"
)

// ProviderSet API层提供者集合
var ProviderSet = wire.NewSet(
	order.NewGrpcHandler,
	cart.NewGrpcHandler,
)
