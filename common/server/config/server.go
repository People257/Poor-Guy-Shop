package config

import (
	"fmt"
)

const (
	EnvDev  = "dev"
	EnvTest = "test"
	EnvProd = "prod"
)

type ServerConfig struct {
	Name string `mapstructure:"name"`
	Env  string `mapstructure:"env"`
	Port uint16 `mapstructure:"port"`
}

// ValidateEnv validates if the environment type is valid
func ValidateEnv(env string) error {
	if env != EnvDev && env != EnvTest && env != EnvProd {
		return fmt.Errorf("invalid environment type: %s, must be one of: %s, %s, %s",
			env, EnvDev, EnvTest, EnvProd)
	}
	return nil
}

func GetServerConfig(cfg *GrpcServerConfig) *ServerConfig {
	if cfg == nil {
		panic("server config is nil")
	}
	return &cfg.Server
}
