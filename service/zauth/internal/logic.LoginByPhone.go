package internal

import (
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/zglobal"
)

// 手机登陆
func Logic_LoginByPhone(ctx *zservice.Context, in *zauth_pb.LoginByPhone_REQ) *zauth_pb.Default_RES {
	if in.Phone == "" || in.VerifyCode == "" || len(in.VerifyCode) != 6 || in.LoginService == "" {
		return &zauth_pb.Default_RES{Code: zglobal.Code_ParamsErr}
	}

	// 检查 token 是否登陆
	at, e := GetToken(ctx.AuthToken)
	if e != nil {
		ctx.LogError(e)
		return &zauth_pb.Default_RES{Code: e.GetCode()}
	}

	if at.UID != 0 { // 已登陆的
		if at.LoginService == in.LoginService {
			return &zauth_pb.Default_RES{Code: zglobal.Code_SUCC}
		}

		return &zauth_pb.Default_RES{Code: zglobal.Code_LoginAgain}
	}

	// 验证手机号
	verifyRes := Logic_SMSVerifyCode(ctx, &zauth_pb.SMSVerifyCode_REQ{
		Phone:      in.Phone,
		VerifyCode: in.VerifyCode,
	})
	if verifyRes.Code != zglobal.Code_SUCC {
		return verifyRes
	}

	// 获取账号信息
	acc, e := GetAccountByPhone(ctx, in.Phone)
	if e != nil {
		if e.GetCode() != zglobal.Code_Zauth_Account_NotFund { // 其他错误
			ctx.LogError(e)
			return &zauth_pb.Default_RES{Code: e.GetCode()}
		} else { // 未找到账号, 进行创建
			acc, e = CreateAccount(ctx)
			if e != nil {
				ctx.LogError(e)
				return &zauth_pb.Default_RES{Code: e.GetCode()}
			}

			acc.Phone = in.Phone
			if e := acc.Save(ctx); e != nil {
				ctx.LogError(e)
				return &zauth_pb.Default_RES{Code: e.GetCode()}
			}
		}
	}

	// 设置关联信息
	at.ExpiresSecond = in.Expires
	at.UID = acc.UID
	at.LoginService = in.LoginService

	// 存储
	if e := at.Save(); e != nil {
		ctx.LogError(e)
		return &zauth_pb.Default_RES{Code: e.GetCode()}
	}

	return &zauth_pb.Default_RES{Code: zglobal.Code_SUCC}

}
