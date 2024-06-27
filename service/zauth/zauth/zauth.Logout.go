package zauth

import (
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/zglobal"
)

func Logout(ctx *zservice.Context) *zauth_pb.Default_RES {

	if res, e := grpcClient.Logout(ctx, &zauth_pb.Logout_REQ{
		Token:     ctx.AuthToken,
		TokenSign: ctx.AuthTokenSign,
	}); e != nil {
		ctx.LogPanic(e)
	} else if res.Code != zglobal.Code_SUCC {
		ctx.LogPanic(res)
	}
	return &zauth_pb.Default_RES{Code: zglobal.Code_SUCC}
}
