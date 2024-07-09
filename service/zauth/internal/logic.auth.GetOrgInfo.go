package internal

import (
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
)

func Logic_GetOrgInfo(ctx *zservice.Context, in *zauth_pb.GetOrgInfo_REQ) *zauth_pb.OrgInfo_RES {
	if in.OrgID != 0 {
		tab, e := GetOrgByID(ctx, in.OrgID)
		if e != nil {
			ctx.LogError(e)
			return &zauth_pb.OrgInfo_RES{Code: e.GetCode()}
		}
		return &zauth_pb.OrgInfo_RES{Code: zservice.Code_SUCC, Info: tab.ToOrgInfo()}
	} else if in.Service != "" && in.OrgName != "" {
		tab, e := GetOrgByName(ctx, in.Service, in.OrgName)
		if e != nil {
			ctx.LogError(e)
			return &zauth_pb.OrgInfo_RES{Code: e.GetCode()}
		}
		return &zauth_pb.OrgInfo_RES{Code: zservice.Code_SUCC, Info: tab.ToOrgInfo()}
	}
	return &zauth_pb.OrgInfo_RES{Code: zservice.Code_ParamsErr}
}
