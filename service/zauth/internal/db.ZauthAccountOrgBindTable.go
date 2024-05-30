package internal

import (
	"fmt"
	"time"
	"zservice/zservice"
	"zservice/zservice/zglobal"

	"gorm.io/gorm"
)

// 账号组织绑定
type ZauthAccountOrgBindTable struct {
	gorm.Model

	OrgID     uint32 // 组ID
	AccountID uint32 // 用户ID
	Expires   uint32 // 过期时间
	State     uint32 `gorm:"default:1"` // 状态 0禁用 1开启
}

// 加入组织
func AccountJoinOrg(ctx *zservice.Context, accountID uint32, orgID uint32, Expires uint32) (*ZauthAccountOrgBindTable, *zservice.Error) {
	// 验证参数是否正确
	if has, e := HasOrgByID(ctx, orgID); e != nil {
		return nil, e
	} else if !has {
		return nil, zservice.NewError("org not found:", orgID).SetCode(zglobal.Code_Zauth_Org_NotFund)
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

	// 准备写入数据
	z := &ZauthAccountOrgBindTable{
		OrgID:     orgID,
		AccountID: accountID,
		Expires:   Expires,
	}

	if e := z.Save(ctx); e != nil {
		return nil, e
	}
	return z, nil
}

// 是否有账号和组织绑定
func HasAccountOrgBindByAOID(ctx *zservice.Context, accountID uint32, orgID uint32) (bool, *zservice.Error) {
	return dbhelper.HasTableValue(ctx, &ZauthAccountOrgBindTable{}, fmt.Sprintf(RK_AOBind_Info, orgID, accountID), fmt.Sprintf("account_id = %v and org_id = %v", accountID, orgID))
}

// 是否过期
func (z *ZauthAccountOrgBindTable) IsExpired() bool {
	if z.Expires == 0 {
		return false
	}
	return time.Now().Unix() < int64(z.Expires)
}

// 是否启动
func (z *ZauthAccountOrgBindTable) IsAllow() bool {
	if z.IsExpired() {
		return false
	}
	return z.State == 1
}

// 存储
func (z *ZauthAccountOrgBindTable) Save(ctx *zservice.Context) *zservice.Error {

	if z.OrgID == 0 || z.AccountID == 0 {
		return zservice.NewError("param error").SetCode(zglobal.Code_ParamsErr)
	}

	rk_info := fmt.Sprintf(RK_AOBind_Info, z.OrgID, z.AccountID)
	un, e := Redis.Lock(rk_info)
	if e != nil {
		return e
	}
	defer un()

	if z.ID == 0 { // 创建
		if e := Mysql.Create(&z).Error; e != nil {
			return zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
		}
	} else { // 更新
		if e := Mysql.Save(&z).Error; e != nil {
			return zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
		}
	}

	// 删缓存
	if e := Redis.Del(rk_info).Err(); e != nil {
		zservice.LogError(zglobal.Code_Redis_DelFail, e)
	}
	return nil
}
