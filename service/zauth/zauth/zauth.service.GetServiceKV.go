package zauth

import (
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/zglobal"
)

func GetServiceKV(ctx *zservice.Context, in *zauth_pb.GetServiceKV_REQ) *zauth_pb.GetServiceKV_RES {
	if res, e := grpcClient.GetServiceKV(ctx, in); e != nil {
		return &zauth_pb.GetServiceKV_RES{Code: zglobal.Code_Fail}
	} else {
		return res
	}
}
