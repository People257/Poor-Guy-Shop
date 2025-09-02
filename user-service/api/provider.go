package api

import (
	"github.com/google/wire"
	"github.com/people257/poor-guy-shop/user-service/api/address"
	"github.com/people257/poor-guy-shop/user-service/api/auth"
	"github.com/people257/poor-guy-shop/user-service/api/info"
)

// APIProviderSet API providers
var APIProviderSet = wire.NewSet(
	auth.NewAuthServer,
	info.NewInfoServer,
	address.NewAddressServer,
)

// HandlerProviderSet Handler providers
var HandlerProviderSet = wire.NewSet()
