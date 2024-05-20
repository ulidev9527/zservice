package internal

import (
	"fmt"
	"time"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/zglobal"
)

// 发送验证码
func SMSSendVerifyCode(ctx *zservice.Context, in *zauth_pb.SMSSendVerifyCode_REQ) *zauth_pb.SMSSendVerifyCode_RES {

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
		return &zauth_pb.SMSSendVerifyCode_RES{Code: zglobal.Code_Zauth_Phone_Ban}
	}

	// CD 检查
	rKeyCD := fmt.Sprintf(RK_PhoneCD, in.Phone)

	if has, e := Redis.Exists(rKeyCD).Result(); e != nil {
		ctx.LogError(e)
		return &zauth_pb.SMSSendVerifyCode_RES{Code: zglobal.Code_ErrorBreakoff}
	} else if has > 0 {
		ctx.LogError(zservice.NewError("phone cd", in.Phone))
		return &zauth_pb.SMSSendVerifyCode_RES{Code: zglobal.Code_Zauth_Phone_CD}
	}

	verifyCode := zservice.IntToString(zservice.RandomIntRange(100000, 999999))
	if e := aliyunSMSSend(ctx, &aliyunSMSSendConfig{
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

	// CD
	e := Redis.SetEx(rKeyCD, time.Now().Format(time.RFC3339), time.Duration(zservice.GetenvInt("SMS_CD_DEF"))*time.Second).Err()
	if e != nil {
		ctx.LogError(e) // 已发送，缓存验证码
	}

	// 验证码
	e = Redis.SetEx(fmt.Sprintf(RK_PhoneCode, in.Phone), verifyCode, time.Duration(zservice.GetenvInt("SMS_CODE_CACHE"))*time.Second).Err()
	if e != nil {
		ctx.LogError(e)
		return &zauth_pb.SMSSendVerifyCode_RES{Code: zglobal.Code_ErrorBreakoff}
	}
	ctx.LogInfo("verify code:", verifyCode)
	return &zauth_pb.SMSSendVerifyCode_RES{Code: zglobal.Code_SUCC, VerifyCode: verifyCode}
}
