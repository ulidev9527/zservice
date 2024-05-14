package internal

const (
	E_PermissionAction_Create = 1
	E_PermissionAction_Delete = 2
	E_PermissionAction_Open   = 3
	E_PermissionAction_Close  = 4

	RK_Token     = "zauth:token:%s"      // toekn 缓存
	RK_TokenLock = "zauth:token:%s_lock" // toekn 缓存锁
)
