package internal

import (
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/zglobal"
)

// 组织更新
func Logic_OrgUpdate(ctx *zservice.Context, in *zauth_pb.OrgInfo) *zauth_pb.Default_RES {
	if in.Name == "" {
		return &zauth_pb.Default_RES{Code: zglobal.Code_ParamsErr}
	}

	// 验证组织是否存在
	if tab, e := GetOrgByID(ctx, uint(in.OrgID)); e != nil {
		ctx.LogError(e)
		return &zauth_pb.Default_RES{Code: e.GetCode()}
	} else if tab == nil {
		return &zauth_pb.Default_RES{Code: zglobal.Code_Zauth_Org_NotFund}
	} else {
		if tab.Name == in.Name && tab.State == in.State {
			return &zauth_pb.Default_RES{Code: zglobal.Code_SUCC}
		} else {
			tab.Name = in.Name
			tab.State = in.State
			if e := tab.Save(ctx); e != nil {
				ctx.LogError(e)
				return &zauth_pb.Default_RES{Code: e.GetCode()}
			}
			return &zauth_pb.Default_RES{Code: zglobal.Code_SUCC}
		}
	}
}
