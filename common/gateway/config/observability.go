package config

import "time"

type ObservabilityConfig struct {
	Port        uint16                     `mapstructure:"port"`
	Pprof       PprofConfig                `mapstructure:"pprof"`
	Metrics     ObservabilityMetricsConfig `mapstructure:"metrics"`
	Trace       ObservabilityTraceConfig   `mapstructure:"trace"`
	Log         ObservabilityLogConfig     `mapstructure:"log"`
	Address     string                     `mapstructure:"address"` // 用于otlp exporter的地址
	Timeout     time.Duration              `mapstructure:"timeout"` // shutdown超时时间
	OTLPHeaders map[string]string          `mapstructure:"otlp_headers"`
}

type PprofConfig struct {
	Enable bool `mapstructure:"enable"`
}

type ObservabilityMetricsConfig struct {
	Enable bool `mapstructure:"enable"`
}

type ObservabilityTraceConfig struct {
	Enable bool `mapstructure:"enable"`
}

type ObservabilityLogConfig struct {
	Enable bool `mapstructure:"enable"`
}

func GetObservabilityConfig(cfg *GatewayConfig) *ObservabilityConfig {
	if cfg == nil {
		panic("observability config is nil")
	}
	return &cfg.Observability
}
