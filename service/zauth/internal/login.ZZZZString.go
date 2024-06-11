package internal

import (
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/zglobal"
)

func Logic_ZZZZString(ctx *zservice.Context, in *zauth_pb.ZZZZString_REQ) *zauth_pb.Default_RES {
	return &zauth_pb.Default_RES{
		Code: zglobal.Code_SUCC,
	}
}
