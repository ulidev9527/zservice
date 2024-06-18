package zauth

import (
	"context"
	"zservice/service/zauth/internal"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/ex/grpcservice"
	"zservice/zservice/zglobal"
)

// 是否有这个账号ID
func HasUID(ctx *zservice.Context, in *zauth_pb.HasUID_REQ) *zauth_pb.Default_RES {

	if res, e := func() (*zauth_pb.Default_RES, error) {
		if zauthInitConfig.ServiceName == zservice.GetServiceName() {
			return internal.Logic_HasUID(ctx, in), nil
		}
		return grpcClient.HasUID(context.WithValue(context.Background(), grpcservice.GRPC_contextEX_Middleware_Key, ctx.ContextS2S), in)
	}(); e != nil {
		return &zauth_pb.Default_RES{Code: zglobal.Code_ErrorBreakoff}
	} else {
		return res
	}
}
