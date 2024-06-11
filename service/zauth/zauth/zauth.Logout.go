package zauth

import (
	"context"
	"zservice/service/zauth/internal"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/ex/grpcservice"
	"zservice/zservice/zglobal"
)

func Logout(ctx *zservice.Context) *zauth_pb.Default_RES {

	if res, e := func() (*zauth_pb.Default_RES, error) {
		in := &zauth_pb.Default_REQ{}
		if zauthInitConfig.ZauthServiceName == zservice.GetServiceName() {
			return internal.Logic_Logout(ctx, in), nil
		}
		return grpcClient.Logout(context.WithValue(context.Background(), grpcservice.GRPC_contextEX_Middleware_Key, ctx.ContextS2S), in)
	}(); e != nil {
		ctx.LogPanic(e)
	} else if res.Code != zglobal.Code_SUCC {
		ctx.LogPanic(res)
	}
	return &zauth_pb.Default_RES{Code: zglobal.Code_SUCC}
}
