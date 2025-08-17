package application

import (
	"github.com/google/wire"

	"github.com/people257/poor-guy-shop/oss-infra/internal/application/file"
)

// AppProviderSet Application providers
var AppProviderSet = wire.NewSet(
	file.NewService,
)
