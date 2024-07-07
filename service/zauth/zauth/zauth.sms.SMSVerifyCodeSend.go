package zauth

import (
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
)

// 发送验证码
func SMSVerifyCodeSend(ctx *zservice.Context, in *zauth_pb.SMSVerifyCodeSend_REQ) *zservice.Error {
	if in.Phone == "" || in.Phone[0] != '+' { // 手机号初步验证
		return zservice.NewError("verify phone fail").SetCode(zservice.Code_ParamsErr)
	}

	if res, e := grpcClient.SMSVerifyCodeSend(ctx, in); e != nil {
		return zservice.NewError("send verify code fail")
	} else if res.Code == zservice.Code_SUCC {
		return nil
	} else {
		return zservice.NewError("send verify code fail").SetCode(res.Code)
	}

}
