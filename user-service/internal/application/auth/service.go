package auth

import (
	"context"

	"github.com/people257/poor-guy-shop/user-service/internal/domain/auth"
	"github.com/people257/poor-guy-shop/user-service/internal/domain/user"
)

// Service 认证应用服务（负责编排不同领域服务）
type Service struct {
	userDomainService *user.DomainService
	authDomainService *auth.DomainService
	userRepo          user.Repository
}

// NewService 创建认证应用服务
func NewService(
	userDomainService *user.DomainService,
	authDomainService *auth.DomainService,
	userRepo user.Repository,
) *Service {
	return &Service{
		userDomainService: userDomainService,
		authDomainService: authDomainService,
		userRepo:          userRepo,
	}
}

// LoginRequest 登录请求
type LoginRequest struct {
	Account  string
	Password string
}

// LoginResponse 登录响应
type LoginResponse struct {
	UserID       string
	AccessToken  string
	RefreshToken string
	ExpiresIn    int32
}

// Login 用户登录（应用服务编排）
func (s *Service) Login(ctx context.Context, req *LoginRequest) (*LoginResponse, error) {
	// 1. 用户认证（委托给用户领域服务）
	authenticatedUser, err := s.userDomainService.AuthenticateUser(ctx, req.Account, req.Password)
	if err != nil {
		return nil, err
	}

	// 2. 生成令牌（委托给认证领域服务）
	tokens, err := s.authDomainService.GenerateUserTokens(ctx, authenticatedUser.ID)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		UserID:       authenticatedUser.ID,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    tokens.AccessExpiresIn,
	}, nil
}

// EmailRegisterRequest 邮箱注册请求
type EmailRegisterRequest struct {
	Email    string
	Password string
	Captcha  string
	Phone    *string
}

// RegisterResponse 注册响应
type RegisterResponse struct {
	UserID       string
	AccessToken  string
	RefreshToken string
	ExpiresIn    int32
}

// EmailRegister 邮箱注册（应用服务编排）
func (s *Service) EmailRegister(ctx context.Context, req *EmailRegisterRequest) (*RegisterResponse, error) {
	// 1. 验证邮箱验证码（委托给认证领域服务）
	if err := s.authDomainService.VerifyEmailCode(ctx, req.Email, req.Captcha, "register"); err != nil {
		return nil, err
	}

	// 2. 注册用户（委托给用户领域服务）
	registeredUser, err := s.userDomainService.RegisterUser(ctx, "", req.Email, req.Password, req.Phone)
	if err != nil {
		return nil, err
	}

	// 3. 保存用户到数据库
	if err := s.userRepo.Create(ctx, registeredUser); err != nil {
		return nil, err
	}

	// 4. 生成令牌（委托给认证领域服务）
	tokens, err := s.authDomainService.GenerateUserTokens(ctx, registeredUser.ID)
	if err != nil {
		return nil, err
	}

	return &RegisterResponse{
		UserID:       registeredUser.ID,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    tokens.AccessExpiresIn,
	}, nil
}

// OTPLoginRequest OTP登录请求
type OTPLoginRequest struct {
	Account string
	Captcha string
}

// OTPLogin OTP登录（应用服务编排）
func (s *Service) OTPLogin(ctx context.Context, req *OTPLoginRequest) (*LoginResponse, error) {
	// 1. 验证验证码（委托给认证领域服务）
	if err := s.authDomainService.VerifyEmailCode(ctx, req.Account, req.Captcha, "login"); err != nil {
		return nil, err
	}

	// 2. 查找并验证用户
	u, err := s.userRepo.FindByAccount(ctx, req.Account)
	if err != nil {
		return nil, err
	}
	if u == nil || !u.IsActive() {
		return nil, user.ErrUserNotFound
	}

	// 3. 生成令牌（委托给认证领域服务）
	tokens, err := s.authDomainService.GenerateUserTokens(ctx, u.ID)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		UserID:       u.ID,
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
		ExpiresIn:    tokens.AccessExpiresIn,
	}, nil
}

