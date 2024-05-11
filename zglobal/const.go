package zglobal

const (
	E_SUCC                = 1    // 成功
	E_Fail                = 2    // 失败
	E_ErrorBreakoff       = 3    // 中断
	E_Phone_NULL          = 1001 // 手机号为空
	E_Phone_VerifyFail    = 1002 // 手机号验证失败
	E_Phone_Ban           = 1003 // 手机号被封禁
	E_Phone_CD            = 1004 // 手机号验证码CD中
	E_Phone_CodeCacheNull = 1005 // 手机号验证码不存在, 没有缓存
	E_Phone_CodeNull      = 1006 // 手机号验证码为空
	E_Phone_CodeLenErr    = 1007 // 手机号验证码长度错误

	E_SMS_ErrorBreakoff = 2001 // 短信发送中断错误
	E_SMS_SendParamsErr = 2002 // 短信发送参数错误

)
