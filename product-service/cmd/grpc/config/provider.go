package config

import (
	"github.com/google/wire"
)

// ConfigProviderSet 配置提供者集合
var ConfigProviderSet = wire.NewSet(
	MustLoad,
	GetRedisConfig,
	GetDBConfig,
	GetGrpcServerConfig,
)
