package info

import (
	"context"

	infopb "github.com/people257/poor-guy-shop/user-service/gen/proto/user/info"
	"github.com/people257/poor-guy-shop/user-service/internal/application/info"
	infraAuth "github.com/people257/poor-guy-shop/user-service/internal/infra/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ infopb.InfoServiceServer = (*InfoServer)(nil)

// InfoServer 用户信息服务器
type InfoServer struct {
	infoService *info.Service
	authInfra   *infraAuth.Auth
}

// NewInfoServer 创建用户信息服务器
func NewInfoServer(infoService *info.Service, authInfra *infraAuth.Auth) *InfoServer {
	return &InfoServer{
		infoService: infoService,
		authInfra:   authInfra,
	}
}

// GetUserInfo 获取用户信息
func (s *InfoServer) GetUserInfo(ctx context.Context, empty *emptypb.Empty) (*infopb.GetUserInfoResponse, error) {
	// 从Token中获取用户ID
	claims, err := s.authInfra.VerifyFromMetadata(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "认证失败")
	}

	// 获取用户信息
	userInfo, err := s.infoService.GetUserInfo(ctx, claims.UserID())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &infopb.GetUserInfoResponse{
		UserId:          userInfo.UserID,
		Nickname:        userInfo.Nickname,
		Phone:           userInfo.Phone,
		Email:           userInfo.Email,
		Status:          userInfo.Status,
		RegisterChannel: userInfo.RegisterChannel,
		CreatedAt:       userInfo.CreatedAt,
		UpdatedAt:       userInfo.UpdatedAt,
	}, nil
}

// UpdateUserInfo 更新用户信息
func (s *InfoServer) UpdateUserInfo(ctx context.Context, req *infopb.UpdateUserInfoRequest) (*infopb.UpdateUserInfoResponse, error) {
	// 从Token中获取用户ID
	claims, err := s.authInfra.VerifyFromMetadata(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "认证失败")
	}

	// 验证请求参数
	if req.Email != "" {
		// 简单的邮箱格式验证
		if !isValidEmail(req.Email) {
			return nil, status.Error(codes.InvalidArgument, "邮箱格式不正确")
		}
	}

	// 更新用户信息
	updateReq := &info.UpdateUserInfoRequest{
		UserID:          claims.UserID(),
		Nickname:        req.Nickname,
		Phone:           req.Phone,
		Email:           req.Email,
		Status:          req.Status,
		RegisterChannel: req.RegisterChannel,
	}

	userInfo, err := s.infoService.UpdateUserInfo(ctx, updateReq)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &infopb.UpdateUserInfoResponse{
		UserId:          userInfo.UserID,
		Nickname:        userInfo.Nickname,
		Phone:           userInfo.Phone,
		Email:           userInfo.Email,
		Status:          userInfo.Status,
		RegisterChannel: userInfo.RegisterChannel,
		CreatedAt:       userInfo.CreatedAt,
		UpdatedAt:       userInfo.UpdatedAt,
	}, nil
}

// isValidEmail 简单的邮箱格式验证
func isValidEmail(email string) bool {
	// 这里可以使用更复杂的邮箱验证逻辑
	return len(email) > 0 && len(email) < 255
}
