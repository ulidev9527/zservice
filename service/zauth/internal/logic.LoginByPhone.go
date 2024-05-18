package internal

import (
	"zservice/service/zauth/zauth_pb"
	"zservice/zglobal"
	"zservice/zservice"
)

func LoginByPhone(ctx *zservice.Context, in *zauth_pb.LoginByPhone_REQ) *zauth_pb.Default_RES {

	return &zauth_pb.Default_RES{Code: zglobal.Code_SUCC}
}