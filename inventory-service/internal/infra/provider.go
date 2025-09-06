package infra

import (
	"github.com/google/wire"

	"github.com/people257/poor-guy-shop/inventory-service/internal/domain/inventory"
	"github.com/people257/poor-guy-shop/inventory-service/internal/domain/reservation"
	"github.com/people257/poor-guy-shop/inventory-service/internal/infra/client"
	"github.com/people257/poor-guy-shop/inventory-service/internal/infra/repository"
)

// ProviderSet 基础设施层依赖注入
var ProviderSet = wire.NewSet(
	// Repository
	repository.NewInventoryRepository,
	repository.NewInventoryLogRepository,
	repository.NewReservationRepository,

	// Client Manager
	client.NewManager,

	// Domain Service
	inventory.NewDomainService,
	reservation.NewDomainService,

	// Wire bindings
	wire.Bind(new(inventory.Repository), new(*repository.InventoryRepository)),
	wire.Bind(new(inventory.LogRepository), new(*repository.InventoryLogRepository)),
	wire.Bind(new(reservation.Repository), new(*repository.ReservationRepository)),
)
