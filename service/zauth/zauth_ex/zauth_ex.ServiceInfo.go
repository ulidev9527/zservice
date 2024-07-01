package zauth_ex

import (
	"zservice/service/zauth/internal"
	"zservice/service/zauth/zauth"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/zglobal"
)

type __serviceInfo struct {
	serviceRegistRES *zauth_pb.ServiceRegist_RES
}

var ServiceInfo = &__serviceInfo{}

// 注册到 zauth 服务
func (si *__serviceInfo) Regist(ctx *zservice.Context, in *zauth_pb.ServiceRegist_REQ, isZauthSelf ...bool) {

	isSelf := false
	if len(isZauthSelf) > 0 {
		isSelf = isZauthSelf[0]
	}

	in.Service = zservice.GetServiceName()

	if res, e := func() (*zauth_pb.ServiceRegist_RES, error) {
		if isSelf {
			return internal.Logic_ServiceRegist(ctx, in), nil
		}
		return zauth.GetGrpcClient().ServiceRegist(ctx, in)
	}(); e != nil {
		ctx.LogPanic(e)
	} else if res.Code != zglobal.Code_SUCC {
		ctx.LogPanic(res)
	} else {
		si.serviceRegistRES = res
	}
}

func (si *__serviceInfo) GetOrgID() uint32 {
	return si.serviceRegistRES.OrgInfo.OrgID
}
