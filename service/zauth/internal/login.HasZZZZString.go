package internal

import (
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/zglobal"
)

func Logic_HasZZZZString(ctx *zservice.Context, in *zauth_pb.HasZZZZString_REQ) *zauth_pb.Default_RES {
	if ZZZZString.Has(ctx, in.Str) {
		return &zauth_pb.Default_RES{Code: zglobal.Code_ZZZZ}
	} else {
		return &zauth_pb.Default_RES{Code: zglobal.Code_SUCC}
	}
}
