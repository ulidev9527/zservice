package internal

import (
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/zglobal"
)

func Logic_LoginCheck(ctx *zservice.Context, in *zauth_pb.Default_REQ) *zauth_pb.Default_RES {
	if at, e := GetToken(ctx, ctx.AuthToken); e != nil {
		return &zauth_pb.Default_RES{Code: zglobal.Code_Zauth_TokenIsNil}
	} else if at.UID > 0 && at.LoginService == ctx.TraceService {
		return &zauth_pb.Default_RES{Code: zglobal.Code_SUCC}
	} else {
		return &zauth_pb.Default_RES{Code: zglobal.Code_Fail}
	}
}
