package config

import (
	"github.com/google/wire"
)

// ConfigProviderSet Config providers
var ConfigProviderSet = wire.NewSet(
	GetDBConfig,
	GetRedisConfig,
)
