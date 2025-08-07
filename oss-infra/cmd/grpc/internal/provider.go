package internal

import "github.com/google/wire"

var InternalProviderSet = wire.NewSet(
	NewDB,
	NewRedisClient,
)
