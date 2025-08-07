package config

import (
	"fmt"
	"github.com/people257/poor-guy-shop/common/server/config"

	"github.com/spf13/viper"
)

type Config struct {
	config.GrpcServerConfig `mapstructure:",squash"`
	Database                DatabaseConfig `mapstructure:"database"`
	Redis                   RedisConfig    `mapstructure:"redis"`
}

func MustLoad(path string) *Config {
	v := viper.New()
	v.SetConfigFile(path)

	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("failed to read config file: %s", err))
	}

	var c Config
	if err := v.Unmarshal(&c); err != nil {
		panic(fmt.Errorf("failed to unmarshal config: %s", err))
	}

	// Validate environment type
	if err := config.ValidateEnv(c.Server.Env); err != nil {
		panic(err)
	}

	return &c
}

func GetGrpcServerConfig(cfg *Config) *config.GrpcServerConfig {
	if cfg == nil {
		panic("grpc server config is nil")
	}
	return &cfg.GrpcServerConfig
}
