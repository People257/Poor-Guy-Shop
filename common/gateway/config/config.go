package config

type GatewayConfig struct {
	Server        ServerConfig        `mapstructure:"server"`
	Observability ObservabilityConfig `mapstructure:"observability"`
	Registry      RegistryConfig      `mapstructure:"registry"`
	Log           LogConfig           `mapstructure:"log"`
}
