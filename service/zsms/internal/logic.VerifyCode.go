package internal

import (
	"zservice/service/zsms/zsms_pb"
	"zservice/zglobal"
	"zservice/zservice"
	"zservice/zservice/ex/redisservice"
)

func VerifyCode(ctx *zservice.Context, in *zsms_pb.VerifyCode_REQ) (code uint32) {
	// 参数检查
	if in.Phone == "" {
		return zglobal.Code_Zsms_Phone_NULL
	}

	if in.Phone[0] != '+' {
		return zglobal.Code_Zsms_Phone_VerifyFail
	}
	if in.Code == "" {
		return zglobal.Code_Zsms_Phone_CodeNull
	}
	if len(in.Code) != 6 {
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

	// 验证
	rk := redisservice.FormatKey(RK_PhoneCode, in.Phone)
	has, e := Redis.Exists(ctx, rk).Result()
	if e != nil {
		ctx.LogError(e)
		return zglobal.Code_ErrorBreakoff
	}
	if has == 0 {
		return zglobal.Code_Zsms_Phone_CodeCacheNull
	}

	codeStr, e := Redis.Get(ctx, rk).Result()
	if e != nil {
		ctx.LogError(e)
		return zglobal.Code_ErrorBreakoff
	}
	isSucc := codeStr == in.Code
	if !isSucc {
		return zglobal.Code_Zsms_Phone_VerifyFail
	}
	// 清除
	e = Redis.Del(ctx, rk).Err()
	if e != nil {
		ctx.LogError(e)
	}
	return zglobal.Code_SUCC

}
