package internal

import (
	"encoding/json"
	"errors"
	"fmt"
	"zservice/zservice"
	"zservice/zservice/ex/gormservice"
	"zservice/zservice/ex/redisservice"
	"zservice/zservice/zglobal"

	"gorm.io/gorm"
)

// 组织表
type OrgTable struct {
	gormservice.Model
	OrgID    uint32 `gorm:"unique"` // 组织ID
	Name     string // 组名
	RootID   uint32 // 根组织ID
	ParentID uint32 // 父级组ID
	State    uint32 `gorm:"default:1"` // 状态 0 禁用 1 开启
}

// 同步组织表缓存
func SyncOrgTableCache(ctx *zservice.Context) *zservice.Error {
	return dbhelper.SyncTableCache(ctx, []OrgTable{}, func(v any) string {
		return fmt.Sprintf(RK_OrgInfo, v.(OrgTable).OrgID)
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
	return dbhelper.HasTableValue(ctx, &OrgTable{}, fmt.Sprintf(RK_OrgInfo, id), fmt.Sprintf("org_id = %v", id))
}

// 根据ID获取一个组织
func GetOrgByID(ctx *zservice.Context, id uint32) (*OrgTable, *zservice.Error) {
	rk_info := fmt.Sprintf(RK_OrgInfo, id)
	tab := &OrgTable{}

	if s, e := Redis.Get(rk_info).Result(); e != nil {
		if !redisservice.IsNilErr(e) {
			return nil, zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
		}
	} else if e := json.Unmarshal([]byte(s), tab); e != nil { // 无法转换为json
		ctx.LogError(e)
	} else if tab.OrgID > 0 { // 大于 0 表示有这个组织
		return tab, nil
	}

	// 未找到 查表
	if e := Mysql.Model(&OrgTable{}).Where("org_id = ?", id).First(tab).Error; e != nil {
		if gormservice.IsNotFound(e) {
			return nil, zservice.NewError(e).SetCode(zglobal.Code_NotFound)
		}
		return nil, zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
	}

	// 设置缓存
	zservice.Go(func() {
		if e := Redis.Set(rk_info, zservice.JsonMustMarshalString(tab)).Err(); e != nil {
			ctx.LogError(e)
		}
	})
	return tab, nil
}

// 获取指定名称的根组织
func GetRootOrgByName(ctx *zservice.Context, name string) (*OrgTable, *zservice.Error) {
	rk_rootName := fmt.Sprintf(RK_OrgRootName, zservice.MD5String(name))
	tab := &OrgTable{}

	if s, e := Redis.Get(rk_rootName).Result(); e != nil {
		if !redisservice.IsNilErr(e) {
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
	if e := Mysql.Model(&OrgTable{}).Where("name = ? AND root_id = 0", name).First(tab).Error; e != nil {
		if !errors.Is(e, gorm.ErrRecordNotFound) {
			return nil, zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
		}
	}

	if tab.OrgID == 0 {
		return nil, nil
	}
	// 更新缓存
	if e := Redis.Set(rk_rootName, zservice.Uint32ToString(tab.OrgID)).Err(); e != nil {
		ctx.LogError(e)
	}

	return tab, nil
}

// 组织存储
func (z *OrgTable) Save(ctx *zservice.Context) *zservice.Error {

	rk_info := fmt.Sprintf(RK_OrgInfo, z.OrgID)

	// 上锁
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
			zservice.LogError(zglobal.Code_Redis_DelFail, e)
		}
	})

	return nil
}
