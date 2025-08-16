package application

import (
	"github.com/google/wire"
	"github.com/people257/poor-guy-shop/user-service/internal/application/auth"
	"github.com/people257/poor-guy-shop/user-service/internal/application/info"
)

// AppProviderSet Application providers
var AppProviderSet = wire.NewSet(
	auth.NewService,
	info.NewService,
)
