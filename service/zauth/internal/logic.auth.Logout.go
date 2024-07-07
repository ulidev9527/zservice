package internal

import (
	"fmt"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
)

// 登出
func Logic_Logout(ctx *zservice.Context, in *zauth_pb.Logout_REQ) *zauth_pb.Default_RES {
	if ctx.AuthToken == "" || in.TokenSign == "" {
		ctx.LogError("param error")
		return &zauth_pb.Default_RES{Code: zservice.Code_ParamsErr}
	}

	at, e := GetTokenInfo(ctx, in.Token)
	if e != nil {
		ctx.LogError(e)
		return &zauth_pb.Default_RES{Code: e.GetCode()}
	}
	if !at.TokenCheck(in.TokenSign) {
		ctx.LogError("logout token check fail AT:", at.Token, at.Sign, "IN:", in.Token, in.TokenSign)
		return &zauth_pb.Default_RES{Code: zservice.Code_Fail}
	}

	if at.UID > 0 {

		// 所有关联服务退出
		for _, service := range at.LoginServices {
			rk := fmt.Sprintf(RK_UserLoginServices, at.UID, service)
			if e := Redis.LRem(rk, 0, at.Token).Err(); e != nil {
				ctx.LogError("redis del fail: ", rk, at.Token)
			}
		}

	}

	at.Del(ctx)
	return &zauth_pb.Default_RES{Code: zservice.Code_SUCC}
}
