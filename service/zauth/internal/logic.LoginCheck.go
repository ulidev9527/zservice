package internal

import (
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/zglobal"
)

func Logic_LoginCheck(ctx *zservice.Context, in *zauth_pb.LogicCheck_REQ) *zauth_pb.Default_RES {

	if at, e := GetToken(ctx, in.Token); e != nil {
		return &zauth_pb.Default_RES{Code: e.GetCode()}
	} else if at.UID > 0 && at.HasLoginService(in.Service) && at.TokenCheck(in.Token, in.Sign) {
		return &zauth_pb.Default_RES{Code: zglobal.Code_SUCC}
	} else {
		return &zauth_pb.Default_RES{Code: zglobal.Code_Fail}
	}
}
