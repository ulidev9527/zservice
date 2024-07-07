package internal

import (
	"fmt"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
)

func Logic_ServiceRegist(ctx *zservice.Context, in *zauth_pb.ServiceRegist_REQ) *zauth_pb.ServiceRegist_RES {

	if in.Service == "" {
		ctx.LogError("param error")
		return &zauth_pb.ServiceRegist_RES{Code: zservice.Code_ParamsErr}
	}

	// 锁
	rk_regist := fmt.Sprintf(RK_Service_Regist, in.Service)
	un, e := Redis.Lock(rk_regist)
	if e != nil {
		ctx.LogError(e)
		return &zauth_pb.ServiceRegist_RES{Code: zservice.Code_Fail}
	}
	defer un()

	// 服务组
	orgInfo := &zauth_pb.OrgInfo{}
	if tab, e := GetRootOrgByName(ctx, in.Service); e != nil {
		if e.GetCode() != zservice.Code_NotFound {
			ctx.LogError(e)
			return &zauth_pb.ServiceRegist_RES{Code: e.GetCode()}
		} else {
			// 创建服务组
			res := Logic_OrgCreate(ctx, &zauth_pb.OrgInfo{Name: in.Service})
			if res.Code == zservice.Code_SUCC {
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
	if tab, e := GetPermissionBySAP(ctx, in.Service, "", ""); e != nil {
		if e.GetCode() != zservice.Code_NotFound {
			ctx.LogError(e)
			return &zauth_pb.ServiceRegist_RES{Code: e.GetCode()}
		} else {
			// 创建权限
			permissionRes := Logic_PermissionCreate(ctx, &zauth_pb.PermissionInfo{
				Name:    in.Service,
				Service: in.Service,
				State:   2,
			})
			// 是否异常
			if permissionRes.Code != zservice.Code_SUCC {
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
		if e.GetCode() != zservice.Code_NotFound {
			ctx.LogError(e)
			return &zauth_pb.ServiceRegist_RES{Code: e.GetCode()}
		} else {

			bindRes := Logic_PermissionBind(ctx, &zauth_pb.PermissionBind_REQ{
				TargetType:   1,
				TargetID:     orgInfo.OrgID,
				PermissionID: permissionInfo.PermissionID,
				State:        1,
			})

			if bindRes.GetCode() != zservice.Code_SUCC {
				return &zauth_pb.ServiceRegist_RES{Code: bindRes.GetCode()}
			}
		}
	}

	// 管理员账号检查
	// 是否有账号和当前组绑定
	adminUserTab := &UserTable{}
	userBindTab := &UserOrgBindTable{}
	if e := Gorm.First(userBindTab, "org_id = ?", orgInfo.OrgID).Error; e != nil {
		if !DBService.IsNotFoundErr(e) { // 其他错误
			ctx.LogError(e)
			return &zauth_pb.ServiceRegist_RES{Code: zservice.Code_Fail}
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
		resultCode := uint32(zservice.Code_SUCC)
		count := 0
		for {
			if count > 10 {
				resultCode = zservice.Code_Limit
				break
			}
			adminName := zservice.RandomString(9)
			adminPass := zservice.RandomString(16)
			adminPassMd5 := zservice.MD5String(adminPass)

			if e := adminUserTab.AddLoginNameAndPassword(ctx, adminName, zservice.MD5String(adminPass)); e != nil {
				if e.GetCode() != zservice.Code_Zauth_UserAlreadyExist_LoginName {
					ctx.LogError(e)
					resultCode = e.GetCode()
					break
				}
				count++
			} else {
				ctx.LogInfo("Create admin user --------------------")
				ctx.LogWarnf("Service: %s, AdminName: %s, AdminPass: %s PassMD5: %s", in.Service, adminName, adminPass, adminPassMd5)
				ctx.LogInfo("Create admin user --------------------")
				break
			}
		}
		if resultCode != zservice.Code_SUCC {
			return &zauth_pb.ServiceRegist_RES{Code: resultCode}
		}
	}

	// 管理员和组绑定
	if res := Logic_UserOrgBind(ctx, &zauth_pb.UserOrgBind_REQ{Uid: adminUserTab.UID, OrgID: orgInfo.OrgID, State: 1}); res.Code != zservice.Code_SUCC && res.Code != zservice.Code_RepetitionErr {
		return &zauth_pb.ServiceRegist_RES{Code: res.Code}
	}

	// 初始化权限添加
	if len(in.InitPermissions) > 0 {

		for _, pInfo := range in.InitPermissions {
			// 登陆权限添加
			if _, e := GetPermissionBySAP(ctx, in.Service, pInfo.Action, pInfo.Path); e != nil {
				if e.GetCode() != zservice.Code_NotFound {
					ctx.LogPanic(e)
				} else { // 创建权限
					permissionRES := Logic_PermissionCreate(ctx, &zauth_pb.PermissionInfo{
						Service: in.Service,
						Action:  pInfo.Action,
						Path:    pInfo.Path,
						State:   pInfo.State,
					})
					if permissionRES.Code != zservice.Code_SUCC {
						return &zauth_pb.ServiceRegist_RES{Code: permissionRES.Code}
					}
				}
			}
		}

	}

	return &zauth_pb.ServiceRegist_RES{
		Code:           zservice.Code_SUCC,
		OrgInfo:        orgInfo,
		PermissionInfo: permissionInfo,
	}
}
