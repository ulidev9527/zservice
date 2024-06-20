package zauth

import (
	"zservice/service/zauth/internal"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/zglobal"
)

func Logout(ctx *zservice.Context) *zauth_pb.Default_RES {

	if res, e := func() (*zauth_pb.Default_RES, error) {
		in := &zauth_pb.Default_REQ{}
		if zauthInitConfig.ServiceName == zservice.GetServiceName() {
			return internal.Logic_Logout(ctx, in), nil
		}
		return grpcClient.Logout(ctx, in)
	}(); e != nil {
		ctx.LogPanic(e)
	} else if res.Code != zglobal.Code_SUCC {
		ctx.LogPanic(res)
	}
	return &zauth_pb.Default_RES{Code: zglobal.Code_SUCC}
}
