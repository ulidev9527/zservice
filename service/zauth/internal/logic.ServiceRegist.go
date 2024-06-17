package internal

import (
	"fmt"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/zglobal"
)

func Logic_ServiceRegist(ctx *zservice.Context, in *zauth_pb.Default_REQ) *zauth_pb.Default_RES {

	if ctx.TraceService == ctx.NowService {
		return &zauth_pb.Default_RES{Code: zglobal.Code_SUCC}
	}

	// 锁
	rk_regist := fmt.Sprintf(RK_Service_Regist, ctx.TraceService)
	un, e := Redis.Lock(rk_regist)
	if e != nil {
		ctx.LogError(e)
		return &zauth_pb.Default_RES{Code: zglobal.Code_ErrorBreakoff}
	}
	defer un()

	// 检查服务组
	orgID := uint32(0)
	if tab, e := GetRootOrgByName(ctx, ctx.TraceService); e != nil {
		ctx.LogError(e)
		return &zauth_pb.Default_RES{Code: e.GetCode()}
	} else if tab != nil {
		orgID = tab.OrgID
	} else {
		// 创建服务组
		res := Logic_OrgCreate(ctx, &zauth_pb.OrgInfo{Name: ctx.TraceService})
		if res.Code == zglobal.Code_SUCC {
			orgID = res.Info.OrgID
		} else {
			return &zauth_pb.Default_RES{Code: res.Code}
		}
	}

	// 创建服务权限
	permissionRes := Logic_PermissionCreate(ctx, &zauth_pb.PermissionInfo{
		Name:    ctx.TraceService,
		Service: ctx.TraceService,
		State:   2,
	})
	// 是否异常
	if permissionRes.Code != zglobal.Code_SUCC && permissionRes.Code != zglobal.Code_Zauth_Permission_Alerady_Exist {
		return &zauth_pb.Default_RES{Code: permissionRes.Code}
	}

	// 服务和权限绑定
	bindRes := Logic_PermissionBind(ctx, &zauth_pb.PermissionBind_REQ{
		TargetType:   1,
		TargetID:     orgID,
		PermissionID: permissionRes.Info.PermissionID,
		State:        1,
	})
	return &zauth_pb.Default_RES{Code: bindRes.Code}
}
