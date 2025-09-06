package user

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/people257/poor-guy-shop/inventory-service/cmd/grpc/config"
)

// UserInfo 用户信息
type UserInfo struct {
	ID       uuid.UUID `json:"id"`
	Username string    `json:"username"`
	Email    string    `json:"email"`
	Status   string    `json:"status"`
	Role     string    `json:"role"`
}

// Client 用户服务客户端接口
type Client interface {
	// GetUser 获取用户信息
	GetUser(ctx context.Context, userID uuid.UUID) (*UserInfo, error)

	// ValidateUser 验证用户是否存在且活跃
	ValidateUser(ctx context.Context, userID uuid.UUID) (bool, error)

	// GetUserPermissions 获取用户权限
	GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]string, error)
}

// GrpcClient 用户服务gRPC客户端实现
type GrpcClient struct {
	conn   *grpc.ClientConn
	config *config.ServiceConfig
}

// NewGrpcClient 创建用户服务gRPC客户端
func NewGrpcClient(servicesConfig *config.ServicesConfig) (Client, error) {
	serviceConfig := servicesConfig.UserService

	// 建立gRPC连接
	addr := fmt.Sprintf("%s:%d", serviceConfig.Host, serviceConfig.Port)
	conn, err := grpc.Dial(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return nil, fmt.Errorf("failed to connect to user service: %w", err)
	}

	return &GrpcClient{
		conn:   conn,
		config: &serviceConfig,
	}, nil
}

// Close 关闭连接
func (c *GrpcClient) Close() error {
	if c.conn != nil {
		return c.conn.Close()
	}
	return nil
}

// GetUser 获取用户信息
func (c *GrpcClient) GetUser(ctx context.Context, userID uuid.UUID) (*UserInfo, error) {
	// TODO: 实现用户信息查询
	// 这里需要根据实际的用户服务proto定义来实现
	return &UserInfo{
		ID:       userID,
		Username: "sample_user",
		Email:    "user@example.com",
		Status:   "active",
		Role:     "user",
	}, nil
}

// ValidateUser 验证用户是否存在且活跃
func (c *GrpcClient) ValidateUser(ctx context.Context, userID uuid.UUID) (bool, error) {
	// TODO: 实现用户验证
	// 这里需要根据实际的用户服务proto定义来实现
	return true, nil
}

// GetUserPermissions 获取用户权限
func (c *GrpcClient) GetUserPermissions(ctx context.Context, userID uuid.UUID) ([]string, error) {
	// TODO: 实现用户权限查询
	// 这里需要根据实际的用户服务proto定义来实现
	return []string{"inventory:read", "inventory:write"}, nil
}

