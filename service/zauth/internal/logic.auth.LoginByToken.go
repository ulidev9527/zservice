package internal

import (
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
)

func Logic_LoginByToken(ctx *zservice.Context, in *zauth_pb.LoginByToken_REQ) *zauth_pb.Login_RES {
	// 参数验证
	if in.Toekn == "" {
		ctx.LogError("param error")
		return &zauth_pb.Login_RES{Code: zservice.Code_ParamsErr}
	}

	// 检查 token 是否登陆
	if res := Logic_LoginCheck(ctx, &zauth_pb.LoginCheck_REQ{
		Token:     in.Toekn,
		TokenSign: in.ToeknSign,
		Service:   in.Service,
	}); res.Code == zservice.Code_SUCC {
		return &zauth_pb.Login_RES{Code: zservice.Code_SUCC, UserInfo: res.UserInfo}
	}

	// 获取token
	at, e := GetTokenInfo(ctx, in.Toekn)
	if e != nil {
		ctx.LogError(e)
		return &zauth_pb.Login_RES{Code: e.GetCode()}
	}

	if at.UID == 0 {
		ctx.LogError("token is not login other service", in)
		return &zauth_pb.Login_RES{Code: zservice.Code_Reject}
	}

	// 获取账号信息/验证账号状态
	user, e := GetUserByUID(ctx, at.UID)
	if e != nil {
		ctx.LogError(e)
		return &zauth_pb.Login_RES{Code: e.GetCode()}
	}

	// 登录
	if e := TokenLogin(ctx, struct {
		Service string
		Expires uint32
	}{
		Service: in.Service,
	}, at, nil); e != nil {
		ctx.LogError(e)
		return &zauth_pb.Login_RES{Code: e.GetCode()}
	}

	return &zauth_pb.Login_RES{Code: zservice.Code_SUCC, UserInfo: user.ToUserInfo()}
}
