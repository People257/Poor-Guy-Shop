package config

type LogConfig struct {
	Level   string           `mapstructure:"level"`
	File    FileLogConfig    `mapstructure:"file"`
	Console ConsoleLogConfig `mapstructure:"console"`
}

type FileLogConfig struct {
	Enable     bool   `mapstructure:"enable"`
	Directory  string `mapstructure:"directory"`
	Name       string `mapstructure:"name"`
	MaxSize    uint32 `mapstructure:"max_size"`
	MaxAge     uint32 `mapstructure:"max_age"`
	MaxBackups uint32 `mapstructure:"max_backups"`
	Compress   bool   `mapstructure:"compress"`
	LocalTime  bool   `mapstructure:"local_time"`
}

type ConsoleLogConfig struct {
	Enable bool   `mapstructure:"enable"`
	Format string `mapstructure:"format"`
}

func GetLogConfig(config *GrpcServerConfig) *LogConfig {
	if config == nil {
		panic("config is nil")
	}
	return &config.Log
}
