package infra

import (
	"github.com/google/wire"
	"github.com/people257/poor-guy-shop/user-service/internal/infra/auth"
	"github.com/people257/poor-guy-shop/user-service/internal/infra/captcha"
	"github.com/people257/poor-guy-shop/user-service/internal/infra/email"
	"github.com/people257/poor-guy-shop/user-service/internal/infra/repository"
)

// InfraProviderSet Infrastructure providers
var InfraProviderSet = wire.NewSet(
	repository.NewUserRepository,
	repository.NewRefreshTokenRepository,
	auth.NewTokenService,
	auth.NewAuth,
	email.NewSMTPService,
	captcha.NewEmailCaptchaService,
	// 配置提供者
	ProvideAuthInfraConfig,
)
