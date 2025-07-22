package config

type RegistryConfig struct {
	Address string `mapstructure:"address"`
}

func GetRegistryConfig(cfg *GrpcServerConfig) *RegistryConfig {
	if cfg == nil {
		panic("registry config is nil")
	}
	return &cfg.Registry
}
