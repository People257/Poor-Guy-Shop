package application

import (
	"github.com/google/wire"

	inventoryApp "github.com/people257/poor-guy-shop/inventory-service/internal/application/inventory"
	reservationApp "github.com/people257/poor-guy-shop/inventory-service/internal/application/reservation"
)

// ProviderSet 应用层依赖注入
var ProviderSet = wire.NewSet(
	inventoryApp.NewService,
	inventoryApp.NewBusinessService,
	inventoryApp.NewEventHandler,
	reservationApp.NewService,
)
