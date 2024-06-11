package zglobal

import "time"

const (
	Code_SUCC          = 1 // 成功
	Code_Fail          = 2 // 失败
	Code_ErrorBreakoff = 3 // 中断
	Code_Zauth_Fail    = 4 // 鉴权失败
	Code_NotImplement  = 5 // 未实现
	Code_ParamsErr     = 6 // 参数错误
	Code_LoginAgain    = 7 // 请重新登录, 重新拉取 token 进行登陆
	Code_NotFound      = 8 // 资源不存在/没找到/数据未查询到
	Code_Wait          = 9 // 等待

	Code_OpenFileErr         = 101 // 打开文件错误
	Code_WiteFileErr         = 102 // 写入文件错误
	Code_CloseFileErr        = 103 // 关闭文件错误
	Code_EmptyFile           = 104 // 文件为空
	Code_Already_Lock        = 105 // 已经上锁
	Code_SyncCacheIncomplete = 106 // 数据同步不完全
	Code_SyncCacheErr        = 107 // 数据同步失败
	Code_GenIDCountMaxErr    = 108 // 生成ID错误次数超上限

	Code_DB_SaveFail   = 201 // 数据存储失败
	Code_Redis_DelFail = 202 // redis 删除失败

	Code_Zauth_Sms_Phone_NULL                = 1001 // 手机号为空
	Code_Zauth_Sms_Phone_VerifyFail          = 1002 // 手机号验证失败
	Code_Zauth_Sms_Phone_Ban                 = 1003 // 手机号被封禁
	Code_Zauth_Sms_Phone_CD                  = 1004 // 手机号验证码CD中
	Code_Zauth_Sms_Phone_VerifyCodeCacheNull = 1005 // 手机号验证码不存在, 没有缓存
	Code_Zauth_Sms_Phone_CodeNull            = 1006 // 手机号验证码为空
	Code_Zauth_Sms_Phone_VerifyCodeLenErr    = 1007 // 手机号验证码长度错误
	Code_Zauth_Sms_Phone_VerifyCodeErr       = 1008 // 手机号验证码错误

	Code_Zauth_TokenSaveFail = 1201 // token 存储失败
	Code_Zauth_TokenIsNil    = 1202 // token 为空
	Code_Zauth_TokenSignFail = 1203 // token 签名错误
	Code_Zauth_TokenDelFail  = 1204 // token 删除失败

	Code_Zauth_OrgCreateRootErr     = 1301 // 创建组织根错误
	Code_Zauth_OrgCreateErr         = 1302
	Code_Zauth_OrgCreateRootIDErr   = 1303 // 创建组织根ID错误
	Code_Zauth_OrgCreateParentIDErr = 1304 // 创建组织父ID错误
	// Code_Zauth_OrgGenIDCountMaxErr  = 1305 // 生成组织ID错误次数超上限
	Code_Zauth_Org_NotFund      = 1306 // 组织不存在
	Code_Zauth_Org_AlreadyExist = 1307 // 组织已经存在

	// Code_Zauth_AccountGenIDCountMaxErr       = 1401 // 生成账号ID错误次数超上限
	Code_Zauth_AccountAlreadyJoin_Org        = 1402 // 账号已经加入组织
	Code_Zauth_AccountAlreadyExist_LoginName = 1403 // 账号已经存在
	Code_Zauth_Account_NotFund               = 1404 // 账号不存在

	// Code_Zauth_PermissionGenIDCountMaxErr     = 1501 // 生成权限ID错误次数超上限
	Code_Zauth_PermissionBind_TargetTypeErr   = 1502 // 权限绑定目标类型错误
	Code_Zauth_PermissionBind_TargetIDErr     = 1503 // 权限绑定目标ID错误
	Code_Zauth_PermissionBind_PermissionIDErr = 1504 // 权限绑定权限ID错误
	Code_Zauth_PermissionBind_Already_Bind    = 1505 // 权限已经绑定
	Code_Zauth_Permission_NotFound            = 1506 // 权限不存在
	Code_Zauth_Permission_ConfigErr           = 1507 // 权限配置错误
	Code_Zauth_Permission_Alerady_Exist       = 1508 // 权限已经存在

	Code_Zauth_Login_Account_NotFund = 1601 // 账号不存在
	Code_Zauth_Login_Pass_Err        = 1602 // 密码错误

	Code_Zauth_config_ParamsErr        = 2001 // 参数错误
	Code_Zauth_config_ParserNotExist   = 2002 // 没有这个解析器
	Code_Zauth_config_FileNotExist     = 2003 // 文件不存在
	Code_Zauth_config_ParserFail       = 2004 // 解析失败
	Code_Zauth_config_PathIsDir        = 2005 // 文件是个目录
	Code_Zauth_config_GetFileMd5Fail   = 2006 // 读取文件 md5 失败
	Code_Zauth_config_FileMd5NotChange = 2007 // 文件 md5 未变化
	Code_Zauth_config_ExcelNoContent   = 2008 // excel 内容为空
	Code_Zauth_config_GetConfigFail    = 2009 // 获取配置失败

	E_ZConfig_Parser_Excel = 1 // excel 解析器

	Time_0     = time.Duration(0)      // 0 秒
	Time_10ms  = time.Millisecond * 10 // 10 毫秒
	Time_1m    = time.Second * 60      // 1 分钟
	Time_10Day = time.Hour * 24 * 10   // 10 天
	Time_3Day  = time.Hour * 24 * 3    // 3 天

	NSQ_Topic_zlog_AddKV = "zlog_add_kv" // 添加 kv 日志
)
