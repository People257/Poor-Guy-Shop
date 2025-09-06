package config

import (
	"github.com/people257/poor-guy-shop/common/db"
	"github.com/people257/poor-guy-shop/common/server/config"
	"github.com/people257/poor-guy-shop/payment-service/internal/infra/payment"
)

// GetGrpcServerConfig 获取gRPC服务器配置
func GetGrpcServerConfig(cfg *Config) *config.GrpcServerConfig {
	return &cfg.GrpcServerConfig
}

// GetDBConfig 获取数据库配置
func GetDBConfig(cfg *Config) *db.DatabaseConfig {
	return &cfg.Database
}

// GetRedisConfig 获取Redis配置
func GetRedisConfig(cfg *Config) *db.RedisConfig {
	return &cfg.Redis
}

// GetAlipayConfig 获取支付宝配置
func GetAlipayConfig(cfg *Config) *payment.AlipayConfig {
	return &cfg.Payment.Alipay
}

// GetServicesConfig 获取服务依赖配置
func GetServicesConfig(cfg *Config) *ServicesConfig {
	return &cfg.Services
}

