package internal

import (
	"fmt"
	"time"
	"zservice/service/zsms/zsms_pb"
	"zservice/zglobal"
	"zservice/zservice"
)

// 发送验证码
func SendVerifyCode(ctx *zservice.Context, in *zsms_pb.SendVerifyCode_REQ) (code uint32) {

	// 参数检查
	if in.Phone == "" {
		return zglobal.Code_Zsms_Phone_NULL
	}

	if in.Phone[0] != '+' {
		return zglobal.Code_Zsms_Phone_VerifyFail
	}

	// 封禁检查
	isBan, e := IsSmsBan(ctx, in.Phone)
	if e != nil {
		ctx.LogError(e)
		ze, ok := e.(*zservice.Error)
		if ok {
			return ze.GetCode()
		}
		return zglobal.Code_ErrorBreakoff
	}
	if isBan {
		return zglobal.Code_Zsms_Phone_Ban
	}

	// CD 检查
	rKeyCD := fmt.Sprintf(RK_PhoneCD, in.Phone)
	has, e := Redis.Exists(ctx, rKeyCD).Result()
	if e != nil {
		ctx.LogError(e)
		return zglobal.Code_ErrorBreakoff
	}
	if has > 0 {
		return zglobal.Code_Zsms_Phone_CD
	}

	vCode := zservice.Convert_IntToString(zservice.RandomIntRange(100000, 999999))
	code = aliyunSMSSend(ctx, &aliyunSMSSendConfig{
		Phone:        in.Phone,
		VerifyCode:   vCode,
		Key:          zservice.Getenv("SMS_KEY"),
		Secret:       zservice.Getenv("SMS_SECRET"),
		TemplateCode: zservice.Getenv("SMS_TEMPLATE_CODE"),
		SignName:     zservice.Getenv("SMS_SIGN_NAME"),
	})

	// 缓存
	if code == zglobal.Code_SUCC {

		// CD
		e := Redis.Set(ctx, rKeyCD, time.Now().Format(time.RFC3339), time.Duration(zservice.GetenvInt("SMS_CD_DEF"))*time.Second).Err()
		if e != nil {
			ctx.LogError(e)
		}

		// 验证码
		e = Redis.Set(ctx, fmt.Sprintf(RK_PhoneCode, in.Phone), code, time.Duration(zservice.GetenvInt("SMS_CODE_CACHE"))*time.Second).Err()
		if e != nil {
			ctx.LogError(e)
			return zglobal.Code_ErrorBreakoff
		}
		ctx.LogInfo("verify code:", vCode)
	}
	return code
}
