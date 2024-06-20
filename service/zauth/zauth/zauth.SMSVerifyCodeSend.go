package zauth

import (
	"zservice/service/zauth/internal"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/zglobal"
)

// 发送验证码
func SMSVerifyCodeSend(ctx *zservice.Context, in *zauth_pb.SMSVerifyCodeSend_REQ) *zservice.Error {
	if in.Phone == "" || in.Phone[0] != '+' { // 手机号初步验证
		return zservice.NewError("verify phone fail").SetCode(zglobal.Code_ParamsErr)
	}

	if res, e := func() (*zauth_pb.SMSSendVerifyCode_RES, error) {
		if zauthInitConfig.ServiceName == zservice.GetServiceName() {
			return internal.Logic_SMSVerifyCodeSend(ctx, in), nil
		}
		return grpcClient.SMSVerifyCodeSend(ctx, in)
	}(); e != nil {
		return zservice.NewError("send verify code fail")
	} else if res.Code == zglobal.Code_SUCC {
		return nil
	} else {
		return zservice.NewError("send verify code fail").SetCode(res.Code)
	}

}
