package config

import (
	"github.com/google/wire"
)

// ConfigProviderSet GrpcServerConfig providers
var ConfigProviderSet = wire.NewSet(
	GetObservabilityConfig,
	GetRegistryConfig,
	GetServerConfig,
	GetLogConfig,
)
