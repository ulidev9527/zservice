package zauth

import (
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/zglobal"
)

// 验证验证码
func SMSVerifyCodeVerifyt(ctx *zservice.Context, phone string, verifyCode string) *zservice.Error {
	if phone == "" || phone[0] != '+' { // 手机号初步验证
		return zservice.NewError("verify phone fail").SetCode(zglobal.Code_ParamsErr)
	}
	if verifyCode == "" || len(verifyCode) != 6 {
		return zservice.NewError("verify code fail").SetCode(zglobal.Code_ParamsErr)
	}

	if res, e := grpcClient.SMSVerifyCodeVerify(ctx, &zauth_pb.SMSVerifyCodeVerify_REQ{
		Phone:      phone,
		VerifyCode: verifyCode,
	}); e != nil {
		return zservice.NewError("verify code fail")
	} else if res.Code == zglobal.Code_SUCC {
		return nil
	} else {
		return zservice.NewError("verify code fail").SetCode(res.Code)
	}
}
