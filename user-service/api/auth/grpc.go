package auth

import (
	"context"
	"time"

	authpb "github.com/people257/poor-guy-shop/user-service/gen/proto/user/auth"
	"github.com/people257/poor-guy-shop/user-service/internal/application/auth"
	infraAuth "github.com/people257/poor-guy-shop/user-service/internal/infra/auth"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

var _ authpb.AuthServiceServer = (*AuthServer)(nil)

type AuthServer struct {
	authService *auth.Service
	authInfra   *infraAuth.Auth
}

// NewAuthServer 创建认证服务器
func NewAuthServer(authService *auth.Service, authInfra *infraAuth.Auth) *AuthServer {
	return &AuthServer{
		authService: authService,
		authInfra:   authInfra,
	}
}

func (s *AuthServer) Login(ctx context.Context, req *authpb.LoginReq) (*authpb.LoginResp, error) {
	// 验证请求参数
	if req.Account == "" {
		return nil, status.Error(codes.InvalidArgument, "账号不能为空")
	}
	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "密码不能为空")
	}

	// 调用应用服务
	loginReq := &auth.LoginRequest{
		Account:  req.Account,
		Password: req.Password,
	}

	resp, err := s.authService.Login(ctx, loginReq)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return &authpb.LoginResp{
		UserId:       resp.UserID,
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		ExpiresIn:    resp.ExpiresIn,
	}, nil
}

func (s *AuthServer) EmailRegister(ctx context.Context, req *authpb.EmailRegisterReq) (*authpb.RegisterResp, error) {
	// 验证请求参数
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "邮箱不能为空")
	}
	if req.Password == "" {
		return nil, status.Error(codes.InvalidArgument, "密码不能为空")
	}
	if req.Captcha == "" {
		return nil, status.Error(codes.InvalidArgument, "验证码不能为空")
	}

	// 调用应用服务
	registerReq := &auth.EmailRegisterRequest{
		Email:    req.Email,
		Password: req.Password,
		Captcha:  req.Captcha,
	}
	if req.Phone != "" {
		registerReq.Phone = &req.Phone
	}

	resp, err := s.authService.EmailRegister(ctx, registerReq)
	if err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &authpb.RegisterResp{
		UserId:       resp.UserID,
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		ExpiresIn:    resp.ExpiresIn,
	}, nil
}

func (s *AuthServer) OTPLogin(ctx context.Context, req *authpb.OTPLoginReq) (*authpb.LoginResp, error) {
	// 验证请求参数
	if req.Account == "" {
		return nil, status.Error(codes.InvalidArgument, "账号不能为空")
	}
	if req.Captcha == "" {
		return nil, status.Error(codes.InvalidArgument, "验证码不能为空")
	}

	// 调用应用服务
	otpReq := &auth.OTPLoginRequest{
		Account: req.Account,
		Captcha: req.Captcha,
	}

	resp, err := s.authService.OTPLogin(ctx, otpReq)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return &authpb.LoginResp{
		UserId:       resp.UserID,
		AccessToken:  resp.AccessToken,
		RefreshToken: resp.RefreshToken,
		ExpiresIn:    resp.ExpiresIn,
	}, nil
}

func (s *AuthServer) RetrievePassword(ctx context.Context, empty *emptypb.Empty) (*emptypb.Empty, error) {
	// TODO: 找回密码功能需要从HTTP请求中获取参数
	// 由于gRPC接口使用Empty，实际参数应该通过HTTP body传递
	return nil, status.Error(codes.Unimplemented, "找回密码功能待实现")
}

func (s *AuthServer) ChangeEmail(ctx context.Context, req *authpb.ChangeEmailReq) (*emptypb.Empty, error) {
	// 从Token中获取用户ID
	claims, err := s.authInfra.VerifyFromMetadata(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "认证失败")
	}

	// 验证请求参数
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "邮箱不能为空")
	}
	if req.Captcha == "" {
		return nil, status.Error(codes.InvalidArgument, "验证码不能为空")
	}

	// 调用应用服务
	changeReq := &auth.ChangeEmailRequest{
		Email:   req.Email,
		Captcha: req.Captcha,
	}

	if err := s.authService.ChangeEmail(ctx, claims.UserID(), changeReq); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *AuthServer) ChangeEmailOTP(ctx context.Context, req *authpb.ChangeEmailOTPReq) (*emptypb.Empty, error) {
	// 验证请求参数
	if req.Email == "" {
		return nil, status.Error(codes.InvalidArgument, "邮箱不能为空")
	}

	// 调用应用服务
	otpReq := &auth.ChangeEmailOTPRequest{
		Email: req.Email,
	}

	if err := s.authService.ChangeEmailOTP(ctx, otpReq); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *AuthServer) AuthenticateRPC(ctx context.Context, empty *emptypb.Empty) (*authpb.AuthenticateRPCResp, error) {
	// 从Metadata中获取Token
	claims, err := s.authInfra.VerifyFromMetadata(ctx)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, "认证失败")
	}

	// 检查是否需要刷新Token
	var newAccessToken string
	var expiresIn int64

	// 如果Token快要过期，生成新Token
	var expiresAtTime time.Time
	if claims.ExpiresAt != nil {
		expiresAtTime = claims.ExpiresAt.Time
	}

	if time.Until(expiresAtTime) < s.authInfra.RefreshThreshold() {
		// 生成新的令牌
		tokens, err := s.authService.GenerateTokens(ctx, claims.UserID())
		if err != nil {
			return nil, status.Error(codes.Internal, "生成新令牌失败")
		}
		newAccessToken = tokens.AccessToken
		expiresIn = time.Now().Add(time.Duration(tokens.AccessExpiresIn) * time.Second).UnixMilli()
	} else {
		expiresIn = expiresAtTime.UnixMilli()
	}

	return &authpb.AuthenticateRPCResp{
		UserId:      claims.UserID(),
		AccessToken: newAccessToken,
		ExpiresIn:   expiresIn,
	}, nil
}

// RefreshToken 刷新令牌
func (s *AuthServer) RefreshToken(ctx context.Context, req *authpb.RefreshTokenReq) (*authpb.RefreshTokenResp, error) {
	// 验证请求参数
	if req.RefreshToken == "" {
		return nil, status.Error(codes.InvalidArgument, "刷新令牌不能为空")
	}

	// 调用应用服务
	refreshReq := &auth.RefreshTokenRequest{
		RefreshToken: req.RefreshToken,
	}

	resp, err := s.authService.RefreshToken(ctx, refreshReq)
	if err != nil {
		return nil, status.Error(codes.Unauthenticated, err.Error())
	}

	return &authpb.RefreshTokenResp{
		UserId:           resp.UserID,
		AccessToken:      resp.AccessToken,
		RefreshToken:     resp.RefreshToken,
		AccessExpiresIn:  resp.AccessExpiresIn,
		RefreshExpiresIn: resp.RefreshExpiresIn,
	}, nil
}
