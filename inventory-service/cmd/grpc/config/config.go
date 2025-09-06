package config

import (
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/people257/poor-guy-shop/common/db"
	"github.com/people257/poor-guy-shop/common/server/config"
)

// ServiceConfig 服务配置
type ServiceConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// ServicesConfig 服务依赖配置
type ServicesConfig struct {
	OrderService   ServiceConfig `mapstructure:"order_service"`
	ProductService ServiceConfig `mapstructure:"product_service"`
	UserService    ServiceConfig `mapstructure:"user_service"`
}

// Config 主配置结构
type Config struct {
	GrpcServerConfig config.GrpcServerConfig `mapstructure:",squash"`
	Database         db.DatabaseConfig       `mapstructure:"database"`
	Redis            db.RedisConfig          `mapstructure:"redis"`
	Services         ServicesConfig          `mapstructure:"services"`
}

// MustLoad 加载配置
func MustLoad(configPath string) *Config {
	k := koanf.New(".")

	// 加载配置文件
	if err := k.Load(file.Provider(configPath), yaml.Parser()); err != nil {
		panic(err)
	}

	// 解析配置
	cfg := &Config{}
	if err := k.Unmarshal("", cfg); err != nil {
		panic(err)
	}

	return cfg
}
