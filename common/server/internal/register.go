package internal

import (
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/people257/poor-guy-shop/common/server/config"
	"go.uber.org/zap"
	"net"
)

type Register struct {
	client    *api.Client
	serverCfg *config.ServerConfig
	ip        string
	serviceID string
}

func NewConsulClient(registryConfig *config.ServerConfig) *api.Client {
	cfg := api.DefaultConfig()
	if registryConfig.Address != "" {
		cfg.Address = registryConfig.Address
	}
	client, err := api.NewClient(cfg)
	if err != nil {
		panic(fmt.Sprintf("failed to create consul client: %v", err))
	}
	return client
}

func NewRegister(serverCfg *config.ServerConfig) *Register {
	ip, err := getOUtboundIP()
	if err != nil {
		panic(fmt.Sprintf("failed to get outbound ip: %v", err))
	}
	serviceID := fmt.Sprintf("%s-%s-%d", serverCfg.Name, ip, serverCfg.Port)
	return &Register{
		client:    NewConsulClient(serverCfg),
		serverCfg: serverCfg,
		ip:        ip,
		serviceID: serviceID,
	}
}

func getOUtboundIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP.String(), nil
}

func (r *Register) RegisterService() error {
	check := &api.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", r.ip, r.serverCfg.Port),
		Timeout:                        "5s",
		Interval:                       "10s",
		DeregisterCriticalServiceAfter: "20s",
	}

	registration := &api.AgentServiceRegistration{
		ID:      r.serviceID,
		Name:    "manage-service",
		Address: r.ip,
		Port:    int(r.serverCfg.Port),
		Check:   check,
	}

	zap.L().Info("Registering service to Consul",
		zap.String("serviceID", r.serviceID),
		zap.String("ip", r.ip),
		zap.Int("port", int(r.serverCfg.Port)),
	)

	return r.client.Agent().ServiceRegister(registration)

}

func (r *Register) DeregisterService() error {
	zap.L().Info("Deregistering service from Consul", zap.String("serviceID", r.serviceID))
	return r.client.Agent().ServiceDeregister(r.serviceID)
}
