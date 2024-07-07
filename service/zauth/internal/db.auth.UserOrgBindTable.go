package internal

import (
	"fmt"
	"time"
	"zservice/zservice"
	"zservice/zserviceex/dbservice"

	"gorm.io/gorm"
)

// 账号组织绑定
type UserOrgBindTable struct {
	gorm.Model

	UID     uint32        // 用户ID
	OrgID   uint32        // 组ID
	Expires zservice.Time // 过期时间
	State   uint32        `gorm:"default:1"` // 状态 0禁用 1开启
}

// 用户和组织绑定
func UserOrgBind(ctx *zservice.Context, uid uint32, orgID uint32, expires int64, state uint32) (*UserOrgBindTable, *zservice.Error) {

	if tab, e := GetUserOrgBind(ctx, uid, orgID); e != nil {
		if e.GetCode() != zservice.Code_NotFound {
			return nil, e.AddCaller()
		}
	} else {
		// 检查是否更新
		if zservice.MD5String(fmt.Sprint(uid, orgID, expires, state)) ==
			zservice.MD5String(fmt.Sprint(tab.UID, tab.OrgID, tab.Expires.UnixMilli(), state)) {
			return tab, nil
		} else {
			tab.Expires = zservice.NewTime(time.UnixMilli(expires))
			tab.State = state
			if e := tab.Save(ctx); e != nil {
				return nil, e.AddCaller()
			}
			return tab, nil
		}
	}

	// 创建新数据

	tab := &UserOrgBindTable{
		OrgID:   orgID,
		UID:     uid,
		Expires: zservice.NewTime(time.UnixMilli(expires)),
		State:   state,
	}

	if e := tab.Save(ctx); e != nil {
		return nil, e.AddCaller()
	}
	return tab, nil
}

// 用户用户和组织绑定信息
func GetUserOrgBind(ctx *zservice.Context, uid uint32, orgID uint32) (*UserOrgBindTable, *zservice.Error) {

	tab := &UserOrgBindTable{}

	if e := DBService.GetTableFirst(ctx, dbservice.GetTableValueOption{
		Tab:      tab,
		RK:       fmt.Sprintf(RK_UserOrgBind_Info, uid, orgID),
		SQLConds: []any{"uid = ? AND org_id = ?", uid, orgID},
	}); e != nil {
		return nil, e.AddCaller()
	}

	return tab, nil
}

// 是否有账号和组织绑定
func HasUserOrgBindByID(ctx *zservice.Context, uid uint32, orgID uint32) (bool, *zservice.Error) {
	return DBService.HasTableValue(ctx, dbservice.HasTableValueOption{Tab: &UserOrgBindTable{}, RK: fmt.Sprintf(RK_UserOrgBind_Info, uid, orgID), SQLConds: []any{"uid = ? AND org_id = ?", uid, orgID}})
}

// 是否过期
func (z *UserOrgBindTable) IsExpired() bool {
	if z.Expires.IsZero() {
		return false
	}
	return z.Expires.Before(time.Now())
}

// 是否启动
func (z *UserOrgBindTable) IsAllow() bool {
	if z.IsExpired() {
		return false
	}
	return z.State == 1
}

// 存储
func (z *UserOrgBindTable) Save(ctx *zservice.Context) *zservice.Error {
	if e := DBService.SaveTableValue(ctx, z, fmt.Sprintf(RK_UserOrgBind_Info, z.UID, z.OrgID)); e != nil {
		return e.AddCaller()
	}
	return nil
}
