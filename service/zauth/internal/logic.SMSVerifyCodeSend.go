package internal

import (
	"fmt"
	"time"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/zglobal"
)

// 发送验证码
func Logic_SMSVerifyCodeSend(ctx *zservice.Context, in *zauth_pb.SMSVerifyCodeSend_REQ) *zauth_pb.SMSSendVerifyCode_RES {

	// 参数检查
	if in.Phone == "" || in.Phone[0] != '+' {
		return &zauth_pb.SMSSendVerifyCode_RES{Code: zglobal.Code_ParamsErr}
	}

	// 封禁检查
	if isBan, e := IsSmsBan(ctx, in.Phone); e != nil {
		ctx.LogError(e)
		return &zauth_pb.SMSSendVerifyCode_RES{Code: e.GetCode()}
	} else if isBan {
		ctx.LogError(zservice.NewError("phone is ban", in.Phone))
		return &zauth_pb.SMSSendVerifyCode_RES{Code: zglobal.Code_Zauth_Sms_Phone_Ban}
	}

	// CD 检查
	rKeyCD := fmt.Sprintf(RK_Sms_PhoneCD, in.Phone)

	if has, e := Redis.Exists(rKeyCD).Result(); e != nil {
		ctx.LogError(e)
		return &zauth_pb.SMSSendVerifyCode_RES{Code: zglobal.Code_Fail}
	} else if has > 0 {
		ctx.LogError(zservice.NewError("phone cd", in.Phone))
		return &zauth_pb.SMSSendVerifyCode_RES{Code: zglobal.Code_Zauth_Sms_Phone_CD}
	}

	verifyCode := zservice.IntToString(zservice.RandomIntRange(100000, 999999))
	if e := SMSSend_aliyun(ctx, &SMSSend_aliyunConfig{
		Phone:        in.Phone,
		VerifyCode:   verifyCode,
		Key:          zservice.Getenv("SMS_KEY"),
		Secret:       zservice.Getenv("SMS_SECRET"),
		TemplateCode: zservice.Getenv("SMS_TEMPLATE_CODE"),
		SignName:     zservice.Getenv("SMS_SIGN_NAME"),
	}); e != nil {
		ctx.LogError(e)
		return &zauth_pb.SMSSendVerifyCode_RES{Code: e.GetCode()}
	}
	ctx.LogInfo("verify code:", in.Phone, verifyCode)

	// 缓存
	// CD
	if e := Redis.SetEX(rKeyCD, " ", time.Minute).Err(); e != nil {
		ctx.LogError(e)
	}

	// 验证码
	if e := Redis.SetEX(fmt.Sprintf(RK_Sms_PhoneCode, in.Phone, verifyCode), " ", 5*time.Minute).Err(); e != nil {
		ctx.LogError(e)
		return &zauth_pb.SMSSendVerifyCode_RES{Code: zglobal.Code_Fail}
	}
	return &zauth_pb.SMSSendVerifyCode_RES{Code: zglobal.Code_SUCC, VerifyCode: fmt.Sprint("****", verifyCode[4:])}
}
