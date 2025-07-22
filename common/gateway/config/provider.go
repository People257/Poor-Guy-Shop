package config

import (
	"github.com/google/wire"
)

// ConfigProviderSet GatewayConfig providers
var ConfigProviderSet = wire.NewSet(
	GetServerConfig,
	GetObservabilityConfig,
	GetRegistryConfig,
	GetLogConfig,
)
