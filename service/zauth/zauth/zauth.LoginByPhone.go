package zauth

import (
	"zservice/service/zauth/internal"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/zglobal"
)

// 手机号登陆
func LoginByPhone(ctx *zservice.Context, in *zauth_pb.LoginByPhone_REQ) *zauth_pb.Login_RES {

	if in.Phone == "" || in.VerifyCode == "" {
		return &zauth_pb.Login_RES{Code: zglobal.Code_ParamsErr}
	}

	if res, e := func() (*zauth_pb.Login_RES, error) {
		if zauthInitConfig.ServiceName == zservice.GetServiceName() {
			return internal.Logic_LoginByPhone(ctx, in), nil
		}
		return grpcClient.LoginByPhone(ctx, in)
	}(); e != nil {
		return &zauth_pb.Login_RES{Code: zglobal.Code_Fail}
	} else {
		return res
	}

}
