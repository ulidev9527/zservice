package zauth

import (
	"context"
	"zservice/service/zauth/internal"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/ex/grpcservice"
	"zservice/zservice/zglobal"
)

// 手机号登陆
func LoginByPhone(ctx *zservice.Context, in *zauth_pb.LoginByPhone_REQ) *zauth_pb.Default_RES {

	if in.Phone == "" || in.VerifyCode == "" {
		return &zauth_pb.Default_RES{Code: zglobal.Code_ParamsErr}
	}

	if res, e := func() (*zauth_pb.Default_RES, error) {
		if zauthInitConfig.ZauthServiceName == zservice.GetServiceName() {
			return internal.Logic_LoginByPhone(ctx, in), nil
		}
		return grpcClient.LoginByPhone(context.WithValue(context.Background(), grpcservice.GRPC_contextEX_Middleware_Key, ctx.ContextS2S), in)
	}(); e != nil {
		return &zauth_pb.Default_RES{Code: zglobal.Code_ErrorBreakoff}
	} else {
		return res
	}

}

// 账号登陆
func LoginByAccount(ctx *zservice.Context, in *zauth_pb.LoginByAccount_REQ) *zauth_pb.Default_RES {

	if in.Account == "" || in.Password == "" {
		return &zauth_pb.Default_RES{Code: zglobal.Code_ParamsErr}
	}

	if res, e := func() (*zauth_pb.Default_RES, error) {
		if zauthInitConfig.ZauthServiceName == zservice.GetServiceName() {
			return internal.Logic_LoginByAccount(ctx, in), nil
		}
		return grpcClient.LoginByAccount(context.WithValue(context.Background(), grpcservice.GRPC_contextEX_Middleware_Key, ctx.ContextS2S), in)
	}(); e != nil {
		return &zauth_pb.Default_RES{Code: zglobal.Code_ErrorBreakoff}
	} else {
		return res
	}
}

// 登陆检查
func LoginCheck(ctx *zservice.Context) bool {
	if ctx.AuthToken == "" {
		return false
	}

	if res, e := func() (*zauth_pb.Default_RES, error) {
		in := &zauth_pb.Default_REQ{}
		if zauthInitConfig.ZauthServiceName == zservice.GetServiceName() {
			return internal.Logic_LoginCheck(ctx, in), nil
		} else {
			return grpcClient.LoginCheck(context.WithValue(context.Background(), grpcservice.GRPC_contextEX_Middleware_Key, ctx.ContextS2S), in)
		}
	}(); e != nil {
		ctx.LogError(e)
		return false
	} else if res.Code == zglobal.Code_SUCC {
		return true
	} else {
		return false
	}
}
