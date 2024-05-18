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

	if has, e := HasAccountByID(ctx, accountID); e != nil {
		return nil, e
	} else if !has {
		return nil, zservice.NewError("account not found:", accountID).SetCode(zglobal.Code_Zauth_Account_NotFund)
	}

	// 是否已经绑定
	if has, e := HasAccountOrgBindByAOID(ctx, accountID, orgID); e != nil {
		return nil, e
	} else if has {
		return nil, zservice.NewError("account already join org:", accountID, orgID).SetCode(zglobal.Code_Zauth_AccountAlreadyJoin_Org)
	}

	rk_lock := fmt.Sprintf(RK_AOBind_CreateLock, orgID, accountID)
	// 上锁
	un, e := Redis.Lock(rk_lock)
	if e != nil {
		return nil, e
	}
	defer un()

	// 准备写入数据
	z := &ZauthAccountOrgBindTable{
		OrgID:     orgID,
		AccountID: accountID,
		Expires:   Expires,
	}

	return z.Save(ctx)
}

// 是否有账号和组织绑定
func HasAccountOrgBindByAOID(ctx *zservice.Context, accountID uint, orgID uint) (bool, *zservice.Error) {
	return HasTableValue(ctx, &ZauthAccountOrgBindTable{}, fmt.Sprintf(RK_AOBind_Info, orgID, accountID), fmt.Sprintf("account_id = %v and org_id = %v", accountID, orgID))
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

	rk_info := fmt.Sprintf(RK_AOBind_Info, z.OrgID, z.AccountID)
	un, e := Redis.Lock(rk_info)
	if e != nil {
		return nil, e
	}
	defer un()

	if z.ID == 0 { // 创建
		if e := Mysql.Create(&z).Error; e != nil {
			return nil, zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
		}
	} else { // 更新
		if e := Mysql.Save(&z).Error; e != nil {
			return nil, zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
		}
	}

	// 删缓存
	if e := Redis.Del(rk_info).Err(); e != nil {
		return z, zservice.NewError(e).SetCode(zglobal.Code_Redis_DelFail)
	}
	return z, nil
}