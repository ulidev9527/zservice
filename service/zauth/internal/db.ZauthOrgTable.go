package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"zservice/zservice"
	"zservice/zservice/ex/gormservice"
	"zservice/zservice/zglobal"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// 组织表
type ZauthOrgTable struct {
	gormservice.TimeModel
	ID       uint32 `gorm:"primaryKey"` // 组织ID
	Name     string // 组名
	RootID   uint32 // 根组织ID
	ParentID uint32 // 父级组ID
	State    uint32 `gorm:"default:1"` // 状态 0 禁用 1 开启
}

// 同步组织表缓存
func SyncOrgTableCache(ctx *zservice.Context) *zservice.Error {
	return dbhelper.SyncTableCache(ctx, []ZauthOrgTable{}, func(v any) string {
		return fmt.Sprintf(RK_OrgInfo, v.(ZauthOrgTable).ID)
	})
}

// 获取一个全新的组织ID
func GetNewOrgID(ctx *zservice.Context) (uint32, *zservice.Error) {
	return dbhelper.GetNewTableID(ctx, func() uint32 {
		return zservice.RandomUInt32Range(100000, 99999999) // 1 是根权限，并且预留 10w 以下的id
	}, HasOrgByID)
}

// 是否存在这个组织
func HasOrgByID(ctx *zservice.Context, id uint32) (bool, *zservice.Error) {
	return dbhelper.HasTableValue(ctx, &ZauthOrgTable{}, fmt.Sprintf(RK_OrgInfo, id), fmt.Sprintf("id = %v", id))
}

// 根据ID获取一个组织
func GetOrgByID(ctx *zservice.Context, id uint32) (*ZauthOrgTable, *zservice.Error) {
	rk_info := fmt.Sprintf(RK_OrgInfo, id)
	tab := &ZauthOrgTable{}

	if s, e := Redis.Get(rk_info).Result(); e != nil {
		if e != redis.Nil {
			return nil, zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
		}
	} else if e := json.Unmarshal([]byte(s), tab); e != nil {
		return nil, zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
	} else if tab.ID > 0 {
		return tab, nil
	}

	// 未找到 查表
	if e := Mysql.Model(&ZauthOrgTable{}).Where("id = ?", id).First(tab).Error; e != nil {
		return nil, zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
	}
	if tab.ID > 0 {
		if e := Redis.Set(rk_info, zservice.JsonMustMarshalString(tab)).Err(); e != nil {
			ctx.LogError(e)
		}
		return tab, nil
	}
	return nil, nil
}

// 是否存在指定名称的根组织
func GetRootOrgByName(ctx *zservice.Context, name string) (*ZauthOrgTable, *zservice.Error) {
	rk_rootName := fmt.Sprintf(RK_OrgRootName, zservice.MD5String(name))
	tab := &ZauthOrgTable{}

	if s, e := Redis.Get(rk_rootName).Result(); e != nil {
		if e != redis.Nil {
			return nil, zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
		}
	} else {
		if tab, e := GetOrgByID(ctx, zservice.StringToUint32(s)); e != nil {
			return nil, e
		} else if tab != nil {
			return tab, nil
		}
	}

	// 未找到 查表
	if e := Mysql.Model(&ZauthOrgTable{}).Where("name = ? AND root_id = 0", name).First(tab).Error; e != nil {
		if !errors.Is(e, gorm.ErrRecordNotFound) {
			return nil, zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
		}
	}

	if tab.ID == 0 {
		return nil, nil
	}
	// 更新缓存
	if e := Redis.Set(rk_rootName, zservice.Uint32ToString(tab.ID)).Err(); e != nil {
		ctx.LogError(e)
	}

	return tab, nil
}

// 组织存储
func (z *ZauthOrgTable) Save(ctx *zservice.Context) *zservice.Error {

	if z.ID == 0 {
		return zservice.NewError("param error").SetCode(zglobal.Code_ParamsErr)
	}

	rk_info := fmt.Sprintf(RK_OrgInfo, z.ID)

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
