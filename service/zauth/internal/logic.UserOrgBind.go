package internal

import (
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/zglobal"
)

func Logic_UserOrgBind(ctx *zservice.Context, in *zauth_pb.UserOrgBind_REQ) *zauth_pb.Default_RES {
	if in.Uid == 0 || in.OrgID == 0 {
		ctx.LogError("params error", in)
		return &zauth_pb.Default_RES{Code: zglobal.Code_ParamsErr}
	}

	if _, e := UserOrgBind(ctx, in.Uid, in.OrgID, in.Expires, in.State); e != nil {
		ctx.LogError(e)
		return &zauth_pb.Default_RES{Code: e.GetCode()}
	}

	return &zauth_pb.Default_RES{Code: zglobal.Code_SUCC}
}