// RetrievePasswordRequest 找回密码请求
type RetrievePasswordRequest struct {
	Email    string
	Captcha  string
	Password string
}

// RetrievePassword 找回密码（应用服务编排）
func (s *Service) RetrievePassword(ctx context.Context, req *RetrievePasswordRequest) error {
	// 1. 验证邮箱验证码（委托给认证领域服务）
	if err := s.authDomainService.VerifyEmailCode(ctx, req.Email, req.Captcha, "reset_password"); err != nil {
		return err
	}

	// 2. 重置密码（委托给用户领域服务）
	return s.userDomainService.ResetUserPassword(ctx, req.Email, req.Password)
}

// ChangeEmailRequest 更换邮箱请求
type ChangeEmailRequest struct {
	Email   string
	Captcha string
}

// ChangeEmail 更换邮箱（应用服务编排）
func (s *Service) ChangeEmail(ctx context.Context, userID string, req *ChangeEmailRequest) error {
	// 1. 验证邮箱验证码（委托给认证领域服务）
	if err := s.authDomainService.VerifyEmailCode(ctx, req.Email, req.Captcha, "change_email"); err != nil {
		return err
	}

	// 2. 更换邮箱（委托给用户领域服务）
	return s.userDomainService.ChangeUserEmail(ctx, userID, req.Email)
}

// ChangeEmailOTPRequest 更换邮箱获取验证码请求
type ChangeEmailOTPRequest struct {
	Email string
}

// ChangeEmailOTP 更换邮箱获取验证码（应用服务编排）
func (s *Service) ChangeEmailOTP(ctx context.Context, req *ChangeEmailOTPRequest) error {
	// 1. 检查邮箱是否已存在
	exists, err := s.userRepo.ExistsByEmail(ctx, req.Email)
	if err != nil {
		return err
	}
	if exists {
		return user.ErrEmailAlreadyExists
	}

	// 2. 发送验证码（委托给认证领域服务）
	return s.authDomainService.SendEmailVerificationCode(ctx, req.Email, "change_email")
}

// AuthenticateRPC RPC认证（应用服务编排）
func (s *Service) AuthenticateRPC(ctx context.Context, token string) (*auth.Claims, error) {
	return s.authDomainService.VerifyUserAccessToken(ctx, token)
}

// GenerateTokens 生成令牌（应用服务编排）
func (s *Service) GenerateTokens(ctx context.Context, userID string) (*auth.TokenPair, error) {
	return s.authDomainService.GenerateUserTokens(ctx, userID)
}

// RefreshTokenRequest 刷新令牌请求
type RefreshTokenRequest struct {
	RefreshToken string
}

// RefreshTokenResponse 刷新令牌响应
type RefreshTokenResponse struct {
	UserID           string
	AccessToken      string
	RefreshToken     string
	AccessExpiresIn  int32
	RefreshExpiresIn int32
}

// RefreshToken 刷新令牌（应用服务编排）
func (s *Service) RefreshToken(ctx context.Context, req *RefreshTokenRequest) (*RefreshTokenResponse, error) {
	// 使用刷新令牌生成新的令牌对
	tokenPair, err := s.authDomainService.RefreshUserTokens(ctx, req.RefreshToken)
	if err != nil {
		return nil, err
	}

	// 从新的访问令牌中解析用户ID
	claims, err := s.authDomainService.VerifyUserAccessToken(ctx, tokenPair.AccessToken)
	if err != nil {
		return nil, err
	}

	return &RefreshTokenResponse{
		UserID:           claims.UserID,
		AccessToken:      tokenPair.AccessToken,
		RefreshToken:     tokenPair.RefreshToken,
		AccessExpiresIn:  tokenPair.AccessExpiresIn,
		RefreshExpiresIn: tokenPair.RefreshExpiresIn,
	}, nil
}
