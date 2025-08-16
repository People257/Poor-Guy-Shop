package auth

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/people257/poor-guy-shop/user-service/internal/domain/auth"
	"github.com/people257/poor-guy-shop/user-service/internal/infra/repository"
)

var _ auth.Service = (*TokenService)(nil)

// TokenService JWT令牌服务实现
type TokenService struct {
	auth             *Auth
	refreshTokenRepo auth.RefreshTokenRepository
	refreshTokenTTL  time.Duration
}

// NewTokenService 创建令牌服务
func NewTokenService(auth *Auth, refreshTokenRepo auth.RefreshTokenRepository) auth.Service {
	return &TokenService{
		auth:             auth,
		refreshTokenRepo: refreshTokenRepo,
		refreshTokenTTL:  24 * 7 * time.Hour, // 7天
	}
}

// GenerateTokens 生成访问令牌和刷新令牌
func (s *TokenService) GenerateTokens(ctx context.Context, userID string) (*auth.TokenPair, error) {
	now := time.Now()
	accessExpiresAt := now.Add(s.auth.ExpireDuration())

	// 生成访问令牌
	claims := &Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(accessExpiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			Subject:   userID,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString(s.auth.secret)
	if err != nil {
		return nil, fmt.Errorf("生成访问令牌失败: %w", err)
	}

	// 生成刷新令牌
	refreshToken, err := auth.NewRefreshToken(userID, s.refreshTokenTTL)
	if err != nil {
		return nil, fmt.Errorf("生成刷新令牌失败: %w", err)
	}

	// 存储刷新令牌
	if err := s.refreshTokenRepo.Store(ctx, refreshToken); err != nil {
		return nil, fmt.Errorf("存储刷新令牌失败: %w", err)
	}

	return &auth.TokenPair{
		AccessToken:      accessToken,
		RefreshToken:     refreshToken.Token,
		AccessExpiresIn:  int32(s.auth.ExpireDuration().Seconds()),
		RefreshExpiresIn: int32(s.refreshTokenTTL.Seconds()),
	}, nil
}

// VerifyAccessToken 验证访问令牌
func (s *TokenService) VerifyAccessToken(ctx context.Context, token string) (*auth.Claims, error) {
	claims, err := s.auth.Verify(token)
	if err != nil {
		return nil, err
	}

	expiresAt := time.Time{}
	issuedAt := time.Time{}

	if claims.ExpiresAt != nil {
		expiresAt = claims.ExpiresAt.Time
	}
	if claims.IssuedAt != nil {
		issuedAt = claims.IssuedAt.Time
	}

	return &auth.Claims{
		UserID:    claims.UserID(),
		ExpiresAt: expiresAt,
		IssuedAt:  issuedAt,
	}, nil
}

// RefreshTokens 使用刷新令牌生成新的令牌对
func (s *TokenService) RefreshTokens(ctx context.Context, refreshTokenStr string) (*auth.TokenPair, error) {
	// 查找刷新令牌
	refreshToken, err := s.refreshTokenRepo.FindByToken(ctx, refreshTokenStr)
	if err != nil {
		return nil, fmt.Errorf("查找刷新令牌失败: %w", err)
	}
	if refreshToken == nil {
		return nil, fmt.Errorf("刷新令牌不存在或已过期")
	}

	// 验证刷新令牌是否有效
	if !refreshToken.IsValid() {
		return nil, fmt.Errorf("刷新令牌无效或已过期")
	}

	// 删除旧的刷新令牌（一次性使用）
	if repo, ok := s.refreshTokenRepo.(*repository.RefreshTokenRepository); ok {
		if err := repo.DeleteToken(ctx, refreshTokenStr, refreshToken.UserID); err != nil {
			return nil, fmt.Errorf("删除旧刷新令牌失败: %w", err)
		}
	}

	// 生成新的令牌对
	return s.GenerateTokens(ctx, refreshToken.UserID)
}

// RevokeRefreshToken 撤销刷新令牌
func (s *TokenService) RevokeRefreshToken(ctx context.Context, refreshTokenStr string) error {
	// 查找刷新令牌
	refreshToken, err := s.refreshTokenRepo.FindByToken(ctx, refreshTokenStr)
	if err != nil {
		return fmt.Errorf("查找刷新令牌失败: %w", err)
	}
	if refreshToken == nil {
		return nil // 令牌不存在，视为已撤销
	}

	// 删除刷新令牌
	if repo, ok := s.refreshTokenRepo.(*repository.RefreshTokenRepository); ok {
		return repo.DeleteToken(ctx, refreshTokenStr, refreshToken.UserID)
	}

	return fmt.Errorf("无法撤销刷新令牌")
}

// RevokeAllTokens 撤销用户的所有令牌
func (s *TokenService) RevokeAllTokens(ctx context.Context, userID string) error {
	return s.refreshTokenRepo.DeleteByUserID(ctx, userID)
}

// generateRefreshToken 生成刷新令牌
func (s *TokenService) generateRefreshToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return hex.EncodeToString(bytes), nil
}
