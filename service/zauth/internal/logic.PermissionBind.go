package internal

import (
	"fmt"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/zglobal"
)

// 权限绑定
// 已有权限会进行更新
func Logic_PermissionBind(ctx *zservice.Context, in *zauth_pb.PermissionBind_REQ) *zauth_pb.Default_RES {

	// 验证参数是否正确
	switch in.TargetType {
	case 1:
		// 组织验证
		if has, e := HasOrgByID(ctx, in.TargetID); e != nil {
			ctx.LogError(e)
			return &zauth_pb.Default_RES{Code: e.GetCode()}
		} else if !has {
			return &zauth_pb.Default_RES{Code: zglobal.Code_Zauth_PermissionBind_TargetIDErr}
		}
	case 2:
		// 账号验证
		if has, e := HasAccountByID(ctx, in.TargetID); e != nil {
			ctx.LogError(e)
			return &zauth_pb.Default_RES{Code: e.GetCode()}
		} else if !has {
			return &zauth_pb.Default_RES{Code: zglobal.Code_Zauth_PermissionBind_TargetIDErr}
		}
	default:
		return &zauth_pb.Default_RES{Code: zglobal.Code_Zauth_PermissionBind_TargetTypeErr}
	}

	// 权限验证
	if has, e := HasPermissionByID(ctx, in.PermissionID); e != nil {
		ctx.LogError(e)
		return &zauth_pb.Default_RES{Code: e.GetCode()}
	} else if !has {
		return &zauth_pb.Default_RES{Code: zglobal.Code_Zauth_PermissionBind_PermissionIDErr}
	}

	// 获取是否有绑定
	if tab, e := GetPermissionBind(ctx, in.TargetType, in.TargetID, in.PermissionID); e != nil {
		if e.GetCode() != zglobal.Code_NotFound {
			ctx.LogError(e)
			return &zauth_pb.Default_RES{Code: e.GetCode()}
		}
	} else if tab != nil { // 有数据 更新

		if zservice.MD5String(fmt.Sprint(in.TargetType, in.TargetID, in.PermissionID, in.Expires, in.State)) ==
			zservice.MD5String(fmt.Sprint(tab.TargetType, tab.TargetID, tab.PermissionID, tab.Expires, tab.State)) {
			return &zauth_pb.Default_RES{Code: zglobal.Code_SUCC}
		} else {

			tab.TargetType = in.TargetType
			tab.TargetID = in.TargetID
			tab.PermissionID = in.PermissionID
			tab.Expires = in.Expires
			tab.State = in.State
			if e := tab.Save(ctx); e != nil {
				ctx.LogError(e)
				return &zauth_pb.Default_RES{Code: e.GetCode()}
			}

			return &zauth_pb.Default_RES{Code: zglobal.Code_SUCC}
		}
	}
	// 创建新数据
	tab := &ZauthPermissionBindTable{
		TargetType:   in.TargetType,
		TargetID:     in.TargetID,
		PermissionID: in.PermissionID,
		Expires:      in.Expires,
		State:        in.State,
	}

	if e := tab.Save(ctx); e != nil {
		ctx.LogError(e)
		return &zauth_pb.Default_RES{Code: e.GetCode()}
	}

	return &zauth_pb.Default_RES{Code: zglobal.Code_SUCC}
}
