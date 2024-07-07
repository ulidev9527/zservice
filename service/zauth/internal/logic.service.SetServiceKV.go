package internal

import (
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
)

func Logic_SetServiceKV(ctx *zservice.Context, in *zauth_pb.SetServiceKV_REQ) *zauth_pb.Default_RES {
	if tab, e := GetOrCreateServiceKVTable(ctx, in.Service, in.Key); e != nil {
		ctx.LogError(e)
		return &zauth_pb.Default_RES{Code: e.GetCode()}
	} else {
		tab.Value = in.Value
		if e := tab.Save(ctx); e != nil {
			ctx.LogError(e)
			return &zauth_pb.Default_RES{Code: e.GetCode()}
		}
		return &zauth_pb.Default_RES{Code: zservice.Code_SUCC}
	}

}
