package zauth

import (
	"context"
	"zservice/service/zauth/internal"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/ex/grpcservice"
	"zservice/zservice/zglobal"
)

// 登陆检查
func LoginCheck(ctx *zservice.Context) bool {
	if ctx.AuthToken == "" {
		return false
	}

	if res, e := func() (*zauth_pb.Default_RES, error) {
		in := &zauth_pb.Default_REQ{}
		if zauthInitConfig.ServiceName == zservice.GetServiceName() {
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
