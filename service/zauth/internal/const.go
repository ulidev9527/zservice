package internal

const (
	E_PermissionAction_Create = 1
	E_PermissionAction_Delete = 2
	E_PermissionAction_Open   = 3
	E_PermissionAction_Close  = 4

	RK_TokenInfo          = "auth:token:info:%v"                // toekn 缓存 %v:token val:AuthToken
	RK_OrgInfo            = "auth:org:info:%v"                  // 组织信息缓存 %v:OrgID val:OrgTable
	RK_OrgRootName        = "auth:org:rootName:%v"              // 组织信息缓存 %v:name_MD5 val:OrgID
	RK_UserInfo           = "auth:user:info:%v"                 // 账号信息缓存 %v:UID val:UserTable
	RK_UserLoginName      = "auth:user:loginName:%v"            // 账号信息缓存 %v:loginName val:UID
	RK_UserLoginPhone     = "auth:user:loginPhone:%v"           // 账号信息缓存 %v:手机号 val:UID
	RK_UserLoginServices  = "auth:user:loginService:%v:%v"      // 账号登陆后的目标service缓存 %v:uid %v:service val:[AuthToken]
	RK_UserOrgBind_Info   = "auth:userOrg:bindInfo:%v:%v"       // 账号组织绑定 %v:UID %v:OrgID val:UserOrgBindTable
	RK_PermissionInfo     = "auth:permission:info:%v"           // 权限信息 %v:PermissionID
	RK_PermissionBindInfo = "auth:permission:bindInfo:%v:%v:%v" // 权限绑定 %v:TargetType %v:TargetID %v:PermissionID val:PermissionBindTable
	RK_PermissionSAP      = "auth:permission:sap:%v:%v:%v"      // 权限绑定 %v:service %v:action %v:path val:PermissionID
	RK_Sms_PhoneBan       = "sms:phoneBan:%s"                   // 手机号封禁 %s:phone %s:BanTime val:1
	RK_Sms_PhoneCD        = "sms:phoneCD:%s"                    // 手机号验证码发送CD %s:手机号
	RK_Sms_PhoneCode      = "sms:phoneCode:%s:%s"               // 手机号验证码 %s:手机号 %s:验证码

	RK_Service_Regist = "service:regist:%s" // 服务注册 %s:ServiceName
	RK_ServiceKVInfo  = "service:kv:%s:%s"  // 服务键值对信息 %s:serviceName %s:key val:value

	RK_AssetInfo       = "asset:assetInfo:%s"          // 资产信息 %s:AssetID val:AssetTable
	RK_ConfigAssetInfo = "asset:configAssetInfo:%s:%s" // 配置资源缓存 %s:ServiceName %s:FileName val:ConfigAssetTable

	KV_Service_SMS_VerifyCodeSend_NoSend = "sms/verifyCodeSend/NoSend" // 服务配置 不发送验证码
	KV_Service_Login_AllowMPOP           = "login/AllowMPOP"           // 是否允许服务进行多点登录
	KV_ZZZZString                        = "ZZZZString"                // ZZZZString

)
