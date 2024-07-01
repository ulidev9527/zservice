package zauth

import (
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/zglobal"
)

func UserOrgBind(ctx *zservice.Context, in *zauth_pb.UserOrgBind_REQ) *zauth_pb.Default_RES {

	if in.Uid == 0 || in.OrgID == 0 {
		return &zauth_pb.Default_RES{Code: zglobal.Code_ParamsErr}
	}

	if res, e := grpcClient.UserOrgBind(ctx, in); e != nil {
		return &zauth_pb.Default_RES{Code: zglobal.Code_Fail}
	} else {
		return res
	}

}
