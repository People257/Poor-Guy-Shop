package config

import "fmt"

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

type GrpcServerConfig struct {
	Server ServerConfig `mapstructure:"server"`
	Log    LogConfig    `mapstructure:"log"`
}

func ValidateEnv(env string) error {
	if env != EnvDev && env != EnvTest && env != EnvProd {
		return fmt.Errorf("invalid env: %s", env)
	}
	return nil
}

func GetServerConfig(config *Config) *ServerConfig {
	if config == nil {
		return nil
	}
	return &config.Server
}
