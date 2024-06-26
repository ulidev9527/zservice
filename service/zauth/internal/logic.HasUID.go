package internal

import (
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/zglobal"
)

func Logic_HasUID(ctx *zservice.Context, in *zauth_pb.HasUID_REQ) *zauth_pb.Default_RES {
	has, e := HasUserByID(ctx, in.Uid)

	if e != nil {
		return &zauth_pb.Default_RES{Code: e.GetCode()}
	}
	if has {
		return &zauth_pb.Default_RES{Code: zglobal.Code_SUCC}
	} else {
		return &zauth_pb.Default_RES{Code: zglobal.Code_Zauth_User_NotFund}
	}
}
