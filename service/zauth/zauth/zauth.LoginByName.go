package zauth

import (
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/zglobal"
)

// 账号登陆
func LoginByName(ctx *zservice.Context, in *zauth_pb.LoginByName_REQ) *zauth_pb.Login_RES {

	if in.User == "" || in.Password == "" {
		ctx.LogError("params error")
		return &zauth_pb.Login_RES{Code: zglobal.Code_ParamsErr}
	}

	in.Service = zservice.GetServiceName()
	in.Toekn = ctx.AuthToken
	in.ToeknSign = ctx.AuthTokenSign

	if res, e := grpcClient.LoginByName(ctx, in); e != nil {
		return &zauth_pb.Login_RES{Code: zglobal.Code_Fail}
	} else {
		return res
	}
}
