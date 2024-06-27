package internal

import (
	"fmt"
	"zservice/zservice"
	"zservice/zservice/ex/gormservice"
	"zservice/zservice/ex/redisservice"
	"zservice/zservice/zglobal"
)

// 权限表
type PermissionTable struct {
	gormservice.Model
	PermissionID uint32 `gorm:"unique"` // 权限ID
	Name         string // 权限名称
	Service      string // 服务名称
	Action       string // 权限动作
	Path         string // 权限路径
	State        uint32 `gorm:"default:3"` // 状态 0禁用 1公开访问 2授权访问 3继承父级(默认)
}

// 同步权限表缓存
func SyncPermissionTableCache(ctx *zservice.Context) *zservice.Error {
	return dbhelper.SyncTableCache(ctx, &[]PermissionTable{}, func(v any) string {
		return fmt.Sprintf(RK_PermissionInfo, v.(*PermissionTable).PermissionID)
	})
}

// 获取一个未使用的权限 ID
func GetNewPermissionID(ctx *zservice.Context) (uint32, *zservice.Error) {
	return dbhelper.GetNewTableID(ctx, func() uint32 {
		return zservice.RandomUInt32Range(1, 9999999)
	}, HasPermissionByID)
}

// 权限是否存在
func HasPermissionByID(ctx *zservice.Context, id uint32) (bool, *zservice.Error) {
	return dbhelper.HasTableValue(ctx, &PermissionTable{}, fmt.Sprintf(RK_PermissionInfo, id), fmt.Sprintf("permission_id = %v", id))
}

// 根据ID获取一个权限
func GetPermissionByID(ctx *zservice.Context, id uint) (*PermissionTable, *zservice.Error) {
	tab := &PermissionTable{}
	if e := dbhelper.GetTableValue(ctx, tab, fmt.Sprintf(RK_PermissionInfo, id), fmt.Sprintf("permission_id = %v", id)); e != nil {
		return nil, e
	}
	return tab, nil
}

// 获取指定权限
func GetPermissionBySAP(ctx *zservice.Context, service, action, path string) (*PermissionTable, *zservice.Error) {
	rk_sap := fmt.Sprintf(RK_PermissionSAP, service, action, path)
	if s, e := Redis.Get(rk_sap).Result(); e != nil {
		if !redisservice.IsNilErr(e) {
			return nil, zservice.NewError(e)
		}

	} else {
		if tab, e := GetPermissionByID(ctx, zservice.StringToUint(s)); e != nil {
			if e.GetCode() != zglobal.Code_NotFound {
				return nil, e.AddCaller()

			}
		} else {
			return tab, nil

		}

	}

	// 未找到 查表
	tab := &PermissionTable{}
	if e := Mysql.Model(&PermissionTable{}).Where("service = ? AND action = ? AND path = ?", service, action, path).First(tab).Error; e != nil {
		if gormservice.IsNotFound(e) {
			return nil, zservice.NewError(e).SetCode(zglobal.Code_NotFound)
		}
		return nil, zservice.NewError(e)
	}

	// 缓存
	zservice.Go(func() {
		if e := Redis.Set(fmt.Sprintf(RK_PermissionInfo, tab.PermissionID), zservice.JsonMustMarshalString(tab)).Err(); e != nil {
			ctx.LogError(e)
		}
		if e := Redis.Set(rk_sap, zservice.Uint32ToString(tab.PermissionID)).Err(); e != nil {
			ctx.LogError(e)
		}
	})

	return tab, nil
}

// 存储
func (z *PermissionTable) Save(ctx *zservice.Context) *zservice.Error {
	rk_info := fmt.Sprintf(RK_PermissionInfo, z.PermissionID)

	// 上锁
	un, e := Redis.Lock(rk_info)
	if e != nil {
		return e
	}
	defer un()

	if e := Mysql.Save(z).Error; e != nil {
		return zservice.NewError(e)
	}

	// 删除缓存
	zservice.Go(func() {
		if e := Redis.Del(rk_info).Err(); e != nil {
			zservice.LogError(zglobal.Code_Redis_DelFail, e)
		}
	})

	return nil
}
