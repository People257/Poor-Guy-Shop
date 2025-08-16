package info

import (
	"context"

	"github.com/people257/poor-guy-shop/user-service/internal/domain/user"
)

// Service 用户信息应用服务
type Service struct {
	userRepo user.Repository
}

// NewService 创建用户信息应用服务
func NewService(userRepo user.Repository) *Service {
	return &Service{
		userRepo: userRepo,
	}
}

// UserInfoResponse 用户信息响应
type UserInfoResponse struct {
	UserID          string
	Nickname        string
	Phone           string
	Email           string
	Status          string
	RegisterChannel string
	CreatedAt       int64 // 毫秒时间戳
	UpdatedAt       int64 // 毫秒时间戳
}

// UpdateUserInfoRequest 更新用户信息请求
type UpdateUserInfoRequest struct {
	UserID          string
	Nickname        string
	Phone           string
	Email           string
	Status          string
	RegisterChannel string
}

// GetUserInfo 获取用户信息
func (s *Service) GetUserInfo(ctx context.Context, userID string) (*UserInfoResponse, error) {
	// 查找用户
	u, err := s.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, user.ErrUserNotFound
	}

	return s.toUserInfoResponse(u), nil
}

// UpdateUserInfo 更新用户信息
func (s *Service) UpdateUserInfo(ctx context.Context, req *UpdateUserInfoRequest) (*UserInfoResponse, error) {
	// 查找用户
	u, err := s.userRepo.FindByID(ctx, req.UserID)
	if err != nil {
		return nil, err
	}
	if u == nil {
		return nil, user.ErrUserNotFound
	}

	// 更新用户信息
	if req.Email != "" && req.Email != *u.Email {
		if err := u.ChangeEmail(req.Email); err != nil {
			return nil, err
		}
	}

	// 保存用户
	if err := s.userRepo.Update(ctx, u); err != nil {
		return nil, err
	}

	return s.toUserInfoResponse(u), nil
}

// toUserInfoResponse 转换为用户信息响应
func (s *Service) toUserInfoResponse(u *user.User) *UserInfoResponse {
	var phone, email string
	if u.PhoneNumber != nil {
		phone = *u.PhoneNumber
	}
	if u.Email != nil {
		email = *u.Email
	}

	return &UserInfoResponse{
		UserID:          u.ID,
		Nickname:        u.Username, // 暂时使用用户名作为昵称
		Phone:           phone,
		Email:           email,
		Status:          s.getUserStatusString(u.Status),
		RegisterChannel: "email", // 暂时固定为邮箱注册
		CreatedAt:       u.CreatedAt.UnixMilli(),
		UpdatedAt:       u.UpdatedAt.UnixMilli(),
	}
}

// getUserStatusString 获取用户状态字符串
func (s *Service) getUserStatusString(status user.UserStatus) string {
	switch status {
	case user.UserStatusNormal:
		return "active"
	case user.UserStatusDisabled:
		return "disabled"
	default:
		return "unknown"
	}
}
