package internal

const (
	E_PermissionAction_Create = 1
	E_PermissionAction_Delete = 2
	E_PermissionAction_Open   = 3
	E_PermissionAction_Close  = 4

	RK_TokenInfo            = "zauth:token:info:%v"                // toekn 缓存
	RK_TokenLock            = "zauth:token:%v_lock"                // toekn 缓存锁
	RK_OrgInfo              = "zauth:org:info:%v"                  // 组织信息缓存 %vOrgID
	RK_OrgRootName          = "zauth:org:rootName:%v"              // 组织信息缓存 %v_nameMD5 存储OrgID
	RK_AccountInfo          = "zauth:account:info:%v"              // 账号信息缓存 存储表的数据
	RK_AccountLoginName     = "zauth:account:loginName:%v"         // 账号信息缓存 AccountID
	RK_AccountLoginPhone    = "zauth:account:loginPhone:%v"        // 账号信息缓存 AccountID
	RK_AccountLoginToken    = "zauth:account:loginToken:%v:%v"     // 账号登陆后的目标token缓存 %v_uid %v_token
	RK_AOBind_CreateLock    = "zauth:AOBind:create_lock:%v:%v"     // 组织绑定 %v组织ID %v账号ID
	RK_AOBind_Info          = "zauth:AOBind:link:%v:%v"            // 账号组织绑定 %v组织ID %v账号ID 存储 BindID
	RK_PermissionCreateLock = "zauth:permission:create_lock"       // 权限创建锁
	RK_PermissionInfo       = "zauth:permission:info:%v"           // 权限信息
	RK_PermissionBindInfo   = "zauth:permission:bindInfo:%v:%v:%v" // 权限绑定 %v目标类型 %v目标ID %v权限ID
	RK_PermissionSAP        = "zauth:permission:sap:%v:%v:%v"      // 权限绑定 %v_service %v_action %v_path
	RK_PhoneBan             = "zsms:PhoneBan:%s"                   // 手机号封禁 %s 手机号, %s 封禁时间
	RK_PhoneCD              = "zsms:PhoneCD:%s"                    // 手机号验证码发送CD %s 手机号
	RK_PhoneCode            = "zsms:PhoneCode:%s"                  // 手机号验证码 %s 手机号
	RK_FileConfig           = "zconfig:fileConfig:%s"              // 文件配置缓存
	RK_FileMD5              = "zconfig:fileMD5:%s"                 // 文件配置的 md5, 用于标识是否需要重新解析

	NSQ_FileConfig_Change = "zconfig_fileConfig_change" // 文件配置变更通知

	FI_StaticRoot = "static" // 静态文件根目录

)
