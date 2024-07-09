package internal

import (
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
)

func Logic_LoginByName(ctx *zservice.Context, in *zauth_pb.LoginByName_REQ) *zauth_pb.Login_RES {

	// 验证参数
	if in.User == "" || in.Password == "" || in.Service == "" || in.Toekn == "" || in.ToeknSign == "" {
		ctx.LogError("param error", in.User, in.Password, in.Service, in.Toekn, in.ToeknSign)
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

	// 获取账号信息/验证账号状态
	user, e := GetUserByLoginName(ctx, in.User)
	if e != nil {
		if e.GetCode() != zservice.Code_NotFound {
			ctx.LogError(e)
			return &zauth_pb.Login_RES{Code: e.GetCode()}
		}
		if !in.IsAutoCreate { // 是否自动创建
			ctx.LogError(e)
			return &zauth_pb.Login_RES{Code: e.GetCode()}
		}

		// 创建
		user, e = CreateUser(ctx)
		if e != nil {
			ctx.LogError(e)
			return &zauth_pb.Login_RES{Code: e.GetCode()}
		}

		if e := user.AddLoginNameAndPassword(ctx, in.User, in.Password); e != nil {
			ctx.LogError(e)
			return &zauth_pb.Login_RES{Code: e.GetCode()}
		}

	} else if isBan, e := user.IsBan(ctx, in.Service); e != nil { // 限制登陆
		ctx.LogError(e)
		return &zauth_pb.Login_RES{Code: e.GetCode()}
	} else if isBan {
		ctx.LogError("login limit", user.UID)
		return &zauth_pb.Login_RES{Code: zservice.Code_Limit}
	} else if !user.VerifyPass(ctx, in.Password) { // 验证
		ctx.LogError("pasword err", in.User, in.Password)
		return &zauth_pb.Login_RES{Code: zservice.Code_Reject}
	}

	// 登录
	if e := TokenLogin(ctx, struct {
		Service string
		Expires uint32
	}{
		Service: in.Service,
	}, at, user); e != nil {
		ctx.LogError(e)
		return &zauth_pb.Login_RES{Code: e.GetCode()}
	}

	return &zauth_pb.Login_RES{Code: zservice.Code_SUCC, UserInfo: user.ToUserInfo()}
}
