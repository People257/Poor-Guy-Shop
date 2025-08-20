package config

import (
	"github.com/google/wire"

	"github.com/people257/poor-guy-shop/common/server/config"
)

// ConfigProviderSet Config providers
var ConfigProviderSet = wire.NewSet(
	GetRedisConfig,
	GetDBConfig,
	GetServerConfig,
)

func GetServerConfig(cfg *Config) *config.ServerConfig {
	if cfg == nil {
		panic("server config is nil")
	}
	return &cfg.Server
}
