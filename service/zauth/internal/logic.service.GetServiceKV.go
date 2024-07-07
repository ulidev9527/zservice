package internal

import (
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
)

func Logic_GetServiceKV(ctx *zservice.Context, in *zauth_pb.GetServiceKV_REQ) *zauth_pb.GetServiceKV_RES {
	if tab, e := GetOrCreateServiceKVTable(ctx, in.Service, in.Key); e != nil {
		ctx.LogError(e)
		return &zauth_pb.GetServiceKV_RES{Code: e.GetCode()}
	} else {
		return &zauth_pb.GetServiceKV_RES{
			Code:  zservice.Code_SUCC,
			Key:   tab.Key,
			Value: tab.Value,
		}
	}

}
