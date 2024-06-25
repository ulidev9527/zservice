package zauth

import (
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/zglobal"
)

func CheckAuth(ctx *zservice.Context, in *zauth_pb.CheckAuth_REQ) *zauth_pb.CheckAuth_RES {
	if res, e := grpcClient.CheckAuth(ctx, in); e != nil {
		return &zauth_pb.CheckAuth_RES{Code: zglobal.Code_Fail}
	} else {
		return res
	}
}
