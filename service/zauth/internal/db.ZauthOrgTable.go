package internal

import (
	"fmt"
	"zservice/zservice"
	"zservice/zservice/zglobal"

	"gorm.io/gorm"
)

// 组织表
type ZauthOrgTable struct {
	gorm.Model
	Name        string // 组名
	OrgID       uint   `gorm:"unique"` // 组织ID
	RootOrgID   uint   // 根组织ID
	ParentOrgID uint   // 父级组ID
	State       uint   `gorm:"default:1"` // 状态 0 禁用 1 开启
}

// 新建一个根组织
func CreateRootOrg(ctx *zservice.Context, name string) (*ZauthOrgTable, *zservice.Error) {

	// 创建锁
	un, e := Redis.Lock(RK_OrgCreateLock)
	if e != nil {
		return nil, e
	}
	defer un()

	// 获取一个未使用的组织 ID
	orgID, e := GetNewOrgID(ctx)
	if e != nil {
		return nil, e
	}

	z := &ZauthOrgTable{
		Name:      name,
		OrgID:     orgID,
		RootOrgID: orgID,
	}
	if e := z.Save(ctx); e != nil {
		return nil, e
	}
	return z, nil
}

// 新建一个组织
func CreateOrg(ctx *zservice.Context, name string, rootOrgID uint, parentOrgID uint) (*ZauthOrgTable, *zservice.Error) {

	// 验证组织是否存在
	// 根组织验证
	if has, e := HasOrgByID(ctx, rootOrgID); e != nil {
		return nil, e
	} else if !has {
		return nil, zservice.NewError("org not found:", rootOrgID).SetCode(zglobal.Code_Zauth_OrgCreateRootIDErr)
	}
	// 父级组织验证
	if rootOrgID != parentOrgID {
		if has, e := HasOrgByID(ctx, parentOrgID); e != nil {
			return nil, e
		} else if !has {
			return nil, zservice.NewError("org not found:", parentOrgID).SetCode(zglobal.Code_Zauth_OrgCreateParentIDErr)
		}
	}

	// 创建锁
	un, e := Redis.Lock(RK_OrgCreateLock)
	if e != nil {
		return nil, e
	}
	defer un()

	// 获取一个未使用的组织 ID
	orgID, e := GetNewOrgID(ctx)
	if e != nil {
		return nil, e
	}

	z := &ZauthOrgTable{
		Name:        name,
		OrgID:       orgID,
		RootOrgID:   rootOrgID,
		ParentOrgID: parentOrgID,
	}

	if e := z.Save(ctx); e != nil {
		return nil, e
	}
	return z, nil
}

// 同步组织表缓存
func SyncOrgTableCache(ctx *zservice.Context) *zservice.Error {
	return dbhelper.SyncTableCache(ctx, []ZauthOrgTable{}, func(v any) string {
		return fmt.Sprintf(RK_OrgInfo, v.(ZauthOrgTable).OrgID)
	})
}

// 获取一个全新的组织ID
func GetNewOrgID(ctx *zservice.Context) (uint, *zservice.Error) {
	return dbhelper.GetNewTableID(ctx, func() uint {
		return uint(zservice.RandomIntRange(100000, 99999999)) // 1 是根权限，并且预留 10w 以下的id
	}, HasOrgByID, func(e *zservice.Error) *zservice.Error {
		if e.GetCode() == zglobal.Code_Zauth_GenIDCountMaxErr {
			return e.SetCode(zglobal.Code_Zauth_OrgGenIDCountMaxErr)
		}
		return e
	})
}

// 是否存在这个组织
func HasOrgByID(ctx *zservice.Context, orgID uint) (bool, *zservice.Error) {
	return dbhelper.HasTableValue(ctx, &ZauthOrgTable{}, fmt.Sprintf(RK_OrgInfo, orgID), fmt.Sprintf("org_id = %v", orgID))
}

// 组织存储
func (z *ZauthOrgTable) Save(ctx *zservice.Context) *zservice.Error {

	if z.OrgID == 0 || z.RootOrgID == 0 {
		return zservice.NewError("param error").SetCode(zglobal.Code_ParamsErr)
	}

	rk_info := fmt.Sprintf(RK_OrgInfo, z.OrgID)

	// 上锁
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
