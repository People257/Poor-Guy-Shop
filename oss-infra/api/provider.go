package api

import (
	"github.com/google/wire"

	"github.com/people257/poor-guy-shop/oss-infra/api/file"
)

// HandlerProviderSet Handler providers
var HandlerProviderSet = wire.NewSet(
	file.NewHandler,
)
