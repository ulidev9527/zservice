package internal

import (
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/zglobal"
)

func LoginByAccount(ctx *zservice.Context, in *zauth_pb.LoginByAccount_REQ) *zauth_pb.Default_RES {

	// 验证参数
	if in.Account == "" || in.Password == "" {
		return &zauth_pb.Default_RES{Code: zglobal.Code_ParamsErr}
	}

	// 检查 token 是否登陆
	at, e := GetToken(ctx.AuthToken)
	if e != nil {
		ctx.LogError(e)
		return &zauth_pb.Default_RES{Code: e.GetCode()}
	}

	if at.UID != 0 { // 已登陆的
		if at.LoginTarget == in.LoginTarget {
			return &zauth_pb.Default_RES{Code: zglobal.Code_SUCC}
		}

		return &zauth_pb.Default_RES{Code: zglobal.Code_LoginAgain}
	}

	// 验证账号
	if has, e := HasAccountByLoginName(ctx, in.Account); e != nil {
		ctx.LogError(e)
		return &zauth_pb.Default_RES{Code: e.GetCode()}
	} else if !has {
		return &zauth_pb.Default_RES{Code: zglobal.Code_Zauth_Login_Account_NotFund}
	}

	// 获取账号信息
	acc, e := GetAccountByLoginName(ctx, in.Account)
	if e != nil {
		ctx.LogError(e)
		return &zauth_pb.Default_RES{Code: e.GetCode()}
	} else if !acc.VerifyPass(ctx, in.Password) { // 验证
		return &zauth_pb.Default_RES{Code: zglobal.Code_Zauth_Login_Pass_Err}
	}

	// 设置关联信息
	at.ExpiresSecond = in.Expires
	at.UID = acc.UID
	at.LoginTarget = in.LoginTarget

	if e := at.Save(); e != nil {
		ctx.LogError(e)
		return &zauth_pb.Default_RES{Code: e.GetCode()}
	}

	return &zauth_pb.Default_RES{Code: zglobal.Code_SUCC}
}
