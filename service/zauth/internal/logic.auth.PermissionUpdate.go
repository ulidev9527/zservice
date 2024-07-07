package internal

import (
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
)

func Logic_PermissionUpdate(ctx *zservice.Context, in *zauth_pb.PermissionInfo) *zauth_pb.Default_RES {

	if in.PermissionID == 0 || in.Name == "" || in.Service == "" {
		ctx.LogError("param error")
		return &zauth_pb.Default_RES{Code: zservice.Code_ParamsErr}
	}

	if in.State > 3 {
		in.State = 3
	}

	// 是否有同名权限
	if tab, e := GetPermissionBySAP(ctx, in.Service, in.Action, in.Path); e != nil {
		if e.GetCode() != zservice.Code_NotFound {
			ctx.LogError(e)
			return &zauth_pb.Default_RES{Code: e.GetCode()}
		}
	} else if tab != nil {

		if tab.PermissionID != in.PermissionID {
			return &zauth_pb.Default_RES{Code: zservice.Code_Zauth_Permission_Alerady_Exist}
		}
	}

	// 检查权限是否存在
	if tab, e := GetPermissionByID(ctx, uint(in.PermissionID)); e != nil {
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
		return &zauth_pb.Default_RES{Code: zservice.Code_SUCC}
	}

}
