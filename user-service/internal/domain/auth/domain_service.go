package auth

import (
	"context"
	"time"
)

// DomainService 认证领域服务
type DomainService struct {
	tokenService     Service
	captchaService   CaptchaService
	refreshTokenRepo RefreshTokenRepository
}

// NewDomainService 创建认证领域服务
func NewDomainService(
	tokenService Service,
	captchaService CaptchaService,
	refreshTokenRepo RefreshTokenRepository,
) *DomainService {
	return &DomainService{
		tokenService:     tokenService,
		captchaService:   captchaService,
		refreshTokenRepo: refreshTokenRepo,
	}
}

// GenerateUserTokens 为用户生成令牌对
func (s *DomainService) GenerateUserTokens(ctx context.Context, userID string) (*TokenPair, error) {
	return s.tokenService.GenerateTokens(ctx, userID)
}

// VerifyUserAccessToken 验证用户访问令牌
func (s *DomainService) VerifyUserAccessToken(ctx context.Context, token string) (*Claims, error) {
	return s.tokenService.VerifyAccessToken(ctx, token)
}

// RefreshUserTokens 使用刷新令牌生成新的令牌对
func (s *DomainService) RefreshUserTokens(ctx context.Context, refreshToken string) (*TokenPair, error) {
	return s.tokenService.RefreshTokens(ctx, refreshToken)
}

// RevokeUserRefreshToken 撤销用户的刷新令牌
func (s *DomainService) RevokeUserRefreshToken(ctx context.Context, refreshToken string) error {
	return s.tokenService.RevokeRefreshToken(ctx, refreshToken)
}

// RevokeAllUserTokens 撤销用户的所有令牌
func (s *DomainService) RevokeAllUserTokens(ctx context.Context, userID string) error {
	return s.tokenService.RevokeAllTokens(ctx, userID)
}

// SendEmailVerificationCode 发送邮箱验证码
func (s *DomainService) SendEmailVerificationCode(ctx context.Context, email, purpose string) error {
	return s.captchaService.SendEmailOTP(ctx, email, purpose)
}

// VerifyEmailCode 验证邮箱验证码
func (s *DomainService) VerifyEmailCode(ctx context.Context, email, code, purpose string) error {
	return s.captchaService.VerifyEmailOTP(ctx, email, code, purpose)
}

// ValidateTokenExpiry 验证令牌是否需要刷新
func (s *DomainService) ValidateTokenExpiry(claims *Claims, refreshThreshold time.Duration) bool {
	return time.Until(claims.ExpiresAt) < refreshThreshold
}
