package zauth

import (
	"context"
	"zservice/service/zauth/internal"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/ex/grpcservice"
	"zservice/zservice/zglobal"
)

func HasZZZZString(ctx *zservice.Context, str string) bool {
	if res, e := func() (*zauth_pb.Default_RES, error) {
		in := &zauth_pb.HasZZZZString_REQ{Str: str}
		if zauthInitConfig.ZauthServiceName == zservice.GetServiceName() {
			return internal.Logic_HasZZZZString(ctx, in), nil
		}
		return grpcClient.HasZZZZString(context.WithValue(context.Background(), grpcservice.GRPC_contextEX_Middleware_Key, ctx.ContextS2S), in)
	}(); e != nil {
		ctx.LogError(e)
		return true
	} else if res.Code != zglobal.Code_SUCC {
		return false
	} else {
		return true
	}
}
