package internal

import (
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/zglobal"
)

// 登出
func Logic_Logout(ctx *zservice.Context, in *zauth_pb.Default_REQ) *zauth_pb.Default_RES {
	if ctx.AuthToken == "" {
		return &zauth_pb.Default_RES{Code: zglobal.Code_ParamsErr}
	}

	at, e := GetToken(ctx.AuthToken)
	if e != nil {
		ctx.LogError(e)
		return &zauth_pb.Default_RES{Code: e.GetCode()}
	}

	if e := at.Del(); e != nil {
		ctx.LogError(e)
		return &zauth_pb.Default_RES{Code: e.GetCode()}
	}
	return &zauth_pb.Default_RES{Code: zglobal.Code_SUCC}
}
