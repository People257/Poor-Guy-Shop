package config

type GrpcServerConfig struct {
	Server        ServerConfig        `mapstructure:"server"`
	Observability ObservabilityConfig `mapstructure:"observability"`
	Registry      RegistryConfig      `mapstructure:"registry"`
	Log           LogConfig           `mapstructure:"log"`
}
