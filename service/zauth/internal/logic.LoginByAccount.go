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

	// 查询账号
	if has, e := HasAccountByLoginName(ctx, in.Account); e != nil {
		ctx.LogError(e)
		return &zauth_pb.Default_RES{Code: e.GetCode()}
	} else if !has {
		return &zauth_pb.Default_RES{Code: zglobal.Code_Zauth_Login_Account_NotFund}
	}

	return &zauth_pb.Default_RES{Code: zglobal.Code_SUCC}
}
