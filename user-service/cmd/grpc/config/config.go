package config

import (
	"github.com/people257/poor-guy-shop/common/conf"
	"github.com/people257/poor-guy-shop/common/db"
	"github.com/people257/poor-guy-shop/common/server/config"
)

type Config struct {
	config.GrpcServerConfig `mapstructure:",squash"`
	Redis                   db.RedisConfig    `mapstructure:"redis"`
	Database                db.DatabaseConfig `mapstructure:"database"`
	Auth                    AuthConfig        `mapstructure:"auth"`
	Captcha                 CaptchaConfig     `mapstructure:"captcha"`
	Email                   EmailConfig       `mapstructure:"email"`
}

func MustLoad(path string) *Config {
	_, c := conf.MustLoad[Config](path)
	return &c
}

func GetGrpcServerConfig(cfg *Config) *config.GrpcServerConfig {
	if cfg == nil {
		panic("grpc server config is nil")
	}
	return &cfg.GrpcServerConfig
}

func GetDBConfig(cfg *Config) *db.DatabaseConfig {
	if cfg == nil {
		panic("database config is nil")
	}
	return &cfg.Database
}

func GetRedisConfig(cfg *Config) *db.RedisConfig {
	if cfg == nil {
		panic("redis config is nil")
	}
	return &cfg.Redis
}

// AuthConfig JWT认证配置
type AuthConfig struct {
	JWT JWTConfig `mapstructure:"jwt"`
}

type JWTConfig struct {
	AccessToken  TokenConfig `mapstructure:"access_token"`
	RefreshToken TokenConfig `mapstructure:"refresh_token"`
}

type TokenConfig struct {
	Secret    string `mapstructure:"secret"`
	ExpiresIn int    `mapstructure:"expires_in"`
	Issuer    string `mapstructure:"issuer"`
}

// CaptchaConfig 验证码配置
type CaptchaConfig struct {
	Email EmailCaptchaConfig `mapstructure:"email"`
}

type EmailCaptchaConfig struct {
	Enabled      bool `mapstructure:"enabled"`
	CodeLength   int  `mapstructure:"code_length"`
	ExpiresIn    int  `mapstructure:"expires_in"`
	SendInterval int  `mapstructure:"send_interval"`
	DailyLimit   int  `mapstructure:"daily_limit"`
}

// EmailConfig 邮件配置
type EmailConfig struct {
	SMTP      SMTPConfig               `mapstructure:"smtp"`
	Templates map[string]EmailTemplate `mapstructure:"templates"`
}

type SMTPConfig struct {
	Host     string `mapstructure:"host"`
	Port     int    `mapstructure:"port"`
	Username string `mapstructure:"username"`
	Password string `mapstructure:"password"`
	From     string `mapstructure:"from"`
	UseTLS   bool   `mapstructure:"use_tls"`
}

type EmailTemplate struct {
	Subject string `mapstructure:"subject"`
	Body    string `mapstructure:"body"`
}
