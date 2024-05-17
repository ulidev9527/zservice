package internal

import "zservice/zservice"

var SystemService *zservice.ZService

// 系统初始化
func SystemInit() {

	SystemDBInit()

	SystemService.StartDone()
}

// 系统数据库数据初始化
func SystemDBInit() {
	ctx := zservice.NewEmptyContext()
	// 检查是否有管理员账号，有管理员账号表示已经初始化过了
	count := int64(0)
	if e := Mysql.Model(&ZauthAccountTable{}).Where("account = ?", "admin").Count(&count).Error; e != nil {
		ctx.LogPanic(e)
	}

	if count > 0 {
		ctx.LogInfo("DB is init")
		return
	}
	ctx.LogInfo("DB init start")

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

	// 添加管理员账号
	admAcc, e := CreateAccount(ctx)
	if e != nil {
		ctx.LogPanic(e)
	} else {
		if e := admAcc.AddLoginNameAndPassword(ctx, "admin", "admin"); e != nil {
			ctx.LogPanic(e)
		}
	}

	// 账号和组绑定
	_, e = AccountJoinOrg(ctx, admAcc.AccountID, sysOrg.OrgID, nil) // 加入系统组
	if e != nil {
		ctx.LogPanic(e)
	}
	_, e = AccountJoinOrg(ctx, admAcc.AccountID, adminOrg.OrgID, nil) // 加入超级管理员组
	if e != nil {
		ctx.LogPanic(e)
	}

	// 添加权限
	func() {

		type T_Permission struct {
			Name    string // 权限名称
			Service string // 权限服务
			Action  string // 权限动作
			Path    string // 权限路径
			State   uint   // 权限状态
			Child   []T_Permission
		}

		pObj := T_Permission{
			Name: "认证服务", Service: "zauth", State: 2,
			Child: []T_Permission{
				{Name: "授权管理", Service: "zauth", Action: "", Path: "/permission", State: 3, Child: []T_Permission{
					{Name: "获取权限", Service: "zauth", Action: "get", Path: "/permission", State: 3},
					{Name: "创建权限", Service: "zauth", Action: "post", Path: "/permission", State: 3},
					{Name: "修改权限", Service: "zauth", Action: "put", Path: "/permission", State: 3},
				}},
				{Name: "组织管理", Service: "zauth", Action: "", Path: "/org", State: 3, Child: []T_Permission{
					{Name: "获取组织", Service: "zauth", Action: "get", Path: "/org", State: 3},
					{Name: "创建组织", Service: "zauth", Action: "post", Path: "/org", State: 3},
					{Name: "修改组织", Service: "zauth", Action: "put", Path: "/org", State: 3},
				}},
			},
		}
		ctx.LogInfo(pObj)

		_, e := CreatePermission(ctx, ZauthPermissionTable{
			Name:    "权限服务",
			Service: "zauth",
			State:   2,
		})
		if e != nil {
			ctx.LogPanic(e)
		}
	}()

	ctx.LogInfo("DB init end")

}
