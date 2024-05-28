package internal

import (
	"encoding/json"
	"fmt"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/zglobal"
)

// 获取权限列表
func Logic_GetPermissionList(ctx *zservice.Context, in *zauth_pb.GetPermissionList_REQ) *zauth_pb.GetPermissionList_RES {

	// 页码处理
	if in.Page < 0 {
		in.Page = 0
	}

	// 限制查询数量
	if in.Size <= 0 {
		in.Size = 30
	} else if in.Size > 100 {
		in.Size = 100
	}

	// 查询字符串长度限制在64个字符
	if len(in.Search) > 32 {
		in.Search = in.Search[:32]
	}

	// 查询数据结构
	tabs := []ZauthPermissionTable{}
	searchStr := fmt.Sprint("%", in.Search, "%")
	if e := Mysql.Model(&ZauthPermissionTable{}).Where("name like ? OR permission_id like ? OR service like ? OR action like ? OR path like ?", searchStr, searchStr, searchStr, searchStr, searchStr).Order("path desc").Offset(int((in.Page - 1) * in.Size)).Limit(int(in.Size)).Find(&tabs).Error; e != nil {
		ctx.LogError(e)
		return &zauth_pb.GetPermissionList_RES{
			Code: zglobal.Code_ErrorBreakoff,
		}
	}

	infos := []*zauth_pb.PermissionInfo{}
	if e := json.Unmarshal(zservice.JsonMustMarshal(tabs), &infos); e != nil {
		ctx.LogError(e)
		return &zauth_pb.GetPermissionList_RES{
			Code: zglobal.Code_ErrorBreakoff,
		}
	}

	return &zauth_pb.GetPermissionList_RES{
		Code: zglobal.Code_SUCC,
		List: infos,
	}

}
