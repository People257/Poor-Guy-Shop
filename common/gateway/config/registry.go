package config

type RegistryConfig struct {
	Service string `mapstructure:"service"`
	Address string `mapstructure:"address"`
}

func GetRegistryConfig(cfg *GatewayConfig) *RegistryConfig {
	if cfg == nil {
		panic("registry config is nil")
	}
	return &cfg.Registry
}
