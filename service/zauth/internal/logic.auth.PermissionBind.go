package internal

import (
	"fmt"
	"time"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
)

// 权限绑定
// 已有权限会进行更新
func Logic_PermissionBind(ctx *zservice.Context, in *zauth_pb.PermissionBind_REQ) *zauth_pb.PermissionBind_RES {

	// 验证参数是否正确
	switch in.TargetType {
	case 1:
		// 组织验证
		if has, e := HasOrgByID(ctx, in.TargetID); e != nil {
			ctx.LogError(e)
			return &zauth_pb.PermissionBind_RES{Code: e.GetCode()}
		} else if !has {
			return &zauth_pb.PermissionBind_RES{Code: zservice.Code_Zauth_PermissionBind_TargetIDErr}
		}
	case 2:
		// 账号验证
		if has, e := HasUserByID(ctx, in.TargetID); e != nil {
			ctx.LogError(e)
			return &zauth_pb.PermissionBind_RES{Code: e.GetCode()}
		} else if !has {
			return &zauth_pb.PermissionBind_RES{Code: zservice.Code_Zauth_PermissionBind_TargetIDErr}
		}
	default:
		return &zauth_pb.PermissionBind_RES{Code: zservice.Code_Zauth_PermissionBind_TargetTypeErr}
	}

	// 权限验证
	if has, e := HasPermissionByID(ctx, in.PermissionID); e != nil {
		ctx.LogError(e)
		return &zauth_pb.PermissionBind_RES{Code: e.GetCode()}
	} else if !has {
		return &zauth_pb.PermissionBind_RES{Code: zservice.Code_Zauth_PermissionBind_PermissionIDErr}
	}

	// 获取是否有绑定
	if tab, e := GetPermissionBind(ctx, in.TargetType, in.TargetID, in.PermissionID); e != nil {
		if e.GetCode() != zservice.Code_NotFound {
			ctx.LogError(e)
			return &zauth_pb.PermissionBind_RES{Code: e.GetCode()}
		}
	} else if tab != nil { // 有数据 进行更新

		// 检查是否更新
		if zservice.MD5String(fmt.Sprint(in.TargetType, in.TargetID, in.PermissionID, in.Expires, in.State)) ==
			zservice.MD5String(fmt.Sprint(tab.TargetType, tab.TargetID, tab.PermissionID, tab.Expires, tab.State)) {
			return &zauth_pb.PermissionBind_RES{
				Code: zservice.Code_SUCC,
				Info: &zauth_pb.PermissionBindInfo{
					TargetType:   tab.TargetType,
					TargetID:     tab.TargetID,
					PermissionID: tab.PermissionID,
					Expires:      tab.Expires.UnixMilli(),
					State:        tab.State,
				},
			}
		} else {

			tab.TargetType = in.TargetType
			tab.TargetID = in.TargetID
			tab.PermissionID = in.PermissionID
			tab.Expires = zservice.NewTime(time.UnixMilli(in.Expires))
			tab.State = in.State
			if e := tab.Save(ctx); e != nil {
				ctx.LogError(e)
				return &zauth_pb.PermissionBind_RES{Code: e.GetCode()}
			}

			return &zauth_pb.PermissionBind_RES{
				Code: zservice.Code_SUCC,
				Info: &zauth_pb.PermissionBindInfo{
					TargetType:   tab.TargetType,
					TargetID:     tab.TargetID,
					PermissionID: tab.PermissionID,
					Expires:      tab.Expires.UnixMilli(),
					State:        tab.State,
				}}
		}
	}
	// 创建新数据
	tab := &PermissionBindTable{
		TargetType:   in.TargetType,
		TargetID:     in.TargetID,
		PermissionID: in.PermissionID,
		Expires:      zservice.NewTime(time.UnixMilli(in.Expires)),
		State:        in.State,
	}

	if e := tab.Save(ctx); e != nil {
		ctx.LogError(e)
		return &zauth_pb.PermissionBind_RES{Code: e.GetCode()}
	}

	return &zauth_pb.PermissionBind_RES{Code: zservice.Code_SUCC, Info: &zauth_pb.PermissionBindInfo{
		TargetType:   tab.TargetType,
		TargetID:     tab.TargetID,
		PermissionID: tab.PermissionID,
		Expires:      tab.Expires.UnixMilli(),
		State:        tab.State,
	}}
}
