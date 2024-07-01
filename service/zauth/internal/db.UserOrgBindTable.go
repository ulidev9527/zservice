package internal

import (
	"database/sql"
	"fmt"
	"time"
	"zservice/zservice"
	"zservice/zservice/zglobal"

	"gorm.io/gorm"
)

// 账号组织绑定
type UserOrgBindTable struct {
	gorm.Model

	UID     uint32       // 用户ID
	OrgID   uint32       // 组ID
	Expires sql.NullTime // 过期时间
	State   uint32       `gorm:"default:1"` // 状态 0禁用 1开启
}

// 用户和组织绑定
func UserOrgBind(ctx *zservice.Context, uid uint32, orgID uint32, Expires int64, state uint32) (*UserOrgBindTable, *zservice.Error) {

	if tab, e := GetUserOrgBind(ctx, uid, orgID); e != nil {
		if e.GetCode() != zglobal.Code_NotFound {
			return nil, e
		}
	} else {
		// 检查是否更新
		if zservice.MD5String(fmt.Sprint(uid, orgID, Expires, state)) ==
			zservice.MD5String(fmt.Sprint(tab.UID, tab.OrgID, tab.Expires.Time.UnixMilli(), state)) {
			return tab, nil
		} else {
			tab.Expires = sql.NullTime{Time: time.UnixMilli(Expires)}
			tab.State = state
			if e := tab.Save(ctx); e != nil {
				return nil, e
			}
			return tab, nil
		}
	}

	// 创建新数据

	tab := &UserOrgBindTable{
		OrgID:   orgID,
		UID:     uid,
		Expires: sql.NullTime{Time: time.UnixMilli(Expires)},
		State:   state,
	}

	if e := tab.Save(ctx); e != nil {
		return nil, e
	}
	return tab, nil
}

// 用户用户和组织绑定信息
func GetUserOrgBind(ctx *zservice.Context, uid uint32, orgID uint32) (*UserOrgBindTable, *zservice.Error) {

	tab := &UserOrgBindTable{}

	if e := dbhelper.GetTableValue(ctx,
		tab,
		fmt.Sprintf(RK_UserOrgBind_Info, uid, orgID),
		fmt.Sprintf("uid = %d AND org_id = %d", uid, orgID),
	); e != nil {
		return nil, e
	}

	return tab, nil
}

// 是否有账号和组织绑定
func HasUserOrgBindByID(ctx *zservice.Context, uid uint32, orgID uint32) (bool, *zservice.Error) {
	return dbhelper.HasTableValue(ctx, &UserOrgBindTable{}, fmt.Sprintf(RK_UserOrgBind_Info, uid, orgID), fmt.Sprintf("uid = %v and org_id = %v", uid, orgID))
}

// 是否过期
func (z *UserOrgBindTable) IsExpired() bool {
	if z.Expires.Time.IsZero() {
		return false
	}
	return z.Expires.Time.After(time.Now())
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

	rk_info := fmt.Sprintf(RK_UserOrgBind_Info, z.UID, z.OrgID)
	un, e := Redis.Lock(rk_info)
	if e != nil {
		return e
	}
	defer un()

	if e := Mysql.Save(&z).Error; e != nil {
		return zservice.NewError(e)
	}

	// 删缓存
	zservice.Go(func() {
		if e := Redis.Del(rk_info).Err(); e != nil {
			ctx.LogError(e)
		}
	})
	return nil
}
