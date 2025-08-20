package infra

import (
	"time"

	"github.com/people257/poor-guy-shop/user-service/cmd/grpc/config"

	"github.com/people257/poor-guy-shop/user-service/internal/infra/auth"
)

// ProvideEmailConfig 提供邮件配置
func ProvideEmailConfig(cfg *config.Config) *config.EmailConfig {
	return &cfg.Email
}

// ProvideCaptchaConfig 提供验证码配置
func ProvideCaptchaConfig(cfg *config.Config) *config.CaptchaConfig {
	return &cfg.Captcha
}

// ProvideAuthConfig 提供认证配置
func ProvideAuthConfig(cfg *config.Config) *config.AuthConfig {
	return &cfg.Auth
}

// ProvideAuthInfraConfig 提供auth包的Config
func ProvideAuthInfraConfig(cfg *config.Config) *auth.Config {
	return &auth.Config{
		Secret:           cfg.Auth.JWT.AccessToken.Secret,
		ExpireDuration:   time.Duration(cfg.Auth.JWT.AccessToken.ExpiresIn) * time.Second,
		RefreshThreshold: 300 * time.Second, // 5分钟刷新阈值
	}
}
