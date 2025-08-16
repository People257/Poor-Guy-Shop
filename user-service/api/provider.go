package api

import (
	"github.com/google/wire"
	"github.com/people257/poor-guy-shop/user-service/api/auth"
	"github.com/people257/poor-guy-shop/user-service/api/info"
)

// APIProviderSet API providers
var APIProviderSet = wire.NewSet(
	auth.NewAuthServer,
	info.NewInfoServer,
)

// HandlerProviderSet Handler providers
var HandlerProviderSet = wire.NewSet()
