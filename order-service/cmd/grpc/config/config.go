package config

import (
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/people257/poor-guy-shop/common/db"
	"github.com/people257/poor-guy-shop/common/server/config"
)

// ServiceConfig 服务连接配置
type ServiceConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// ServicesConfig 所有服务配置
type ServicesConfig struct {
	UserService      ServiceConfig `mapstructure:"user_service"`
	ProductService   ServiceConfig `mapstructure:"product_service"`
	PaymentService   ServiceConfig `mapstructure:"payment_service"`
	InventoryService ServiceConfig `mapstructure:"inventory_service"`
}

// Config 应用配置
type Config struct {
	GrpcServerConfig config.GrpcServerConfig `mapstructure:",squash"`
	Database         db.DatabaseConfig       `mapstructure:"database"`
	Redis            db.RedisConfig          `mapstructure:"redis"`
	Services         ServicesConfig          `mapstructure:"services"`
}

// MustLoad 加载配置
func MustLoad(path string) *Config {
	k := koanf.New(".")

	if err := k.Load(file.Provider(path), yaml.Parser()); err != nil {
		panic(err)
	}

	var cfg Config
	if err := k.Unmarshal("", &cfg); err != nil {
		panic(err)
	}

	return &cfg
}
