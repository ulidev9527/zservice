package internal

const (
	E_PermissionAction_Create = 1
	E_PermissionAction_Delete = 2
	E_PermissionAction_Open   = 3
	E_PermissionAction_Close  = 4

	RK_Token                = "zauth:token:%v"                   // toekn 缓存
	RK_TokenLock            = "zauth:token:%v_lock"              // toekn 缓存锁
	RK_OrgCreateLock        = "zauth:org:create_lock"            // 组织创建锁
	RK_OrgInfo              = "zauth:org:info:%v"                // 组织信息缓存 %vOrgID
	RK_AccountCreateLock    = "zauth:account:create_lock"        // 账号创建锁
	RK_AccountInfo          = "zauth:account:info:%v"            // 账号信息缓存 存储表的数据
	RK_AccountLoginName     = "zauth:account:loginName:%v"       // 账号信息缓存 AccountID
	RK_PermissionCreateLock = "zauth:permission:create_lock"     // 权限创建锁
	RK_PermissionInfo       = "zauth:permission:info:%v"         // 权限信息
	RK_AccountBindOrgInfo   = "zauth:account:bindOrg:info:%v:%v" // 账号组织绑定 %v组织%v账号
)
