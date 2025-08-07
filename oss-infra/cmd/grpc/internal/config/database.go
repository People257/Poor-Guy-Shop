package config

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     uint16 `mapstructure:"port"`
	Database string `mapstructure:"database"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

func GetDBConfig(cfg *Config) *DatabaseConfig {
	if cfg == nil {
		panic("database config is nil")
	}
	return &cfg.Database
}
