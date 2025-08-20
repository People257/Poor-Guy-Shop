package auth

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/people257/poor-guy-shop/common/auth"

	"github.com/golang-jwt/jwt/v5"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
)

type Auth struct {
	secret           []byte        // secret key
	expireDuration   time.Duration // token expire duration
	refreshThreshold time.Duration // refresh threshold
}

func NewAuth(cfg *Config) *Auth {
	secret := []byte(cfg.Secret)
	if len(secret) < 32 {
		panic("secret must be at least 32 bytes")
	}
	return &Auth{
		secret:           secret,
		expireDuration:   cfg.ExpireDuration,
		refreshThreshold: cfg.RefreshThreshold,
	}
}

// VerifyFromMetadata verifies the JWT token in the grpc metadata authorization header.
func (a *Auth) VerifyFromMetadata(ctx context.Context) (*Claims, error) {
	tokens := metadata.ValueFromIncomingContext(ctx, auth.GrpcTokenMetadataKey)
	if len(tokens) == 0 {
		return nil, ErrEmptyToken
	}

	token := tokens[0]
	token = strings.TrimPrefix(token, "Bearer ")
	return a.Verify(token)
}

// Verify verifies the JWT token.
func (a *Auth) Verify(token string) (*Claims, error) {
	c := &Claims{}

	_, err := jwt.ParseWithClaims(token, c, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, ErrInvalidToken
		}
		return a.secret, nil
	}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Name}),
		jwt.WithIssuedAt(),
		jwt.WithExpirationRequired(),
	)
	if err != nil {
		if !errors.Is(err, jwt.ErrTokenExpired) {
			zap.L().Warn("parse token error", zap.Error(err), zap.String("token", token))
		}
		return nil, ErrInvalidToken
	}

	return c, nil
}

// ExpireDuration returns the duration after which the token will expire.
func (a *Auth) ExpireDuration() time.Duration {
	return a.expireDuration
}

// RefreshThreshold returns the threshold duration for token refresh.
func (a *Auth) RefreshThreshold() time.Duration {
	return a.refreshThreshold
}
