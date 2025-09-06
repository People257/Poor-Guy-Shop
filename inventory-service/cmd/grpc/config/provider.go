package config

import (
	"github.com/google/wire"
	"github.com/people257/poor-guy-shop/common/db"
	"github.com/people257/poor-guy-shop/common/server/config"
)

// ProviderSet 配置依赖注入
var ProviderSet = wire.NewSet(
	NewConfig,
	GetDatabaseConfig,
	GetRedisConfig,
	GetGrpcServerConfig,
	GetServicesConfig,
)

// GetDatabaseConfig 获取数据库配置
func GetDatabaseConfig(cfg *Config) *db.DatabaseConfig {
	return &cfg.Database
}

// GetRedisConfig 获取Redis配置
func GetRedisConfig(cfg *Config) *db.RedisConfig {
	return &cfg.Redis
}

// GetGrpcServerConfig 获取gRPC服务配置
func GetGrpcServerConfig(cfg *Config) *config.GrpcServerConfig {
	return &cfg.GrpcServerConfig
}

// GetServicesConfig 获取服务依赖配置
func GetServicesConfig(cfg *Config) *ServicesConfig {
	return &cfg.Services
}

// NewConfig 创建配置
func NewConfig(configPath string) *Config {
	return MustLoad(configPath)
}
