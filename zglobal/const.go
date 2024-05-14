package zglobal

import "time"

const (
	Code_SUCC          = 1 // 成功
	Code_Fail          = 2 // 失败
	Code_ErrorBreakoff = 3 // 中断
	Code_AuthFail      = 4 // 鉴权失败
	Code_NotImplement  = 5 // 未实现

	Code_OpenFileErr      = 101 // 打开文件错误
	Code_WiteFileErr      = 102 // 写入文件错误
	Code_CloseFileErr     = 103 // 关闭文件错误
	Code_EmptyFile        = 104 // 文件为空
	Code_RedisKeyLockFail = 105 // redis 锁失败

	Code_Zsms_Phone_NULL          = 1001 // 手机号为空
	Code_Zsms_Phone_VerifyFail    = 1002 // 手机号验证失败
	Code_Zsms_Phone_Ban           = 1003 // 手机号被封禁
	Code_Zsms_Phone_CD            = 1004 // 手机号验证码CD中
	Code_Zsms_Phone_CodeCacheNull = 1005 // 手机号验证码不存在, 没有缓存
	Code_Zsms_Phone_CodeNull      = 1006 // 手机号验证码为空
	Code_Zsms_Phone_CodeLenErr    = 1007 // 手机号验证码长度错误
	Code_Zsms_ErrorBreakoff       = 1008 // 短信中断错误
	Code_Zsms_SendParamsErr       = 1009 // 短信发送参数错误

	Code_Zconfig_ParamsErr        = 2001 // 参数错误
	Code_Zconfig_ParserNotExist   = 2002 // 没有这个解析器
	Code_Zconfig_FileNotExist     = 2003 // 文件不存在
	Code_Zconfig_ParserFail       = 2004 // 解析失败
	Code_Zconfig_PathIsDir        = 2005 // 文件是个目录
	Code_Zconfig_GetFileMd5Fail   = 2006 // 读取文件 md5 失败
	Code_Zconfig_FileMd5NotChange = 2007 // 文件 md5 未变化
	Code_Zconfig_ExcelNoContent   = 2008 // excel 内容为空
	Code_Zconfig_GetConfigFail    = 2009 // 获取配置失败

	Code_Zauth_TokenSaveFail = 3101 // token 存储失败
	Code_Zauth_TokenIsNil    = 3102 // token 为空
	Code_Zauth_TokenSignFail = 3103 // token 签名错误

	E_ZConfig_Parser_Excel = 1 // excel 解析器

	Time_1m = time.Second * 60 // 1 分钟
)
