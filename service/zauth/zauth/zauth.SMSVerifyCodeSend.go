package zauth

import (
	"context"
	"zservice/service/zauth/internal"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/ex/grpcservice"
	"zservice/zservice/zglobal"
)

// 发送验证码
func SMSVerifyCodeSend(ctx *zservice.Context, in *zauth_pb.SMSVerifyCodeSend_REQ) *zservice.Error {
	if in.Phone == "" || in.Phone[0] != '+' { // 手机号初步验证
		return zservice.NewError("verify phone fail").SetCode(zglobal.Code_ParamsErr)
	}

	if res, e := func() (*zauth_pb.SMSSendVerifyCode_RES, error) {
		if zauthInitConfig.ZauthServiceName == zservice.GetServiceName() {
			return internal.Logic_SMSVerifyCodeSend(ctx, in), nil
		}
		return grpcClient.SMSVerifyCodeSend(context.WithValue(context.Background(), grpcservice.GRPC_contextEX_Middleware_Key, ctx.ContextS2S), in)
	}(); e != nil {
		return zservice.NewError("send verify code fail").SetCode(zglobal.Code_ErrorBreakoff)
	} else if res.Code == zglobal.Code_SUCC {
		return nil
	} else {
		return zservice.NewError("send verify code fail").SetCode(res.Code)
	}

}
