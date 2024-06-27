package internal

const (
	E_PermissionAction_Create = 1
	E_PermissionAction_Delete = 2
	E_PermissionAction_Open   = 3
	E_PermissionAction_Close  = 4

	RK_TokenInfo                   = "token:info:%v"                      // toekn 缓存 %v:token val:AuthToken
	RK_OrgInfo                     = "org:info:%v"                        // 组织信息缓存 %v:OrgID val:OrgTable
	RK_OrgRootName                 = "org:rootName:%v"                    // 组织信息缓存 %v:name_MD5 val:OrgID
	RK_UserInfo                    = "user:info:%v"                       // 账号信息缓存 %v:UID val:UserTable
	RK_UserLoginName               = "user:loginName:%v"                  // 账号信息缓存 %v:loginName val:UID
	RK_UserLoginPhone              = "user:loginPhone:%v"                 // 账号信息缓存 %v:手机号 val:UID
	RK_UserLoginServices           = "user:loginService:%v:%v"            // 账号登陆后的目标service缓存 %v:uid %v:service val:[AuthToken]
	RK_AOBind_Info                 = "AOBind:link:%v:%v"                  // 账号组织绑定 %v:OrgID %v:UID val:BindID
	RK_PermissionInfo              = "permission:info:%v"                 // 权限信息 %v:PermissionID
	RK_PermissionBindInfo          = "permission:bindInfo:%v:%v:%v"       // 权限绑定 %v:TargetType %v:TargetID %v:PermissionID val:PermissionBindTable
	RK_PermissionSAP               = "permission:sap:%v:%v:%v"            // 权限绑定 %v:service %v:action %v:path val:PermissionID
	RK_Sms_PhoneBan                = "sms:phoneBan:%s"                    // 手机号封禁 %s:phone %s:BanTime val:1
	RK_Sms_PhoneCD                 = "sms:phoneCD:%s"                     // 手机号验证码发送CD %s:手机号
	RK_Sms_PhoneCode               = "sms:phoneCode:%s:%s"                // 手机号验证码 %s:手机号 %s:验证码
	RK_Config_ServiceFileConfig    = "config:service:%s:fileConfig:%s"    // 文件配置缓存
	RK_Config_ServiceFileConfigMD5 = "config:service:%s:fileConfigMD5:%s" // 文件配置的 md5, 用于标识是否需要重新解析
	RK_Config_ServiceEnvAuth       = "config:serviceEnvAuth:%s"           // 服务环境配置缓存
	RK_Service_Regist              = "service:regist:%s"                  // 服务注册 %s:ServiceName
	RK_AssetToken                  = "asset:token:%s"                     // 资产信息 %s:token val:AssetTable
	RK_AssetMd5                    = "asset:md5:%s"                       // 资产信息 %s:md5 val:1
	RK_ServiceKVInfo               = "service:kv:%s:%s"                   // 服务键值对信息 %s:serviceName %s:key val:value

	KV_Service_SMS_VerifyCodeSend_NoSend = "sms/verifyCodeSend/NoSend" // 服务配置 不发送验证码
	KV_Service_Login_AllowMPOP           = "login/AllowMPOP"           // 是否允许服务进行多点登录

	EV_Config_ServiceFileConfigChange = "ConfigService_%s_FileConfigChange" // 服务文件配置变更通知
	EV_Config_ZZZZStringChange        = "ZZZZStringChange"                  // ZZZZ 字符串变更

	FI_ServiceConfigFile = "static/fileConfig/%s/%s" // 服务配置文件路径
	FI_ServiceEnvFile    = "static/envConfig/%s.env" // 服务环境文件路径
	FI_UploadDir         = "static/upload/%s"        // 上传文件路径
	FI_ZZZZStringFile    = "static/zzzzString.txt"   // zzzz 字符串
)
