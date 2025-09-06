package config

import (
	"github.com/people257/poor-guy-shop/common/conf"
	"github.com/people257/poor-guy-shop/common/db"
	"github.com/people257/poor-guy-shop/common/server/config"
)

// Config 产品服务配置
type Config struct {
	config.GrpcServerConfig `mapstructure:",squash"`
	Database                db.DatabaseConfig `mapstructure:"database"`
	Redis                   db.RedisConfig    `mapstructure:"redis"`
}

// MustLoad 加载配置文件
func MustLoad(path string) *Config {
	_, c := conf.MustLoad[Config](path)
	return &c
}

// GetGrpcServerConfig 获取gRPC服务器配置
func GetGrpcServerConfig(cfg *Config) *config.GrpcServerConfig {
	if cfg == nil {
		panic("grpc server config is nil")
	}
	return &cfg.GrpcServerConfig
}

// GetDBConfig 获取数据库配置
func GetDBConfig(cfg *Config) *db.DatabaseConfig {
	if cfg == nil {
		panic("database config is nil")
	}
	return &cfg.Database
}

// GetRedisConfig 获取Redis配置
func GetRedisConfig(cfg *Config) *db.RedisConfig {
	if cfg == nil {
		panic("redis config is nil")
	}
	return &cfg.Redis
}
