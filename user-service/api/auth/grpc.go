package auth

import authpb "github.com/people257/poor-guy-shop/user-service/gen/proto/user/auth"

var _ authpb.AuthServiceServer = (*AuthServer)(nil)

type AuthServer struct {
}

func (s *AuthServer) Login(ctx context.Context, req *authpb.LoginReq) (*authpb.LoginResp, error) {
	//TODO implement me
	panic("implement me")
}

func (s *AuthServer) EmailRegister(ctx context.Context, req *authpb.EmailRegisterReq) (*authpb.RegisterResp, error) {
	//TODO implement me
	panic("implement me")
}

func (s *AuthServer) OTPLogin(ctx context.Context, req *authpb.OTPLoginReq) (*authpb.LoginResp, error) {
	//TODO implement me
	panic("implement me")
}

func (s *AuthServer) RetrievePassword(ctx context.Context, empty *emptypb.Empty) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (s *AuthServer) ChangeEmail(ctx context.Context, req *authpb.ChangeEmailReq) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (s *AuthServer) ChangeEmailOTP(ctx context.Context, req *authpb.ChangeEmailOTPReq) (*emptypb.Empty, error) {
	//TODO implement me
	panic("implement me")
}

func (s *AuthServer) AuthenticateRPC(ctx context.Context, empty *emptypb.Empty) (*authpb.AuthenticateRPCResp, error) {
	//TODO implement me
	panic("implement me")
}
