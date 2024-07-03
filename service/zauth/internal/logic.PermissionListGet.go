package internal

import (
	"encoding/json"
	"fmt"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/zglobal"
)

// 获取权限列表
func Logic_PermissionListGet(ctx *zservice.Context, in *zauth_pb.PermissionListGet_REQ) *zauth_pb.PermissionInfoList_RES {

	// 限制查询数量
	if in.Size == 0 {
		in.Size = 30
	} else if in.Size > 100 {
		in.Size = 100
	}

	// 查询字符串长度限制在64个字符
	if len(in.Search) > 32 {
		in.Search = in.Search[:32]
	}

	// 查询数据结构
	tabs := []PermissionTable{}
	searchStr := fmt.Sprint("%", in.Search, "%")
	if e := Gorm.Model(&PermissionTable{}).Where("name like ? OR id like ? OR service like ? OR action like ? OR path like ?", searchStr, searchStr, searchStr, searchStr, searchStr).Order("created_at DESC").Offset(int((in.Page - 1) * in.Size)).Limit(int(in.Size)).Find(&tabs).Error; e != nil {
		ctx.LogError(e)
		return &zauth_pb.PermissionInfoList_RES{
			Code: zglobal.Code_Fail,
		}
	}

	infos := []*zauth_pb.PermissionInfo{}
	if e := json.Unmarshal(zservice.JsonMustMarshal(tabs), &infos); e != nil {
		ctx.LogError(e)
		return &zauth_pb.PermissionInfoList_RES{
			Code: zglobal.Code_Fail,
		}
	}

	if len(infos) == 0 {
		return &zauth_pb.PermissionInfoList_RES{
			Code: zglobal.Code_NotFound,
		}
	}

	return &zauth_pb.PermissionInfoList_RES{
		Code: zglobal.Code_SUCC,
		List: infos,
	}

}
