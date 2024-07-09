package zauth

import (
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
)

// 注册到 zauth 服务
func ServiceRegist(ctx *zservice.Context, in *zauth_pb.ServiceRegist_REQ) {

	in.Service = zservice.GetServiceName()

	if res, e := grpcClient.ServiceRegist(ctx, in); e != nil {
		ctx.LogPanic(e)
	} else if res.Code != zservice.Code_SUCC {
		ctx.LogPanic(res)
	} else {
		serviceInfo.serviceRegistRES = res
	}
}
