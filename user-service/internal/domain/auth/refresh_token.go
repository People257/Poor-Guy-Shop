package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"time"
)

// RefreshToken 刷新令牌实体
type RefreshToken struct {
	ID        string    // 令牌ID
	UserID    string    // 用户ID
	Token     string    // 令牌值
	ExpiresAt time.Time // 过期时间
	CreatedAt time.Time // 创建时间
	IsUsed    bool      // 是否已使用
}

// RefreshTokenRepository 刷新令牌仓储接口
type RefreshTokenRepository interface {
	// Store 存储刷新令牌
	Store(ctx context.Context, token *RefreshToken) error

	// FindByToken 根据令牌值查找
	FindByToken(ctx context.Context, token string) (*RefreshToken, error)

	// MarkAsUsed 标记令牌为已使用
	MarkAsUsed(ctx context.Context, tokenID string) error

	// DeleteByUserID 删除用户的所有刷新令牌
	DeleteByUserID(ctx context.Context, userID string) error

	// DeleteExpired 删除过期的令牌
	DeleteExpired(ctx context.Context) error
}

// NewRefreshToken 创建新的刷新令牌
func NewRefreshToken(userID string, expiresIn time.Duration) (*RefreshToken, error) {
	token, err := generateRefreshToken()
	if err != nil {
		return nil, err
	}

	id, err := generateTokenID()
	if err != nil {
		return nil, err
	}

	now := time.Now()
	return &RefreshToken{
		ID:        id,
		UserID:    userID,
		Token:     token,
		ExpiresAt: now.Add(expiresIn),
		CreatedAt: now,
		IsUsed:    false,
	}, nil
}

// IsExpired 检查令牌是否过期
func (rt *RefreshToken) IsExpired() bool {
	return time.Now().After(rt.ExpiresAt)
}

// IsValid 检查令牌是否有效
func (rt *RefreshToken) IsValid() bool {
	return !rt.IsUsed && !rt.IsExpired()
}

// MarkUsed 标记令牌为已使用
func (rt *RefreshToken) MarkUsed() {
	rt.IsUsed = true
}

// generateRefreshToken 生成刷新令牌
func generateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}

// generateTokenID 生成令牌ID
func generateTokenID() (string, error) {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
