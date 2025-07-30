package client

import (
	gatewayconfig "poor-guy-shop/common/gateway/config"
	grpcconfig "poor-guy-shop/common/server/config"
)

type Config struct {
	registryAddr string // 注册中心地址
}

func NewConfigFromGRPCConfig(cfg *grpcconfig.GrpcServerConfig) *Config {
	if cfg == nil {
		panic("grpc config is nil")
	}
	return &Config{
		registryAddr: cfg.Registry.Address,
	}
}

func NewConfigFromGatewayConfig(cfg *gatewayconfig.GatewayConfig) *Config {
	if cfg == nil {
		panic("gateway config is nil")
	}
	return &Config{
		registryAddr: cfg.Registry.Address,
	}
}

func (c *Config) RegistryAddr() string {
	return c.registryAddr
}
