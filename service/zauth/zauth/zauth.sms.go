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
func SendVerifyCode(ctx *zservice.Context, phone string) *zservice.Error {
	if phone == "" || phone[0] != '+' { // 手机号初步验证
		return zservice.NewError("verify phone fail").SetCode(zglobal.Code_ParamsErr)
	}
	req := &zauth_pb.SMSSendVerifyCode_REQ{
		Phone:  phone,
		Serive: zservice.GetServiceName(),
	}

	if res, e := func() (*zauth_pb.SMSSendVerifyCode_RES, error) {
		if zauthInitConfig.ZauthServiceName == zservice.GetServiceName() {
			return internal.Logic_SMSSendVerifyCode(ctx, req), nil
		}
		return grpcClient.SMSSendVerifyCode(context.WithValue(context.Background(), grpcservice.GRPC_contextEX_Middleware_Key, ctx.ContextS2S), req)
	}(); e != nil {
		return zservice.NewError("send verify code fail").SetCode(zglobal.Code_ErrorBreakoff)
	} else if res.Code == zglobal.Code_SUCC {
		return nil
	} else {
		return zservice.NewError("send verify code fail").SetCode(res.Code)
	}

}

// 验证验证码
func VerifyCode(ctx *zservice.Context, phone string, verifyCode string) *zservice.Error {
	if phone == "" || phone[0] != '+' { // 手机号初步验证
		return zservice.NewError("verify phone fail").SetCode(zglobal.Code_ParamsErr)
	}
	if verifyCode == "" || len(verifyCode) != 6 {
		return zservice.NewError("verify code fail").SetCode(zglobal.Code_ParamsErr)
	}

	req := &zauth_pb.SMSVerifyCode_REQ{
		Phone:      phone,
		VerifyCode: verifyCode,
		Serive:     zservice.GetServiceName(),
	}

	if res, e := func() (*zauth_pb.Default_RES, error) {
		if zauthInitConfig.ZauthServiceName == zservice.GetServiceName() {
			return internal.Logic_SMSVerifyCode(ctx, req), nil
		}
		return grpcClient.SMSVerifyCode(context.WithValue(context.Background(), grpcservice.GRPC_contextEX_Middleware_Key, ctx.ContextS2S), req)
	}(); e != nil {
		return zservice.NewError("verify code fail").SetCode(zglobal.Code_ErrorBreakoff)
	} else if res.Code == zglobal.Code_SUCC {
		return nil
	} else {
		return zservice.NewError("verify code fail").SetCode(res.Code)
	}
}
