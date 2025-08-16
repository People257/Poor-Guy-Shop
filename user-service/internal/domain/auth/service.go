package auth

import (
	"context"
	"time"
)

// TokenPair 令牌对
type TokenPair struct {
	AccessToken      string
	RefreshToken     string
	AccessExpiresIn  int32 // access token 过期时间（秒）
	RefreshExpiresIn int32 // refresh token 过期时间（秒）
}

// AuthResult 认证结果
type AuthResult struct {
	UserID    string
	TokenPair *TokenPair
}

// Service 认证服务接口
type Service interface {
	// GenerateTokens 生成访问令牌和刷新令牌
	GenerateTokens(ctx context.Context, userID string) (*TokenPair, error)

	// VerifyAccessToken 验证访问令牌
	VerifyAccessToken(ctx context.Context, token string) (*Claims, error)

	// RefreshTokens 使用刷新令牌生成新的令牌对
	RefreshTokens(ctx context.Context, refreshToken string) (*TokenPair, error)

	// RevokeRefreshToken 撤销刷新令牌
	RevokeRefreshToken(ctx context.Context, refreshToken string) error

	// RevokeAllTokens 撤销用户的所有令牌
	RevokeAllTokens(ctx context.Context, userID string) error
}

// Claims JWT声明
type Claims struct {
	UserID    string    `json:"user_id"`
	ExpiresAt time.Time `json:"exp"`
	IssuedAt  time.Time `json:"iat"`
}

// ExpiresAtMilli 返回毫秒时间戳
func (c *Claims) ExpiresAtMilli() int64 {
	return c.ExpiresAt.UnixMilli()
}

// IssuedAtMilli 返回毫秒时间戳
func (c *Claims) IssuedAtMilli() int64 {
	return c.IssuedAt.UnixMilli()
}

// CaptchaService 验证码服务接口
type CaptchaService interface {
	// SendEmailOTP 发送邮箱验证码
	SendEmailOTP(ctx context.Context, email string, purpose string) error

	// VerifyEmailOTP 验证邮箱验证码
	VerifyEmailOTP(ctx context.Context, email, otp, purpose string) error
}
