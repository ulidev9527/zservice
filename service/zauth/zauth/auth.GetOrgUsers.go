package zauth

import (
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
)

func GetOrgUsers(ctx *zservice.Context, in *zauth_pb.GetOrgUsers_REQ) *zauth_pb.GetOrgUsers_RES {

	if res, e := grpcClient.GetOrgUsers(ctx, in); e != nil {
		ctx.LogError(e)
		return &zauth_pb.GetOrgUsers_RES{Code: zservice.Code_Fail}
	} else {
		return res
	}
}
