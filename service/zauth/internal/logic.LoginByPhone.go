package internal

import (
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/zglobal"
)

// 手机登陆
func Logic_LoginByPhone(ctx *zservice.Context, in *zauth_pb.LoginByPhone_REQ) *zauth_pb.Login_RES {
	if in.Phone == "" || in.VerifyCode == "" || len(in.VerifyCode) != 6 || in.Service == "" || in.Toekn == "" || in.ToeknSign == "" {
		ctx.LogError("param error", in.Phone, in.VerifyCode, in.Service, in.Toekn, in.ToeknSign)
		return &zauth_pb.Login_RES{Code: zglobal.Code_ParamsErr}
	}

	// 检查 token 是否登陆
	if res := Logic_LoginCheck(ctx, &zauth_pb.LoginCheck_REQ{
		Token:     in.Toekn,
		TokenSign: in.ToeknSign,
		Service:   in.Service,
	}); res.Code == zglobal.Code_SUCC {
		return &zauth_pb.Login_RES{Code: zglobal.Code_SUCC, UserInfo: res.UserInfo}
	}

	// 获取token
	at, e := GetToken(ctx, in.Toekn)
	if e != nil {
		ctx.LogError(e)
		return &zauth_pb.Login_RES{Code: e.GetCode()}
	}

	// 验证手机号
	if verifyRes := Logic_SMSVerifyCodeVerify(ctx, &zauth_pb.SMSVerifyCodeVerify_REQ{
		Phone:      in.Phone,
		VerifyCode: in.VerifyCode,
	}); verifyRes.Code != zglobal.Code_SUCC {
		ctx.LogError("phone verify code fail", in.Phone, in.VerifyCode)
		return &zauth_pb.Login_RES{Code: verifyRes.GetCode()}
	}

	// 获取账号信息/ 验证账号状态
	user, e := GetUserByPhone(ctx, in.Phone)
	if e != nil {
		if e.GetCode() != zglobal.Code_NotFound { // 其他错误
			ctx.LogError(e)
			return &zauth_pb.Login_RES{Code: e.GetCode()}
		} else { // 未找到账号, 进行创建
			user, e = CreateUser(ctx)
			if e != nil {
				ctx.LogError(e)
				return &zauth_pb.Login_RES{Code: e.GetCode()}
			}

			user.Phone = in.Phone
			if e := user.Save(ctx); e != nil {
				ctx.LogError(e)
				return &zauth_pb.Login_RES{Code: e.GetCode()}
			}
		}
	} else if user.State == 0 {
		ctx.LogError("login limit", user.UID)
		return &zauth_pb.Login_RES{Code: zglobal.Code_Limit}
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

	return &zauth_pb.Login_RES{Code: zglobal.Code_SUCC, UserInfo: &zauth_pb.UserInfo{
		Uid:       user.UID,
		LoginName: user.LoginName,
		Phone:     user.Phone,
		State:     user.State,
	}}

}
