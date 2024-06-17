package internal

import (
	"fmt"
	"time"
	"zservice/zservice"
	"zservice/zservice/ex/gormservice"
	"zservice/zservice/zglobal"
)

// 账号组织绑定
type UserOrgBindTable struct {
	gormservice.Model

	OrgID   uint32 // 组ID
	UID     uint32 // 用户ID
	Expires uint64 // 过期时间
	State   uint32 `gorm:"default:1"` // 状态 0禁用 1开启
}

// 加入组织
func UserJoinOrg(ctx *zservice.Context, uid uint32, orgID uint32, Expires uint64) (*UserOrgBindTable, *zservice.Error) {
	// 验证参数是否正确
	if has, e := HasOrgByID(ctx, orgID); e != nil {
		return nil, e
	} else if !has {
		return nil, zservice.NewError("org not found:", orgID).SetCode(zglobal.Code_Zauth_Org_NotFund)
	}

	if has, e := HasUserByID(ctx, uid); e != nil {
		return nil, e
	} else if !has {
		return nil, zservice.NewError("user not found:", uid).SetCode(zglobal.Code_Zauth_User_NotFund)
	}

	// 是否已经绑定
	if has, e := HasUserOrgBindByAOID(ctx, uid, orgID); e != nil {
		return nil, e
	} else if has {
		return nil, zservice.NewError("user already join org:", uid, orgID).SetCode(zglobal.Code_Zauth_UserAlreadyJoin_Org)
	}

	// 准备写入数据
	z := &UserOrgBindTable{
		OrgID:   orgID,
		UID:     uid,
		Expires: Expires,
	}

	if e := z.Save(ctx); e != nil {
		return nil, e
	}
	return z, nil
}

// 是否有账号和组织绑定
func HasUserOrgBindByAOID(ctx *zservice.Context, uid uint32, orgID uint32) (bool, *zservice.Error) {
	return dbhelper.HasTableValue(ctx, &UserOrgBindTable{}, fmt.Sprintf(RK_AOBind_Info, orgID, uid), fmt.Sprintf("uid = %v and org_id = %v", uid, orgID))
}

// 是否过期
func (z *UserOrgBindTable) IsExpired() bool {
	if z.Expires == 0 {
		return false
	}
	return time.Now().Unix() < int64(z.Expires)
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

	rk_info := fmt.Sprintf(RK_AOBind_Info, z.OrgID, z.UID)
	un, e := Redis.Lock(rk_info)
	if e != nil {
		return e
	}
	defer un()

	if e := Mysql.Save(&z).Error; e != nil {
		return zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
	}

	// 删缓存
	zservice.Go(func() {
		if e := Redis.Del(rk_info).Err(); e != nil {
			ctx.LogError(zglobal.Code_Redis_DelFail, e)
		}
	})
	return nil
}
