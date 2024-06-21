package internal

import (
	"fmt"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/ex/gormservice"
	"zservice/zservice/zglobal"
)

func Logic_ServiceRegist(ctx *zservice.Context, in *zauth_pb.ServiceRegist_REQ) *zauth_pb.ServiceRegist_RES {

	if ctx.TraceService == "" {
		return &zauth_pb.ServiceRegist_RES{Code: zglobal.Code_ParamsErr}
	}

	// 锁
	rk_regist := fmt.Sprintf(RK_Service_Regist, ctx.TraceService)
	un, e := Redis.Lock(rk_regist)
	if e != nil {
		ctx.LogError(e)
		return &zauth_pb.ServiceRegist_RES{Code: zglobal.Code_Fail}
	}
	defer un()

	// 服务组
	orgInfo := &zauth_pb.OrgInfo{}
	if tab, e := GetRootOrgByName(ctx, ctx.TraceService); e != nil {
		if e.GetCode() != zglobal.Code_NotFound {
			ctx.LogError(e)
			return &zauth_pb.ServiceRegist_RES{Code: e.GetCode()}
		} else {
			// 创建服务组
			res := Logic_OrgCreate(ctx, &zauth_pb.OrgInfo{Name: ctx.TraceService})
			if res.Code == zglobal.Code_SUCC {
				orgInfo = res.Info
			} else {
				return &zauth_pb.ServiceRegist_RES{Code: res.Code}
			}
		}
	} else if tab != nil {
		orgInfo = &zauth_pb.OrgInfo{
			OrgID:    tab.OrgID,
			Name:     tab.Name,
			RootID:   tab.RootID,
			ParentID: tab.ParentID,
			State:    tab.State,
		}
	}

	// 服务权限
	permissionInfo := &zauth_pb.PermissionInfo{}
	if tab, e := GetPermissionBySAP(ctx, ctx.TraceService, "", ""); e != nil {
		if e.GetCode() != zglobal.Code_NotFound {
			ctx.LogError(e)
			return &zauth_pb.ServiceRegist_RES{Code: e.GetCode()}
		} else {
			// 创建权限
			permissionRes := Logic_PermissionCreate(ctx, &zauth_pb.PermissionInfo{
				Name:    ctx.TraceService,
				Service: ctx.TraceService,
				State:   2,
			})
			// 是否异常
			if permissionRes.Code != zglobal.Code_SUCC {
				return &zauth_pb.ServiceRegist_RES{Code: permissionRes.Code}
			}
			permissionInfo = permissionRes.Info
		}
	} else {
		permissionInfo = &zauth_pb.PermissionInfo{
			PermissionID: tab.PermissionID,
			Name:         tab.Name,
			Service:      tab.Service,
			Path:         tab.Path,
			State:        tab.State,
		}
	}

	// 服务和权限绑定
	if _, e := GetPermissionBind(ctx, 1, orgInfo.OrgID, permissionInfo.PermissionID); e != nil {
		if e.GetCode() != zglobal.Code_NotFound {
			ctx.LogError(e)
			return &zauth_pb.ServiceRegist_RES{Code: e.GetCode()}
		} else {

			bindRes := Logic_PermissionBind(ctx, &zauth_pb.PermissionBind_REQ{
				TargetType:   1,
				TargetID:     orgInfo.OrgID,
				PermissionID: permissionInfo.PermissionID,
				State:        1,
			})

			if bindRes.GetCode() != zglobal.Code_SUCC {
				return &zauth_pb.ServiceRegist_RES{Code: bindRes.GetCode()}
			}
		}
	}

	// 管理员账号检查
	// 是否有账号和当前组绑定
	adminUserTab := &UserTable{}
	userBindTab := &UserOrgBindTable{}
	if e := Mysql.Model(&UserOrgBindTable{}).Where("org_id = ?", orgInfo.OrgID).First(userBindTab).Error; e != nil {
		if !gormservice.IsNotFound(e) { // 其他错误
			ctx.LogError(e)
			return &zauth_pb.ServiceRegist_RES{Code: zglobal.Code_Fail}
		} else {
			// 创建管理员账号
			if tab, e := CreateUser(ctx); e != nil {
				ctx.LogError(e)
				return &zauth_pb.ServiceRegist_RES{Code: e.GetCode()}
			} else {
				adminUserTab = tab
			}
		}
	} else {
		if tab, e := GetUserByUID(ctx, userBindTab.UID); e != nil {
			ctx.LogError("user get error")
			ctx.LogError(e)
			return &zauth_pb.ServiceRegist_RES{Code: e.GetCode()}
		} else {
			adminUserTab = tab
		}
	}

	// 添加登陆名和密码
	if adminUserTab.LoginName == "" {
		resultCode := uint32(zglobal.Code_SUCC)
		count := 0
		for {
			if count > 10 {
				resultCode = zglobal.Code_Limit
				break
			}
			adminName := zservice.RandomString(9)
			adminPass := zservice.RandomString(16)
			adminPassMd5 := zservice.MD5String(adminPass)

			if e := adminUserTab.AddLoginNameAndPassword(ctx, adminName, zservice.MD5String(adminPass)); e != nil {
				if e.GetCode() != zglobal.Code_Zauth_UserAlreadyExist_LoginName {
					ctx.LogError(e)
					resultCode = e.GetCode()
					break
				}
				count++
			} else {
				ctx.LogInfo("Create admin user --------------------")
				ctx.LogWarnf("Service: %s, AdminName: %s, AdminPass: %s PassMD5: %s", ctx.TraceService, adminName, adminPass, adminPassMd5)
				ctx.LogInfo("Create admin user --------------------")
				break
			}
		}
		if resultCode != zglobal.Code_SUCC {
			return &zauth_pb.ServiceRegist_RES{Code: resultCode}
		}
	}

	// 管理员和组绑定
	if res := Logic_AddUserToOrg(ctx, &zauth_pb.AddUserToOrg_REQ{Uid: adminUserTab.UID, OrgID: orgInfo.OrgID}); res.Code != zglobal.Code_SUCC && res.Code != zglobal.Code_RepetitionErr {
		return &zauth_pb.ServiceRegist_RES{Code: res.Code}
	} else {
		return &zauth_pb.ServiceRegist_RES{
			Code:           zglobal.Code_SUCC,
			OrgInfo:        orgInfo,
			PermissionInfo: permissionInfo,
		}
	}
}
