package config

type LogConfig struct {
	// Log level configuration
	Level string `mapstructure:"level"`

	// File output configuration
	File struct {
		Enable    bool   `mapstructure:"enable"`
		Directory string `mapstructure:"directory"`
		Name      string `mapstructure:"name"`

		// Rotation configuration
		MaxSize    int  `mapstructure:"max_size"`    // Maximum size in megabytes
		MaxAge     int  `mapstructure:"max_age"`     // Maximum age in days
		MaxBackups int  `mapstructure:"max_backups"` // Maximum number of backups
		Compress   bool `mapstructure:"compress"`    // Compress rotated files
		LocalTime  bool `mapstructure:"local_time"`  // Use local time for rotation
	} `mapstructure:"file"`

	// Console output configuration
	Console struct {
		Enable bool   `mapstructure:"enable"`
		Format string `mapstructure:"format"` // json or console
	} `mapstructure:"console"`
}

func GetLogConfig(cfg *GrpcServerConfig) *LogConfig {
	if cfg == nil {
		panic("log config is nil")
	}
	return &cfg.Log
}
