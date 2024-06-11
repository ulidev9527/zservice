package zauth

import (
	"context"
	"zservice/service/zauth/internal"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/ex/grpcservice"
	"zservice/zservice/zglobal"
)

// 账号登陆
func LoginByAccount(ctx *zservice.Context, in *zauth_pb.LoginByAccount_REQ) *zauth_pb.Login_RES {

	if in.Account == "" || in.Password == "" {
		return &zauth_pb.Login_RES{Code: zglobal.Code_ParamsErr}
	}

	if res, e := func() (*zauth_pb.Login_RES, error) {
		if zauthInitConfig.ZauthServiceName == zservice.GetServiceName() {
			return internal.Logic_LoginByAccount(ctx, in), nil
		}
		return grpcClient.LoginByAccount(context.WithValue(context.Background(), grpcservice.GRPC_contextEX_Middleware_Key, ctx.ContextS2S), in)
	}(); e != nil {
		return &zauth_pb.Login_RES{Code: zglobal.Code_ErrorBreakoff}
	} else {
		return res
	}
}
