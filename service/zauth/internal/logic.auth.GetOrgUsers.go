package internal

import (
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
)

func Logic_GetOrgUsers(ctx *zservice.Context, in *zauth_pb.GetOrgUsers_REQ) *zauth_pb.GetOrgUsers_RES {

	tabs, e := GetOrgUsers(ctx, in.OrgID, in.Page, 0)

	if e != nil {
		ctx.LogError(e)
		return &zauth_pb.GetOrgUsers_RES{Code: zservice.Code_Fail}
	}

	res := &zauth_pb.GetOrgUsers_RES{Code: zservice.Code_SUCC}

	idList := []uint32{}
	for _, v := range tabs {
		idList = append(idList, v.UID)
	}

	userTabs := []*UserTable{}
	if e := Gorm.Model(&UserTable{}).Where("uid in (?)", idList).Find(&userTabs).Error; e != nil {
		ctx.LogError(e)
		return &zauth_pb.GetOrgUsers_RES{Code: zservice.Code_Fail}
	}

	for _, v := range userTabs {
		res.List = append(res.List, v.ToUserInfo())
	}

	return res
}
