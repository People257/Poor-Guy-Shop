package config

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

func GetServerConfig(cfg *GatewayConfig) *ServerConfig {
	if cfg == nil {
		panic("server config is nil")
	}
	return &cfg.Server
}
