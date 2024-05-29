package internal

import (
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/zglobal"
)

var ZauthInitService *zservice.ZService

// 系统初始化
func ZAuthInit() {

	ZauthDBInit()

	ZauthInitService.StartDone()
}

// 系统数据库数据初始化
func ZauthDBInit() {
	ctx := zservice.NewEmptyContext()
	// 检查账号表是否为空，为空表示未初始化
	count := int64(0)
	if e := Mysql.Model(&ZauthAccountTable{}).Count(&count).Error; e != nil {
		ctx.LogPanic(e)
	}

	if count > 0 {
		ctx.LogInfo("DB is init")
		return
	}
	ctx.LogInfo("DB init start")

	// 添加管理员账号
	admAcc, e := CreateAccount(ctx)
	if e != nil {
		ctx.LogPanic(e)
	} else {
		adminName := zservice.RandomString(9)
		adminPass := zservice.RandomString(16)
		adminPassMd5 := zservice.MD5String(adminPass)

		ctx.LogInfof("Create admin account succ, AdminName: %s, AdminPass: %s PassMD5: %s", adminName, adminPass, adminPassMd5)
		if e := admAcc.AddLoginNameAndPassword(ctx, adminName, zservice.MD5String(adminPass)); e != nil {
			ctx.LogPanic(e)
		}
	}

	// 创建系统组
	sysOrg, e := CreateRootOrg(ctx, "系统管理")
	if e != nil {
		ctx.LogPanic(e)
	}
	// 超级管理员
	adminOrg, e := CreateOrg(ctx, "超级管理员", sysOrg.OrgID)
	if e != nil {
		ctx.LogPanic(e)
	}

	// 账号和组绑定
	_, e = AccountJoinOrg(ctx, admAcc.UID, sysOrg.OrgID, nil) // 加入系统组
	if e != nil {
		ctx.LogPanic(e)
	}
	_, e = AccountJoinOrg(ctx, admAcc.UID, adminOrg.OrgID, nil) // 加入超级管理员组
	if e != nil {
		ctx.LogPanic(e)
	}

	// 添加权限/权限绑定
	func() {
		// 创建权限
		pt := Logic_PermissionCreate(ctx, &zauth_pb.PermissionInfo{
			Name:    "授权系统",
			Service: zservice.GetServiceName(),
			State:   2,
		})
		if pt.Code != zglobal.Code_SUCC {
			ctx.LogPanic(pt)
		}

		// 权限绑定
		_, e = PermissionBind(ctx, 1, adminOrg.OrgID, uint(pt.Info.PermissionID), nil, true)
		if e != nil {
			ctx.LogPanic(e)
		}

		// 开放登陆接口
		pt = Logic_PermissionCreate(ctx, &zauth_pb.PermissionInfo{
			Name:    "授权系统登陆",
			Service: zservice.GetServiceName(),
			Path:    "/login",
			State:   1,
		})
		if pt.Code != zglobal.Code_SUCC {
			ctx.LogPanic(pt)
		}
	}()

	ctx.LogInfo("DB init end")

}
