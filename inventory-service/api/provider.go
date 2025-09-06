package api

import (
	"github.com/google/wire"

	"github.com/people257/poor-guy-shop/inventory-service/api/inventory"
)

// ProviderSet API层依赖注入
var ProviderSet = wire.NewSet(
	inventory.NewServer,
)
