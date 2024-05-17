package internal

import (
	"fmt"
	"time"
	"zservice/zglobal"
	"zservice/zservice"

	"gorm.io/gorm"
)

// 账号组织绑定
type ZauthAccountOrgBindTable struct {
	gorm.Model

	OrgID     uint       // 组ID
	AccountID uint       // 用户ID
	Expires   *time.Time // 过期时间
	State     uint       `gorm:"default:1"` // 状态 0禁用 1开启
}

// 加入组织
func AccountJoinOrg(ctx *zservice.Context, accountID uint, orgID uint, Expires *time.Time) (*ZauthAccountOrgBindTable, *zservice.Error) {
	// 验证参数是否正确
	if has, e := HasOrgByID(ctx, orgID); e != nil {
		return nil, e
	} else if !has {
		return nil, zservice.NewError("org not found:", orgID).SetCode(zglobal.Code_Zauth_OrgNotFund)
	}

	// 是否已经绑定
	rk_info := fmt.Sprintf(RK_AccountBindOrgInfo, orgID, accountID)
	// 上锁
	un, e := Redis.Lock(rk_info)
	if e != nil {
		return nil, e
	}
	defer un()

	// 查数数据是否已经存在
	// 查缓存

	if has, e := HasTableValue(ctx, &ZauthAccountOrgBindTable{}, rk_info, fmt.Sprintf("account_id = %v and org_id = %v", accountID, orgID)); e != nil {
		return nil, e
	} else if has {
		return nil, zservice.NewError("account already join org:", accountID, orgID).SetCode(zglobal.Code_Zauth_AccountAlreadyJoin_Org)
	}

	// 准备写入数据
	z := &ZauthAccountOrgBindTable{
		OrgID:     orgID,
		AccountID: accountID,
		Expires:   Expires,
	}

	if e := Mysql.Create(&z).Error; e != nil {
		return nil, zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
	}

	// 存 redis
	if e := Redis.HMSet(rk_info, &z).Err(); e != nil {
		ctx.LogError(e)
	}
	return z, nil
}

// 是否过期
func (z *ZauthAccountOrgBindTable) IsExpired() bool {
	if z.Expires == nil {
		return false
	}
	return z.Expires.Before(time.Now())
}

// 存储
func (z *ZauthAccountOrgBindTable) Save(ctx *zservice.Context) (*ZauthAccountOrgBindTable, *zservice.Error) {

	if z.OrgID == 0 || z.AccountID == 0 {
		return nil, zservice.NewError("param error").SetCode(zglobal.Code_ParamsErr)
	}

	rk_info := fmt.Sprintf(RK_AccountBindOrgInfo, z.OrgID, z.AccountID)
	un, e := Redis.Lock(rk_info)
	if e != nil {
		return nil, e
	}
	defer un()

	if e := Mysql.Save(&z).Error; e != nil {
		return nil, zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
	}

	// 存 redis
	if e := Redis.HMSet(rk_info, &z).Err(); e != nil {
		ctx.LogError(e)
	}
	return z, nil
}
