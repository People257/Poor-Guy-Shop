package config

import (
	"github.com/people257/poor-guy-shop/common/db"
	"github.com/people257/poor-guy-shop/common/server/config"
)

// Config 应用配置
type Config struct {
	Server        config.GrpcServerConfig    `mapstructure:"server"`
	Log           config.LogConfig           `mapstructure:"log"`
	Database      db.DatabaseConfig          `mapstructure:"database"`
	Redis         db.RedisConfig             `mapstructure:"redis"`
	Registry      config.RegistryConfig      `mapstructure:"registry"`
	Observability config.ObservabilityConfig `mapstructure:"observability"`
}
