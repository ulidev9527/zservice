package zauth

import (
	"zservice/service/zauth/internal"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/zglobal"
)

// 账号登陆
func LoginByName(ctx *zservice.Context, in *zauth_pb.LoginByUser_REQ) *zauth_pb.Login_RES {

	if in.User == "" || in.Password == "" {
		return &zauth_pb.Login_RES{Code: zglobal.Code_ParamsErr}
	}

	if res, e := func() (*zauth_pb.Login_RES, error) {
		if zauthInitConfig.ServiceName == zservice.GetServiceName() {
			return internal.Logic_LoginByName(ctx, in), nil
		}
		return grpcClient.LoginByUser(ctx, in)
	}(); e != nil {
		return &zauth_pb.Login_RES{Code: zglobal.Code_Fail}
	} else {
		return res
	}
}
