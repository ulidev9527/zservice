package internal

import (
	"fmt"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/zglobal"
)

func Logic_SMSVerifyCode(ctx *zservice.Context, in *zauth_pb.SMSVerifyCode_REQ) *zauth_pb.Default_RES {

	// 参数检查
	if in.Phone == "" || in.Phone[0] != '+' || in.VerifyCode == "" || len(in.VerifyCode) != 6 {
		return &zauth_pb.Default_RES{Code: zglobal.Code_ParamsErr}
	}

	// 封禁检查
	if isBan, e := IsSmsBan(ctx, in.Phone); e != nil {
		ctx.LogError(e)
		return &zauth_pb.Default_RES{Code: e.GetCode()}
	} else if isBan {
		ctx.LogError(zservice.NewError("phone is ban", in.Phone))
		return &zauth_pb.Default_RES{Code: zglobal.Code_Zauth_Phone_Ban}
	}

	// 验证
	rk := fmt.Sprintf(RK_PhoneCode, in.Phone)

	if has, e := Redis.Exists(rk).Result(); e != nil {
		ctx.LogError(e)
		return &zauth_pb.Default_RES{Code: zglobal.Code_ErrorBreakoff}
	} else if has == 0 {
		ctx.LogError(zservice.NewError("phone code not found", in.Phone))
		return &zauth_pb.Default_RES{Code: zglobal.Code_Zauth_Phone_VerifyCodeCacheNull}
	}

	if codeStr, e := Redis.Get(rk).Result(); e != nil {
		ctx.LogError(e)
		return &zauth_pb.Default_RES{Code: zglobal.Code_ErrorBreakoff}
	} else if codeStr != in.VerifyCode {
		return &zauth_pb.Default_RES{Code: zglobal.Code_Zauth_Phone_VerifyCodeErr}
	} else {
		// 清除
		if e := Redis.Del(rk).Err(); e != nil {
			ctx.LogError(e)
		}
		return nil
	}
}
