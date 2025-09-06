package config

import (
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
	"github.com/people257/poor-guy-shop/common/db"
	"github.com/people257/poor-guy-shop/common/server/config"
	"github.com/people257/poor-guy-shop/payment-service/internal/infra/payment"
)

// PaymentConfig 支付配置
type PaymentConfig struct {
	Alipay payment.AlipayConfig `mapstructure:"alipay"`
	// Wechat WechatConfig `mapstructure:"wechat"`  // TODO: 添加微信支付配置
}

// ServiceConfig 服务配置
type ServiceConfig struct {
	Host string `mapstructure:"host"`
	Port int    `mapstructure:"port"`
}

// ServicesConfig 服务依赖配置
type ServicesConfig struct {
	OrderService ServiceConfig `mapstructure:"order_service"`
}

// Config 主配置结构
type Config struct {
	GrpcServerConfig config.GrpcServerConfig `mapstructure:",squash"`
	Database         db.DatabaseConfig       `mapstructure:"database"`
	Redis            db.RedisConfig          `mapstructure:"redis"`
	Payment          PaymentConfig           `mapstructure:"payment"`
	Services         ServicesConfig          `mapstructure:"services"`
}

// MustLoad 加载配置
func MustLoad(configPath string) *Config {
	k := koanf.New(".")

	// 加载配置文件
	if err := k.Load(file.Provider(configPath), nil); err != nil {
		panic(err)
	}

	// 解析配置
	cfg := &Config{}
	if err := k.Unmarshal("", cfg); err != nil {
		panic(err)
	}

	return cfg
}
