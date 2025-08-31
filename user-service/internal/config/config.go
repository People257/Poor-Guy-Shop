package config

import (
	"time"
)

// JWTConfig JWT认证配置
type JWTConfig struct {
	Secret                   string        `mapstructure:"secret"`
	ExpireDuration           time.Duration `mapstructure:"expire_duration"`            // Token 过期时间
	RefreshThresholdDuration time.Duration `mapstructure:"refresh_threshold_duration"` //刷新门限
}

// CaptchaConfig 验证码配置
type CaptchaConfig struct {
	Provider       string             `mapstructure:"provider"`        // Captcha 服务提供商
	Secret         string             `mapstructure:"secret"`          // Captcha 密钥
	Endpoint       string             `mapstructure:"endpoint"`        // Captcha 服务端点
	ExpectedDomain string             `mapstructure:"expected_domain"` // 预期的域名
	Email          EmailCaptchaConfig `mapstructure:"email"`           // 邮箱验证码配置
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
