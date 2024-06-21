package internal

import (
	"context"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
)

type ZauthServer struct {
	zauth_pb.UnimplementedZauthServer
}

// 登出
func (s *ZauthServer) Logout(ctx context.Context, in *zauth_pb.Default_REQ) (*zauth_pb.Default_RES, error) {
	return Logic_Logout(ctx.(*zservice.Context), in), nil
}

// 手机号登陆
func (s *ZauthServer) LoginByPhone(ctx context.Context, in *zauth_pb.LoginByPhone_REQ) (*zauth_pb.Login_RES, error) {
	return Logic_LoginByPhone(ctx.(*zservice.Context), in), nil
}

// 账号密码登陆
func (s *ZauthServer) LoginByUser(ctx context.Context, in *zauth_pb.LoginByUser_REQ) (*zauth_pb.Login_RES, error) {
	return Logic_LoginByName(ctx.(*zservice.Context), in), nil
}

// 是否有账号
func (s *ZauthServer) HasUID(ctx context.Context, in *zauth_pb.HasUID_REQ) (*zauth_pb.Default_RES, error) {
	return Logic_HasUID(ctx.(*zservice.Context), in), nil
}

// 登陆检查
func (s *ZauthServer) LoginCheck(ctx context.Context, in *zauth_pb.Default_REQ) (*zauth_pb.Default_RES, error) {
	return Logic_LoginCheck(ctx.(*zservice.Context), in), nil
}

// 检查权限
func (s *ZauthServer) CheckAuth(ctx context.Context, in *zauth_pb.CheckAuth_REQ) (*zauth_pb.CheckAuth_RES, error) {
	return Logic_CheckAuth(ctx.(*zservice.Context), in), nil
}

// 短信验证码发送
func (s *ZauthServer) SMSVerifyCodeSend(ctx context.Context, in *zauth_pb.SMSVerifyCodeSend_REQ) (*zauth_pb.SMSSendVerifyCode_RES, error) {
	return Logic_SMSVerifyCodeSend(ctx.(*zservice.Context), in), nil
}

// 短信验证码校验
func (s *ZauthServer) SMSVerifyCodeVerify(ctx context.Context, in *zauth_pb.SMSVerifyCodeVerify_REQ) (*zauth_pb.Default_RES, error) {
	return Logic_SMSVerifyCodeVerify(ctx.(*zservice.Context), in), nil
}

// ZZZZ字符串验证
func (s *ZauthServer) HasZZZZString(ctx context.Context, in *zauth_pb.HasZZZZString_REQ) (*zauth_pb.Default_RES, error) {
	return Logic_HasZZZZString(ctx.(*zservice.Context), in), nil
}

// 注册服务
func (s *ZauthServer) ServiceRegist(ctx context.Context, in *zauth_pb.ServiceRegist_REQ) (*zauth_pb.ServiceRegist_RES, error) {
	return Logic_ServiceRegist(ctx.(*zservice.Context), in), nil
}
func (s *ZauthServer) AddUserToOrg(ctx context.Context, in *zauth_pb.AddUserToOrg_REQ) (*zauth_pb.Default_RES, error) {
	return Logic_AddUserToOrg(ctx.(*zservice.Context), in), nil
}

// 获取文件配置
func (s *ZauthServer) ConfigGetFileConfig(ctx context.Context, in *zauth_pb.ConfigGetFileConfig_REQ) (*zauth_pb.ConfigGetFileConfig_RES, error) {

	return Logic_ConfigGetFileConfig(ctx.(*zservice.Context), in), nil
}
