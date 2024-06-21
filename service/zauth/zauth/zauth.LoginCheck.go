package zauth

import (
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/zglobal"
)

// 登陆检查
func LoginCheck(ctx *zservice.Context) bool {
	if ctx.AuthToken == "" {
		return false
	}

	if res, e := grpcClient.LoginCheck(ctx, &zauth_pb.Default_REQ{}); e != nil {
		ctx.LogError(e)
		return false
	} else if res.Code == zglobal.Code_SUCC {
		return true
	} else {
		return false
	}
}
