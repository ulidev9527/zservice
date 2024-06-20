package zauth

import (
	"zservice/service/zauth/internal"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/zglobal"
)

func HasZZZZString(ctx *zservice.Context, str string) bool {
	if res, e := func() (*zauth_pb.Default_RES, error) {
		in := &zauth_pb.HasZZZZString_REQ{Str: str}
		if zauthInitConfig.ServiceName == zservice.GetServiceName() {
			return internal.Logic_HasZZZZString(ctx, in), nil
		}
		return grpcClient.HasZZZZString(ctx, in)
	}(); e != nil {
		ctx.LogError(e)
		return true
	} else if res.Code != zglobal.Code_SUCC {
		return false
	} else {
		return true
	}
}
