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
func HasAccountID(ctx *zservice.Context, in *zauth_pb.HasAccountID_REQ) *zauth_pb.Default_RES {

	if res, e := func() (*zauth_pb.Default_RES, error) {
		if zauthInitConfig.ZauthServiceName == zservice.GetServiceName() {
			return internal.Logic_HasAccountID(ctx, in), nil
		}
		return grpcClient.HasAccountID(context.WithValue(context.Background(), grpcservice.GRPC_contextEX_Middleware_Key, ctx.ContextS2S), in)
	}(); e != nil {
		return &zauth_pb.Default_RES{Code: zglobal.Code_ErrorBreakoff}
	} else {
		return res
	}
}
