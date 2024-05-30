package internal

const (
	E_PermissionAction_Create = 1
	E_PermissionAction_Delete = 2
	E_PermissionAction_Open   = 3
	E_PermissionAction_Close  = 4

	RK_TokenInfo           = "token:info:%v"                // toekn 缓存
	RK_TokenLock           = "token:%v_lock"                // toekn 缓存锁
	RK_OrgInfo             = "org:info:%v"                  // 组织信息缓存 %vOrgID
	RK_OrgRootName         = "org:rootName:%v"              // 组织信息缓存 %v_nameMD5 存储OrgID
	RK_AccountInfo         = "account:info:%v"              // 账号信息缓存 存储表的数据
	RK_AccountLoginName    = "account:loginName:%v"         // 账号信息缓存 存储 AccountID
	RK_AccountLoginPhone   = "account:loginPhone:%v"        // 账号信息缓存 存储 AccountID
	RK_AccountLoginService = "account:loginService:%v:%v"   // 账号登陆后的目标service缓存 %v_uid %v_service 存储 AuthToken
	RK_AOBind_Info         = "AOBind:link:%v:%v"            // 账号组织绑定 %v组织ID %v账号ID 存储 BindID
	RK_PermissionInfo      = "permission:info:%v"           // 权限信息
	RK_PermissionBindInfo  = "permission:bindInfo:%v:%v:%v" // 权限绑定 %v目标类型 %v目标ID %v权限ID
	RK_PermissionSAP       = "permission:sap:%v:%v:%v"      // 权限绑定 %v_service %v_action %v_path
	RK_PhoneBan            = "sms:phoneBan:%s"              // 手机号封禁 %s 手机号, %s 封禁时间
	RK_PhoneCD             = "sms:phoneCD:%s"               // 手机号验证码发送CD %s 手机号
	RK_PhoneCode           = "sms:phoneCode:%s"             // 手机号验证码 %s 手机号
	RK_FileConfig          = "config:fileConfig:%s"         // 文件配置缓存
	RK_FileMD5             = "config:fileMD5:%s"            // 文件配置的 md5, 用于标识是否需要重新解析

	NSQ_FileConfig_Change = "zconfig_fileConfig_change" // 文件配置变更通知

	FI_StaticRoot = "static" // 静态文件根目录

)
