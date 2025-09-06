package infra

import (
	"github.com/google/wire"

	"github.com/people257/poor-guy-shop/order-service/internal/infra/client"
	"github.com/people257/poor-guy-shop/order-service/internal/infra/repository"
)

// ProviderSet 基础设施层提供者集合
var ProviderSet = wire.NewSet(
	repository.NewOrderRepository,
	repository.NewCartRepository,
	client.ClientProviderSet,
)
