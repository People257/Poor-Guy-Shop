package auth

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"go.uber.org/zap"
	"time"
)

// GenTokenCmd 生成 jwt token 的命令
type GenTokenCmd struct {
	UserID string // 用户id
}

// Token 生成的 JWT TOKEN
type Token struct {
	ID          string        // jwt id
	AccessToken string        // jwt token
	ExpiresIn   time.Duration // jwt 过期时间
	Claims      jwt.Claims    // jwt 自定义字段
}

// GenerateToken 生成 jwt token
func (a *Auth) GenerateToken(t *GenTokenCmd) (*Token, error) {
	tokenID := uuid.New().String()
	claims := jwt.RegisteredClaims{
		Subject:   t.UserID,
		IssuedAt:  jwt.NewNumericDate(time.Now()),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(a.ExpireDuration())),
		ID:        tokenID,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(a.secret)
	if err != nil {
		zap.L().Error("generate token failed", zap.Error(err), zap.Any("claims", claims))
		return nil, ErrGenerateToken
	}
	return &Token{
		ID:          tokenID,
		AccessToken: tokenString,
		ExpiresIn:   a.ExpireDuration(),
		Claims:      claims,
	}, nil
}
