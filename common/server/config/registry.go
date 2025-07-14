package config

type RegistryConfig struct {
	Address string `mapstructure:"address"`
}

func GetRegistryConfig(config *GrpcServerConfig) *RegistryConfig {
	if config == nil {
		panic("registry config is nil")
	}
	return &config.Registry
}
