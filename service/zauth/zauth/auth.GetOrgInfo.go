package zauth

import (
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
)

func GetOrgInfo(ctx *zservice.Context, in *zauth_pb.GetOrgInfo_REQ) *zauth_pb.OrgInfo_RES {
	if res, e := grpcClient.GetOrgInfo(ctx, in); e != nil {
		ctx.LogError(e)
		return &zauth_pb.OrgInfo_RES{Code: zservice.Code_Fail}
	} else {
		return res
	}
}
