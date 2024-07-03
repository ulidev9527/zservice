package internal

import (
	"encoding/json"
	"fmt"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/zglobal"
)

func Logic_OrgListGet(ctx *zservice.Context, in *zauth_pb.OrgListGet_REQ) *zauth_pb.OrgInfoList_RES {

	// 限制查询数量
	if in.Size == 0 {
		in.Size = 30
	} else if in.Size > 100 {
		in.Size = 100
	}

	if in.Page == 0 {
		in.Page = 1
	}

	// 查询字符串长度限制在64个字符
	if len(in.Search) > 32 {
		in.Search = in.Search[:32]
	}

	// 查询数据结构
	tabs := []OrgTable{}
	searchStr := fmt.Sprint("%", in.Search, "%")
	if e := Gorm.Model(&OrgTable{}).Where("name like ? OR id like ? OR root_id like ? OR parent_id like ?", searchStr, searchStr, searchStr, searchStr).Order("created_at DESC").Offset(int((in.Page - 1) * in.Size)).Limit(int(in.Size)).Find(&tabs).Error; e != nil {
		ctx.LogError(e)
		return &zauth_pb.OrgInfoList_RES{Code: zglobal.Code_Fail}
	}

	infos := []*zauth_pb.OrgInfo{}
	if e := json.Unmarshal(zservice.JsonMustMarshal(tabs), &infos); e != nil {
		ctx.LogError(e)
		return &zauth_pb.OrgInfoList_RES{Code: zglobal.Code_Fail}
	}

	if len(infos) == 0 {
		return &zauth_pb.OrgInfoList_RES{Code: zglobal.Code_NotFound}
	}

	return &zauth_pb.OrgInfoList_RES{
		Code: zglobal.Code_SUCC,
		List: infos,
	}

}
