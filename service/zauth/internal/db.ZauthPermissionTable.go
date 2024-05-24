package internal

import (
	"errors"
	"fmt"
	"zservice/zservice"
	"zservice/zservice/zglobal"

	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// 权限表
type ZauthPermissionTable struct {
	gorm.Model
	Name         string // 权限名称
	PermissionID uint   `gorm:"unique"` // 权限ID
	Service      string // 权限服务
	Action       string // 权限动作
	Path         string // 权限路径
	State        uint   `gorm:"default:3"` // 状态 0禁用 1公开访问 2授权访问 3继承父级
}

// 创建权限的配置
type CreatePermissionConfig struct {
	Name    string // 权限名称
	Service string // 权限服务
	Action  string // 权限动作
	Path    string // 权限路径
	State   uint   `gorm:"default:"` // 状态 0禁用 1公开访问 2授权访问 3继承父级
}

// 新建权限
func CreatePermission(ctx *zservice.Context, param CreatePermissionConfig) (*ZauthPermissionTable, *zservice.Error) {

	// 锁
	un, e := Redis.Lock(RK_PermissionCreateLock)
	if e != nil {
		return nil, e
	}
	defer un()

	// 获取一个未使用的权限 ID
	pid, e := GetNewPermissionID(ctx)
	if e != nil {
		return nil, e
	}

	z := &ZauthPermissionTable{
		Name:         param.Name,
		PermissionID: pid,
		Service:      param.Service,
		Action:       param.Action,
		Path:         param.Path,
		State:        param.State,
	}

	if e := z.Save(ctx); e != nil {
		return nil, e
	}

	return z, nil
}

// 同步权限表缓存
func SyncPermissionTableCache(ctx *zservice.Context) *zservice.Error {
	return dbhelper.SyncTableCache(ctx, &[]ZauthPermissionTable{}, func(v any) string {
		return fmt.Sprintf(RK_PermissionInfo, v.(*ZauthPermissionTable).PermissionID)
	})
}

// 获取一个未使用的权限 ID
func GetNewPermissionID(ctx *zservice.Context) (uint, *zservice.Error) {
	return dbhelper.GetNewTableID(ctx, func() uint {
		return uint(zservice.RandomIntRange(1, 9999999))
	}, HasPermissionByID, func(e *zservice.Error) *zservice.Error {
		if e.GetCode() == zglobal.Code_Zauth_GenIDCountMaxErr {
			return e.SetCode(zglobal.Code_Zauth_PermissionGenIDCountMaxErr)
		}
		return e
	})
}

// 权限是否存在
func HasPermissionByID(ctx *zservice.Context, id uint) (bool, *zservice.Error) {
	return dbhelper.HasTableValue(ctx, &ZauthPermissionTable{}, fmt.Sprintf(RK_PermissionInfo, id), fmt.Sprintf("permission_id = %v", id))
}

// 根据ID获取一个权限
func GetPermissionByID(ctx *zservice.Context, id uint) (*ZauthPermissionTable, *zservice.Error) {
	tab := &ZauthPermissionTable{}
	if e := dbhelper.GetTableValue(ctx, tab, fmt.Sprintf(RK_PermissionInfo, id), fmt.Sprintf("permission_id = %v", id)); e != nil {
		return nil, e
	}
	return tab, nil
}

// 获取指定权限
func GetPermissionBySAP(ctx *zservice.Context, service, action, path string) (*ZauthPermissionTable, *zservice.Error) {
	rk_sap := fmt.Sprintf(RK_PermissionSAP, service, action, path)
	if s, e := Redis.Get(rk_sap).Result(); e != nil {
		if e != redis.Nil {
			return nil, zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
		}

	} else {
		if tab, e := GetPermissionByID(ctx, zservice.StringToUint(s)); e != nil {
			if e.GetCode() != zglobal.Code_DB_NotFound {
				return nil, zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
			}
		} else {
			return tab, nil
		}
	}

	// 未找到 查表
	tab := &ZauthPermissionTable{}
	if e := Mysql.Model(&ZauthPermissionTable{}).Where("service = ? AND action = ? AND path = ?", service, action, path).First(tab).Error; e != nil {
		if !errors.Is(e, gorm.ErrRecordNotFound) {
			return nil, zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
		}
	}

	if tab.ID == 0 {
		return nil, zservice.NewError("not found").SetCode(zglobal.Code_DB_NotFound)
	}
	// 缓存
	if e := Redis.Set(fmt.Sprintf(RK_PermissionInfo, tab.PermissionID), zservice.JsonMustMarshalString(tab)).Err(); e != nil {
		ctx.LogError(e)
	}
	if e := Redis.Set(rk_sap, zservice.UIntToString(tab.PermissionID)).Err(); e != nil {
		ctx.LogError(e)
	}

	return tab, nil

}

// 存储
func (z *ZauthPermissionTable) Save(ctx *zservice.Context) *zservice.Error {
	rk_info := fmt.Sprintf(RK_PermissionInfo, z.PermissionID)

	// 上锁
	un, e := Redis.Lock(rk_info)
	if e != nil {
		return e
	}
	defer un()

	if z.ID == 0 { // 创建
		if e := Mysql.Create(z).Error; e != nil {
			return zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
		}
	} else { // 更新
		if e := Mysql.Save(z).Error; e != nil {
			return zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
		}
	}

	// 删除缓存
	if e := Redis.Del(rk_info).Err(); e != nil {
		zservice.LogError(zglobal.Code_Redis_DelFail, e)
	}

	return nil
}
