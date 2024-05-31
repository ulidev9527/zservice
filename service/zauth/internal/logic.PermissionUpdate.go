package internal

import (
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/zglobal"
)

func Logic_PermissionUpdate(ctx *zservice.Context, in *zauth_pb.PermissionInfo) *zauth_pb.Default_RES {

	if in.Id == 0 || in.Name == "" || in.Service == "" {
		return &zauth_pb.Default_RES{Code: zglobal.Code_ParamsErr}
	}

	if in.State > 3 {
		in.State = 3
	}

	// 是否有同名权限
	if tab, e := GetPermissionBySAP(ctx, in.Service, in.Action, in.Path); e == nil {
		if e.GetCode() != zglobal.Code_NotFound {
			ctx.LogError(e)
			return &zauth_pb.Default_RES{Code: e.GetCode()}
		}
	} else if tab != nil {
		return &zauth_pb.Default_RES{Code: zglobal.Code_Zauth_Permission_Alerady_Exist}
	}

	// 检查权限是否存在
	if tab, e := GetPermissionByID(ctx, uint(in.Id)); e != nil {
		ctx.LogError(e)
		return &zauth_pb.Default_RES{Code: e.GetCode()}
	} else {
		tab.Name = in.Name
		tab.Service = in.Service
		tab.Action = in.Action
		tab.Path = in.Path
		tab.State = in.State
		if e := tab.Save(ctx); e != nil {
			ctx.LogError(e)
			return &zauth_pb.Default_RES{Code: e.GetCode()}
		}
		return &zauth_pb.Default_RES{Code: zglobal.Code_SUCC}
	}

}
