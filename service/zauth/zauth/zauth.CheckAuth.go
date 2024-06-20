package zauth

import (
	"zservice/service/zauth/internal"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/zglobal"
)

// 检查权限, 没返回错误表示检查成功
func CheckAuth(ctx *zservice.Context, in *zauth_pb.CheckAuth_REQ) *zservice.Error {
	if res, e := func() (*zauth_pb.CheckAuth_RES, error) {
		if zauthInitConfig.ServiceName == zservice.GetServiceName() {
			return internal.Logic_CheckAuth(ctx, in), nil
		}
		return grpcClient.CheckAuth(ctx, in)
	}(); e != nil {
		return zservice.NewError(e)
	} else if res.Code != zglobal.Code_SUCC {
		return zservice.NewError("check auth fail").SetCode(res.Code)
	} else {
		if res.IsTokenRefresh {
			ctx.AuthToken = res.Token
		}
		ctx.AuthID = res.Uid
		return nil
	}
}
