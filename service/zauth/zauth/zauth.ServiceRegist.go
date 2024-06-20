package zauth

import (
	"zservice/service/zauth/internal"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/zglobal"
)

// 服务注册
func ServiceRegist(ctx *zservice.Context) {

	if res, e := func() (*zauth_pb.Default_RES, error) {
		in := &zauth_pb.Default_REQ{}
		if zauthInitConfig.ServiceName == zservice.GetServiceName() {
			return internal.Logic_ServiceRegist(ctx, in), nil
		}
		return grpcClient.ServiceRegist(ctx, in)
	}(); e != nil {
		ctx.LogPanic(e)
	} else if res.Code != zglobal.Code_SUCC {
		ctx.LogPanic(res)
	}
}
