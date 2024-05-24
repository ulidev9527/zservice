package internal

import (
	"context"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice/ex/grpcservice"
)

type ZauthServer struct {
	zauth_pb.UnimplementedZauthServer
}

// 登出
func (s *ZauthServer) Logout(ctx context.Context, in *zauth_pb.Default_REQ) (*zauth_pb.Default_RES, error) {
	return Logic_Logout(grpcservice.GetCtxEX(ctx), in), nil
}

// 手机号登陆
func (s *ZauthServer) LoginByPhone(ctx context.Context, in *zauth_pb.LoginByPhone_REQ) (*zauth_pb.Default_RES, error) {
	return Logic_LoginByPhone(grpcservice.GetCtxEX(ctx), in), nil
}

// 检查权限
func (s *ZauthServer) CheckAuth(ctx context.Context, in *zauth_pb.CheckAuth_REQ) (*zauth_pb.CheckAuth_RES, error) {
	return Logic_CheckAuth(grpcservice.GetCtxEX(ctx), in), nil
}

// 短信验证码发送
func (s *ZauthServer) SMSSendVerifyCode(ctx context.Context, in *zauth_pb.SMSSendVerifyCode_REQ) (*zauth_pb.SMSSendVerifyCode_RES, error) {
	return Logic_SMSSendVerifyCode(grpcservice.GetCtxEX(ctx), in), nil
}

// 短信验证码校验
func (s *ZauthServer) SMSVerifyCode(ctx context.Context, in *zauth_pb.SMSVerifyCode_REQ) (*zauth_pb.Default_RES, error) {
	return Logic_SMSVerifyCode(grpcservice.GetCtxEX(ctx), in), nil
}

// 获取文件配置
func (s *ZauthServer) GetFileConfig(ctx context.Context, in *zauth_pb.GetFileConfig_REQ) (*zauth_pb.GetFileConfig_RES, error) {
	return Logic_GetFileConfig(grpcservice.GetCtxEX(ctx), in), nil
}
