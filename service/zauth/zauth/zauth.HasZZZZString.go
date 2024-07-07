package zauth

import (
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
)

func HasZZZZString(ctx *zservice.Context, str string) bool {
	if res, e := grpcClient.HasZZZZString(ctx, &zauth_pb.HasZZZZString_REQ{Str: str}); e != nil {
		ctx.LogError(e)
		return true
	} else if res.Code != zservice.Code_SUCC {
		return false
	} else {
		return true
	}
}
