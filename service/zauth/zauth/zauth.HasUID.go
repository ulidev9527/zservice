package zauth

import (
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/zglobal"
)

// 是否有这个账号ID
func HasUID(ctx *zservice.Context, in *zauth_pb.HasUID_REQ) *zauth_pb.Default_RES {

	if res, e := grpcClient.HasUID(ctx, in); e != nil {
		return &zauth_pb.Default_RES{Code: zglobal.Code_Fail}
	} else {
		return res
	}
}
