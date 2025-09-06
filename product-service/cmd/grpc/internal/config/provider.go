package config

import (
	"github.com/google/wire"

	"github.com/people257/poor-guy-shop/common/db"
)

// ConfigProviderSet 配置提供者集合
var ConfigProviderSet = wire.NewSet(
	GetDatabaseConfig,
	GetRedisConfig,
)

// GetDatabaseConfig 获取数据库配置
func GetDatabaseConfig(cfg *Config) *db.DatabaseConfig {
	return &cfg.Database
}

// GetRedisConfig 获取Redis配置
func GetRedisConfig(cfg *Config) *db.RedisConfig {
	return &cfg.Redis
}
