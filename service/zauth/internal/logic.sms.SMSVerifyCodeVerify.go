package internal

import (
	"fmt"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
)

func Logic_SMSVerifyCodeVerify(ctx *zservice.Context, in *zauth_pb.SMSVerifyCodeVerify_REQ) *zauth_pb.Default_RES {

	// 参数检查
	if in.Phone == "" || in.Phone[0] != '+' || in.VerifyCode == "" || len(in.VerifyCode) != 6 {
		return &zauth_pb.Default_RES{Code: zservice.Code_ParamsErr}
	}

	// 封禁检查
	if isBan, e := IsSmsBan(ctx, in.Phone); e != nil {
		ctx.LogError(e)
		return &zauth_pb.Default_RES{Code: e.GetCode()}
	} else if isBan {
		ctx.LogError(zservice.NewError("phone is ban", in.Phone))
		return &zauth_pb.Default_RES{Code: zservice.Code_Zauth_Sms_Phone_Ban}
	}

	// 验证
	rk := fmt.Sprintf(RK_Sms_PhoneCode, in.Phone, in.VerifyCode)

	if has, e := Redis.Exists(rk).Result(); e != nil {
		ctx.LogError(e)
		return &zauth_pb.Default_RES{Code: zservice.Code_Fail}
	} else if has == 0 {
		return &zauth_pb.Default_RES{Code: zservice.Code_Zauth_Sms_Phone_VerifyFail}
	}

	if e := Redis.Del(rk).Err(); e != nil {
		ctx.LogError(e)
	}
	return &zauth_pb.Default_RES{Code: zservice.Code_SUCC}
}
