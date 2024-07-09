package zauth

import (
	"zservice/service/zauth/zauth_pb"
)

type __serviceInfo struct {
	serviceRegistRES *zauth_pb.ServiceRegist_RES
}

// 服务信息
var serviceInfo = &__serviceInfo{}

// 获取服务的组织ID
func GetServiceOrgID() uint32 {
	return serviceInfo.serviceRegistRES.OrgInfo.OrgID
}
