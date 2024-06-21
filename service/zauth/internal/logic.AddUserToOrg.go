package internal

import (
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/zglobal"
)

func Logic_AddUserToOrg(ctx *zservice.Context, in *zauth_pb.AddUserToOrg_REQ) *zauth_pb.Default_RES {
	if in.Uid == 0 || in.OrgID == 0 {
		return &zauth_pb.Default_RES{Code: zglobal.Code_ParamsErr}
	}

	if has, e := HasUserOrgBindByID(ctx, in.Uid, in.OrgID); e != nil {
		ctx.LogError(e)
		return &zauth_pb.Default_RES{Code: e.GetCode()}
	} else if has {
		return &zauth_pb.Default_RES{Code: zglobal.Code_RepetitionErr}
	}

	if _, e := UserOrgBind(ctx, in.Uid, in.OrgID, in.Expire); e != nil {
		ctx.LogError(e)
		return &zauth_pb.Default_RES{Code: e.GetCode()}
	}

	return &zauth_pb.Default_RES{Code: zglobal.Code_SUCC}

}
