package internal

import "zservice/zservice"

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

		ctx.LogInfof("Create admin account succ, AdminName: %s, AdminPass: %s", adminName, adminPass)
		if e := admAcc.AddLoginNameAndPassword(ctx, adminName, adminPass); e != nil {
			ctx.LogPanic(e)
		}
	}

	// 创建系统组
	sysOrg, e := CreateRootOrg(ctx, "系统管理")
	if e != nil {
		ctx.LogPanic(e)
	}
	// 超级管理员
	adminOrg, e := CreateOrg(ctx, "超级管理员", sysOrg.OrgID, sysOrg.OrgID)
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
		pt, e := CreatePermission(ctx, ZauthPermissionTable{
			Name:    "权限服务",
			Service: "zauth",
			State:   2,
		})
		if e != nil {
			ctx.LogPanic(e)
		}

		// 权限绑定
		_, e = PermissionBind(ctx, 1, adminOrg.OrgID, pt.PermissionID, nil, 1)
		if e != nil {
			ctx.LogPanic(e)
		}

	}()

	ctx.LogInfo("DB init end")

}
