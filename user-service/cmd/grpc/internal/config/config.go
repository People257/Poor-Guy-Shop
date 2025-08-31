package config

import (
	"time"

	"github.com/people257/poor-guy-shop/common/conf"
	"github.com/people257/poor-guy-shop/common/db"
	"github.com/people257/poor-guy-shop/common/server/config"
)

type DatabaseConfig struct {
	Host     string `mapstructure:"host"`
	Port     uint16 `mapstructure:"port"`
	Database string `mapstructure:"database"`
	User     string `mapstructure:"user"`
	Password string `mapstructure:"password"`
}

type Config struct {
	config.GrpcServerConfig `mapstructure:",squash"`
	Database                db.DatabaseConfig `mapstructure:"database"`
	Redis                   db.RedisConfig    `mapstructure:"redis"`
	JWT                     JWTConfig         `mapstructure:"jwt"`
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

// JWTConfig JWT认证配置
type JWTConfig struct {
	Secret                   string        `mapstructure:"secret"`
	ExpireDuration           time.Duration `mapstructure:"expire_duration"`            // Token 过期时间
	RefreshThresholdDuration time.Duration `mapstructure:"refresh_threshold_duration"` //刷新门限
}

func GetJWTConfig(cfg *Config) *JWTConfig {
	if cfg == nil {
		panic("jwt config is nil")
	}
	return &cfg.JWT
}

// CaptchaConfig 验证码配置
type CaptchaConfig struct {
	Provider       string `mapstructure:"provider"`        // Captcha 服务提供商
	Secret         string `mapstructure:"secret"`          // Captcha 密钥
	Endpoint       string `mapstructure:"endpoint"`        // Captcha 服务端点
	ExpectedDomain string `mapstructure:"expected_domain"` // 预期的域名
}

func GetCaptchaConfig(cfg *Config) *CaptchaConfig {
	if cfg == nil {
		panic("captcha config is nil")
	}
	return &cfg.Captcha
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

func GetEmailConfig(cfg *Config) *EmailConfig {
	if cfg == nil {
		panic("email config is nil")
	}
	return &cfg.Email
}
