package config

import (
	"github.com/google/wire"
	"github.com/people257/poor-guy-shop/common/db"
	"github.com/people257/poor-guy-shop/common/server/config"
)

// ConfigProviderSet Config providers
var ConfigProviderSet = wire.NewSet(
	GetRedisConfig,
	GetDBConfig,
	GetGrpcServerConfig,
	GetJWTConfig,
	GetCaptchaConfigNew,
	GetEmailConfig,
)

func GetRedisConfig(cfg *Config) *db.RedisConfig {
	return &cfg.Redis
}

func GetDBConfig(cfg *Config) *db.DatabaseConfig {
	return &cfg.Database
}

func GetGrpcServerConfig(cfg *Config) *config.GrpcServerConfig {
	return &cfg.GrpcServerConfig
}
