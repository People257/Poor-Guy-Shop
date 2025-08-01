package config

type RedisConfig struct {
	Host     string `mapstructure:"host"`
	Port     uint16 `mapstructure:"port"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

func GetRedisConfig(cfg *Config) *RedisConfig {
	if cfg == nil {
		panic("redis config is nil")
	}
	return &cfg.Redis
}
