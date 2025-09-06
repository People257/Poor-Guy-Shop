package middleware

import (
	"context"
	"strings"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"

	"github.com/people257/poor-guy-shop/order-service/internal/infra/client"
)

// AuthInterceptor 认证拦截器
type AuthInterceptor struct {
	userClient *client.UserServiceClient
}

// NewAuthInterceptor 创建认证拦截器
func NewAuthInterceptor(userClient *client.UserServiceClient) *AuthInterceptor {
	return &AuthInterceptor{
		userClient: userClient,
	}
}

// UnaryInterceptor gRPC一元调用认证拦截器
func (a *AuthInterceptor) UnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		// 跳过不需要认证的方法
		if a.skipAuth(info.FullMethod) {
			return handler(ctx, req)
		}

		// 从metadata中获取token
		token, err := a.extractToken(ctx)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "认证失败: %v", err)
		}

		// 验证token
		userID, err := a.userClient.AuthenticateUser(ctx, token)
		if err != nil {
			return nil, status.Errorf(codes.Unauthenticated, "token验证失败: %v", err)
		}

		// 将用户ID添加到context中
		ctx = context.WithValue(ctx, "user_id", userID)

		return handler(ctx, req)
	}
}

// StreamInterceptor gRPC流调用认证拦截器
func (a *AuthInterceptor) StreamInterceptor() grpc.StreamServerInterceptor {
	return func(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
		// 跳过不需要认证的方法
		if a.skipAuth(info.FullMethod) {
			return handler(srv, ss)
		}

		// 从metadata中获取token
		token, err := a.extractToken(ss.Context())
		if err != nil {
			return status.Errorf(codes.Unauthenticated, "认证失败: %v", err)
		}

		// 验证token
		userID, err := a.userClient.AuthenticateUser(ss.Context(), token)
		if err != nil {
			return status.Errorf(codes.Unauthenticated, "token验证失败: %v", err)
		}

		// 创建新的context
		ctx := context.WithValue(ss.Context(), "user_id", userID)

		// 包装ServerStream
		wrappedStream := &wrappedServerStream{
			ServerStream: ss,
			ctx:          ctx,
		}

		return handler(srv, wrappedStream)
	}
}

// extractToken 从metadata中提取token
func (a *AuthInterceptor) extractToken(ctx context.Context) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "missing metadata")
	}

	values := md.Get("authorization")
	if len(values) == 0 {
		return "", status.Error(codes.Unauthenticated, "missing authorization header")
	}

	authHeader := values[0]
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return "", status.Error(codes.Unauthenticated, "invalid authorization header format")
	}

	return strings.TrimPrefix(authHeader, "Bearer "), nil
}

// skipAuth 判断是否跳过认证
func (a *AuthInterceptor) skipAuth(method string) bool {
	// 这里可以定义不需要认证的方法
	skipMethods := []string{
		// 健康检查等公共方法
		"/grpc.health.v1.Health/Check",
		"/grpc.health.v1.Health/Watch",
	}

	for _, skipMethod := range skipMethods {
		if method == skipMethod {
			return true
		}
	}

	return false
}

// wrappedServerStream 包装ServerStream以提供新的context
type wrappedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (w *wrappedServerStream) Context() context.Context {
	return w.ctx
}

// GetUserIDFromContext 从context中获取用户ID
func GetUserIDFromContext(ctx context.Context) (string, bool) {
	userID, ok := ctx.Value("user_id").(string)
	return userID, ok
}
