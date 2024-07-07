package zauth

import (
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
)

// 手机号登陆
func LoginByPhone(ctx *zservice.Context, in *zauth_pb.LoginByPhone_REQ) *zauth_pb.Login_RES {

	if in.Phone == "" || in.VerifyCode == "" {
		ctx.LogError("params error")
		return &zauth_pb.Login_RES{Code: zservice.Code_ParamsErr}
	}

	in.Service = zservice.GetServiceName()
	in.Toekn = ctx.AuthToken
	in.ToeknSign = ctx.AuthTokenSign

	if res, e := grpcClient.LoginByPhone(ctx, in); e != nil {
		return &zauth_pb.Login_RES{Code: zservice.Code_Fail}
	} else {
		return res
	}

}
