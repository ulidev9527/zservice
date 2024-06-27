package internal

import (
	"strings"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/zglobal"
)

// 创建逻辑
func Logic_PermissionCreate(ctx *zservice.Context, in *zauth_pb.PermissionInfo) *zauth_pb.PermissionInfo_RES {

	if in == nil || in.Name == "" || in.Service == "" {
		ctx.LogError("param error")
		return &zauth_pb.PermissionInfo_RES{Code: zglobal.Code_ParamsErr}
	}

	if in.State > 3 {
		in.State = 3
	}

	// 检查权限是否存在
	if _, e := GetPermissionBySAP(ctx, in.Service, in.Action, in.Path); e != nil {
		if e.GetCode() != zglobal.Code_NotFound {
			ctx.LogError(e)
			return &zauth_pb.PermissionInfo_RES{Code: e.GetCode()}
		}
	} else {
		return &zauth_pb.PermissionInfo_RES{Code: zglobal.Code_Zauth_Permission_Alerady_Exist}
	}

	// 获取一个未使用的权限 ID
	pid, e := GetNewPermissionID(ctx)
	if e != nil {
		ctx.LogError(e)
		return &zauth_pb.PermissionInfo_RES{Code: e.GetCode()}
	}

	z := &PermissionTable{
		Name:         in.Name,
		PermissionID: pid,
		Service:      in.Service,
		Action:       strings.ToLower(in.Action),
		Path:         in.Path,
		State:        in.State,
	}

	if e := z.Save(ctx); e != nil {
		ctx.LogError(e)
		return &zauth_pb.PermissionInfo_RES{Code: e.GetCode()}
	}

	return &zauth_pb.PermissionInfo_RES{Code: zglobal.Code_SUCC, Info: &zauth_pb.PermissionInfo{
		PermissionID: z.PermissionID,
		Name:         z.Name,
		Service:      z.Service,
		Action:       z.Action,
		Path:         z.Path,
		State:        z.State,
	}}
}
