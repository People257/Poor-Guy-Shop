package auth

import (
	"context"
	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/labstack/echo/v4"
	authpb "github.com/people257/poor-guy-shop/user-service/gen/proto/user/auth"
	semconv "go.opentelemetry.io/otel/semconv/v1.32.0"
	"go.opentelemetry.io/otel/trace"
	"go.uber.org/zap"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
	"strconv"
	"strings"
)

func BuildMetadataMiddleware(client authpb.) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			token := c.Request().Header.Get("Authorization")
			// 清除所有 Grpc-Metadata- 开头的 header
			for key := range c.Request().Header {
				if strings.HasPrefix(key, runtime.MetadataHeaderPrefix) {
					c.Request().Header.Del(key)
				}
			}
			if token == "" {
				return next(c)
			}

			ctx := metadata.AppendToOutgoingContext(c.Request().Context(), GrpcTokenMetadataKey, token)
			resp, err := client.AuthenticateRPC(ctx, &emptypb.Empty{})
			if err != nil {
				zap.L().Error("authenticateRPC failed", zap.Error(err))
				return next(c)
			}
			// 如果有新的 token，设置到响应头
			if resp.NewToken != nil {
				c.Response().Header().Set("New-Access-Token", resp.NewToken.AccessToken)
				c.Response().Header().Set("New-Expires-In", strconv.FormatInt(resp.NewToken.ExpiresIn, 10))
			}
			// 将 user_id 写入 HTTP header, 在请求时 Gateway 会自动将其写入 gRPC metadata
			c.Request().Header.Set(HttpUserIDHeaderKey, resp.UserID)

			c.SetRequest(c.Request().WithContext(context.WithValue(c.Request().Context(), userIDCtxKey, resp.UserID)))

			if span := trace.SpanFromContext(c.Request().Context()); span != nil && span.IsRecording() {
				span.SetAttributes(semconv.UserID(resp.UserID))
			}

			return next(c)
		}
	}
}
