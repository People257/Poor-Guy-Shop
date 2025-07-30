package config

import (
	"fmt"
	"github.com/people257/poor-guy-shop/common/gateway/config"
	"github.com/spf13/viper"
)

type Config struct {
	config.GatewayConfig `mapstructure:",squash"`
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

	return &c
}

func GetGatewayConfig(cfg *Config) *config.GatewayConfig {
	if cfg == nil {
		panic("gateway config is nil")
	}
	return &cfg.GatewayConfig
}
