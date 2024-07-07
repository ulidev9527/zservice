package zauth

import (
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
)

func SetServiceKV(ctx *zservice.Context, in *zauth_pb.SetServiceKV_REQ) *zauth_pb.Default_RES {
	if res, e := grpcClient.SetServiceKV(ctx, in); e != nil {
		return &zauth_pb.Default_RES{Code: zservice.Code_Fail}
	} else {
		return res
	}
}
