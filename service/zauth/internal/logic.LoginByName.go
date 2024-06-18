package internal

import (
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/zglobal"
)

func Logic_LoginByName(ctx *zservice.Context, in *zauth_pb.LoginByUser_REQ) *zauth_pb.Login_RES {

	// 验证参数
	if in.User == "" || in.Password == "" {
		return &zauth_pb.Login_RES{Code: zglobal.Code_ParamsErr}
	}

	// 检查 token 是否登陆
	at, e := GetToken(ctx.AuthToken)
	if e != nil {
		ctx.LogError(e)
		return &zauth_pb.Login_RES{Code: e.GetCode()}
	}

	if at.UID != 0 { // 已登陆的
		if at.LoginService == ctx.TraceService {
			return &zauth_pb.Login_RES{Code: zglobal.Code_SUCC, Uid: at.UID}
		}

		return &zauth_pb.Login_RES{Code: zglobal.Code_LoginAgain}
	}

	// 验证账号
	if has, e := HasUserByLoginName(ctx, in.User); e != nil {
		ctx.LogError(e)
		return &zauth_pb.Login_RES{Code: e.GetCode()}
	} else if !has {
		return &zauth_pb.Login_RES{Code: zglobal.Code_Zauth_Login_User_NotFund}
	}

	// 获取账号信息
	acc, e := GetUserByLoginName(ctx, in.User)
	if e != nil {
		ctx.LogError(e)
		return &zauth_pb.Login_RES{Code: e.GetCode()}
	} else if !acc.VerifyPass(ctx, in.Password) { // 验证
		return &zauth_pb.Login_RES{Code: zglobal.Code_Zauth_Login_Pass_Err}
	}

	// 设置关联信息
	at.ExpiresSecond = in.Expires
	at.UID = acc.UID
	at.LoginService = ctx.TraceService

	if e := at.Save(); e != nil {
		ctx.LogError(e)
		return &zauth_pb.Login_RES{Code: e.GetCode()}
	}

	return &zauth_pb.Login_RES{Code: zglobal.Code_SUCC, Uid: acc.UID}
}