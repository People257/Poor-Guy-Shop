package client

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"

	"github.com/people257/poor-guy-shop/order-service/cmd/grpc/config"
)

// 临时定义，后续需要引入user-service的proto
type AuthService interface {
	AuthenticateRPC(ctx context.Context, in *emptypb.Empty, opts ...grpc.CallOption) (*AuthenticateRPCResp, error)
}

type AuthenticateRPCResp struct {
	UserID      string
	AccessToken string
	ExpiresIn   int64
}

// UserServiceClient 用户服务客户端
type UserServiceClient struct {
	conn        *grpc.ClientConn
	authService AuthService
}

// NewUserServiceClient 创建用户服务客户端
func NewUserServiceClient(cfg *config.ServiceConfig) (*UserServiceClient, error) {
	addr := fmt.Sprintf("%s:%d", cfg.Host, cfg.Port)

	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %w", err)
	}

	// 这里需要使用真实的user-service proto client
	// authService := authpb.NewAuthServiceClient(conn)

	return &UserServiceClient{
		conn: conn,
		// authService: authService,
	}, nil
}

// AuthenticateUser 验证用户身份
func (c *UserServiceClient) AuthenticateUser(ctx context.Context, token string) (string, error) {
	// 将token放入metadata
	md := metadata.Pairs("authorization", "Bearer "+token)
	ctx = metadata.NewOutgoingContext(ctx, md)

	// 调用用户服务验证
	// resp, err := c.authService.AuthenticateRPC(ctx, &emptypb.Empty{})
	// if err != nil {
	//     return "", fmt.Errorf("failed to authenticate user: %w", err)
	// }

	// 临时返回固定用户ID，实际应该从resp中获取
	return "temp-user-id", nil
}

// Close 关闭连接
func (c *UserServiceClient) Close() error {
	return c.conn.Close()
}
