package main

import (
	"github.com/people257/poor-guy-shop/user-service/cmd/grpc/internal/config"
	internalConfig "github.com/people257/poor-guy-shop/user-service/internal/config"
)

// ProvideInternalEmailConfig 转换邮件配置
func ProvideInternalEmailConfig(cfg *config.Config) *internalConfig.EmailConfig {
	return &internalConfig.EmailConfig{
		SMTP: internalConfig.SMTPConfig{
			Host:     cfg.Email.SMTP.Host,
			Port:     cfg.Email.SMTP.Port,
			Username: cfg.Email.SMTP.Username,
			Password: cfg.Email.SMTP.Password,
			From:     cfg.Email.SMTP.From,
			UseTLS:   cfg.Email.SMTP.UseTLS,
		},
		Templates: make(map[string]internalConfig.EmailTemplate),
	}
}

// ProvideInternalCaptchaConfig 转换验证码配置
func ProvideInternalCaptchaConfig(cfg *config.Config) *internalConfig.CaptchaConfig {
	return &internalConfig.CaptchaConfig{
		Provider:       cfg.Captcha.Provider,
		Secret:         cfg.Captcha.Secret,
		Endpoint:       cfg.Captcha.Endpoint,
		ExpectedDomain: cfg.Captcha.ExpectedDomain,
		Email: internalConfig.EmailCaptchaConfig{
			Enabled:      true, // 默认启用
			CodeLength:   6,    // 默认6位数字
			ExpiresIn:    300,  // 默认5分钟过期
			SendInterval: 60,   // 默认60秒发送间隔
			DailyLimit:   10,   // 默认每日10次限制
		},
	}
}

// ProvideInternalJWTConfig 转换JWT配置
func ProvideInternalJWTConfig(cfg *config.Config) *internalConfig.JWTConfig {
	return &internalConfig.JWTConfig{
		Secret:                   cfg.JWT.Secret,
		ExpireDuration:           cfg.JWT.ExpireDuration,
		RefreshThresholdDuration: cfg.JWT.RefreshThresholdDuration,
	}
}
