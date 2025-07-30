package internal

import (
	"fmt"
	capi "github.com/hashicorp/consul/api"
	"go.uber.org/zap"
	"net"
	"poor-guy-shop/common/server/config"
)

type Register struct {
	client    *capi.Client
	serverCfg *config.ServerConfig
	ip        net.IP
}

func NewRegister(client *capi.Client, serverCfg *config.ServerConfig) *Register {
	ip, err := getOutboundIP()
	if err != nil {
		panic(err)
	}
	return &Register{client: client, serverCfg: serverCfg, ip: ip}
}

func (r *Register) SetIP(ip net.IP) {
	r.ip = ip
}

func (r *Register) CheckAndReRegisterService() error {
	client := r.client
	serverCfg := r.serverCfg
	ip := r.ip

	services, err := client.Agent().Services()
	if err != nil {
		return fmt.Errorf("failed to get services from consul: %w", err)
	}

	serviceID := fmt.Sprintf("%s-%s-%d", serverCfg.Name, ip, serverCfg.Port)

	if _, exists := services[serviceID]; !exists {
		zap.L().Info("service not registered, registering now", zap.String("service_id", serviceID))
		return r.RegisterService()
	}

	return nil
}

func (r *Register) RegisterService() error {
	client := r.client
	serverCfg := r.serverCfg
	ip, err := getOutboundIP()
	if err != nil {
		return err
	}

	check := &capi.AgentServiceCheck{
		GRPC:                           fmt.Sprintf("%s:%d", ip, serverCfg.Port),
		Timeout:                        "5s",
		Interval:                       "10s",
		DeregisterCriticalServiceAfter: "20s",
	}

	zap.L().Info("registering service to consul", zap.String("name", serverCfg.Name), zap.String("ip", ip.String()), zap.Uint16("port", serverCfg.Port))
	err = client.Agent().ServiceRegister(&capi.AgentServiceRegistration{
		ID:      r.GetServiceIDString(),
		Name:    serverCfg.Name,
		Address: ip.String(),
		Port:    int(serverCfg.Port),
		Check:   check,
	})
	if err != nil {
		return err
	}

	return nil
}

func (r *Register) DeregisterService() error {
	return r.client.Agent().ServiceDeregister(r.GetServiceIDString())
}

func (r *Register) GetServiceIDString() string {
	serverCfg := r.serverCfg
	ip := r.ip

	return fmt.Sprintf("%s-%s-%d", serverCfg.Name, ip.String(), serverCfg.Port)
}

func NewConsulClient(registryConfig *config.RegistryConfig) *capi.Client {
	cfg := capi.DefaultConfig()
	if registryConfig.Address != "" {
		cfg.Address = registryConfig.Address
	}
	client, err := capi.NewClient(cfg)
	if err != nil {
		panic(err)
	}

	return client
}

func getOutboundIP() (net.IP, error) {
	conn, err := net.Dial("udp", "1.1.1.1:1")
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	localAddr := conn.LocalAddr().(*net.UDPAddr)
	return localAddr.IP, nil
}
