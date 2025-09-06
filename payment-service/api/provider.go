package api

import (
	"github.com/google/wire"
	"github.com/people257/poor-guy-shop/payment-service/api/payment"
)

// ProviderSet API层提供者集合
var ProviderSet = wire.NewSet(
	payment.NewGrpcHandler,
)
