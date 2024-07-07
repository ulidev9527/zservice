package zauth

import (
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
)

func LoginByToken(ctx *zservice.Context) *zauth_pb.Login_RES {
	in := &zauth_pb.LoginByToken_REQ{
		Service:   zservice.GetServiceName(),
		Toekn:     ctx.AuthToken,
		ToeknSign: ctx.AuthTokenSign,
	}

	if res, e := grpcClient.LoginByToken(ctx, in); e != nil {
		return &zauth_pb.Login_RES{Code: zservice.Code_Fail}
	} else {
		return res
	}
}
