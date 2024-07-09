package internal

import (
	"strings"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
)

func Logic_UserBan(ctx *zservice.Context, in *zauth_pb.UserBan_REQ) *zauth_pb.Default_RES {

	if in.Uid == 0 || in.Msg == "" {
		ctx.LogError("params error", in)
		return &zauth_pb.Default_RES{Code: zservice.Code_ParamsErr}
	}

	// 添加日志
	if _, e := NewUserBanLogTable(ctx, UserBanLogTable{
		UID:      in.Uid,
		Services: strings.Join(in.Service, ","),
		Msg:      in.Msg,
	}); e != nil {
		ctx.LogError(e)
		return &zauth_pb.Default_RES{Code: e.GetCode()}
	}

	for _, service := range in.Service {

		tab, e := GetOrCreateUserBanTable(ctx, in.Uid, service)
		if e != nil {
			ctx.LogError(e)
			continue
		}
		tab.Expire = zservice.TimeUnixMilli(in.Expire)
		tab.TraceID = ctx.TraceID
		if e := tab.Save(ctx); e != nil {
			ctx.LogError(e)
		}
	}

	return &zauth_pb.Default_RES{Code: zservice.Code_SUCC}
}
