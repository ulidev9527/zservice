package internal

import (
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/zglobal"
)

func Logic_HasAccountID(ctx *zservice.Context, in *zauth_pb.HasAccountID_REQ) *zauth_pb.Default_RES {
	has, e := HasAccountByID(ctx, in.AccountID)

	if e != nil {
		return &zauth_pb.Default_RES{Code: e.GetCode()}
	}
	if has {
		return &zauth_pb.Default_RES{Code: zglobal.Code_SUCC}
	} else {
		return &zauth_pb.Default_RES{Code: zglobal.Code_Zauth_Account_NotFund}
	}
}
