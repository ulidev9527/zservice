package zauth

import (
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/zglobal"
)

// 账号登陆
func LoginByName(ctx *zservice.Context, in *zauth_pb.LoginByName_REQ) *zauth_pb.Login_RES {

	if in.User == "" || in.Password == "" {
		return &zauth_pb.Login_RES{Code: zglobal.Code_ParamsErr}
	}

	if res, e := grpcClient.LoginByName(ctx, in); e != nil {
		return &zauth_pb.Login_RES{Code: zglobal.Code_Fail}
	} else {
		return res
	}
}
