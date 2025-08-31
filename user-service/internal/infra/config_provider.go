package infra

import (
	"github.com/people257/poor-guy-shop/user-service/internal/config"
	"github.com/people257/poor-guy-shop/user-service/internal/infra/auth"
)

// ProvideAuthInfraConfig 提供auth包的Config
func ProvideAuthInfraConfig(jwtCfg *config.JWTConfig) *auth.Config {
	return &auth.Config{
		Secret:           jwtCfg.Secret,
		ExpireDuration:   jwtCfg.ExpireDuration,
		RefreshThreshold: jwtCfg.RefreshThresholdDuration,
	}
}
