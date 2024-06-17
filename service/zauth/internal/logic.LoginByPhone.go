package internal

import (
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/zglobal"
)

// 手机登陆
func Logic_LoginByPhone(ctx *zservice.Context, in *zauth_pb.LoginByPhone_REQ) *zauth_pb.Login_RES {
	if in.Phone == "" || in.VerifyCode == "" || len(in.VerifyCode) != 6 || ctx.TraceService == "" {
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

	// 验证手机号
	verifyRes := Logic_SMSVerifyCodeVerify(ctx, &zauth_pb.SMSVerifyCodeVerify_REQ{
		Phone:      in.Phone,
		VerifyCode: in.VerifyCode,
	})
	if verifyRes.Code != zglobal.Code_SUCC {
		return &zauth_pb.Login_RES{Code: verifyRes.Code}
	}

	// 获取账号信息
	acc, e := GetUserByPhone(ctx, in.Phone)
	if e != nil {
		if e.GetCode() != zglobal.Code_Zauth_User_NotFund { // 其他错误
			ctx.LogError(e)
			return &zauth_pb.Login_RES{Code: e.GetCode()}
		} else { // 未找到账号, 进行创建
			acc, e = CreateUser(ctx)
			if e != nil {
				ctx.LogError(e)
				return &zauth_pb.Login_RES{Code: e.GetCode()}
			}

			acc.Phone = in.Phone
			if e := acc.Save(ctx); e != nil {
				ctx.LogError(e)
				return &zauth_pb.Login_RES{Code: e.GetCode()}
			}
		}
	}

	// 设置关联信息
	at.ExpiresSecond = in.Expires
	at.UID = acc.UID
	at.LoginService = ctx.TraceService

	// 存储
	if e := at.Save(); e != nil {
		ctx.LogError(e)
		return &zauth_pb.Login_RES{Code: e.GetCode()}
	}

	return &zauth_pb.Login_RES{Code: zglobal.Code_SUCC, Uid: at.UID}

}
