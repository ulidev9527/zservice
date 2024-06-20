package zauth

import (
	"zservice/service/zauth/internal"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/zglobal"
)

// 是否有这个账号ID
func HasUID(ctx *zservice.Context, in *zauth_pb.HasUID_REQ) *zauth_pb.Default_RES {

	if res, e := func() (*zauth_pb.Default_RES, error) {
		if zauthInitConfig.ServiceName == zservice.GetServiceName() {
			return internal.Logic_HasUID(ctx, in), nil
		}
		return grpcClient.HasUID(ctx, in)
	}(); e != nil {
		return &zauth_pb.Default_RES{Code: zglobal.Code_Fail}
	} else {
		return res
	}
}
