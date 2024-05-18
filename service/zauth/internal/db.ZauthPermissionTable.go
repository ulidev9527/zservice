package internal

import (
	"fmt"
	"time"
	"zservice/zglobal"
	"zservice/zservice"

	"gorm.io/gorm"
)

// 权限表
type ZauthPermissionTable struct {
	gorm.Model
	Name         string     // 权限名称
	Expires      *time.Time // 过期时间
	PermissionID uint       `gorm:"unique"` // 权限ID
	ParentID     uint       // 父级权限ID
	Service      string     // 权限服务
	Action       string     // 权限动作
	Path         string     // 权限路径
	State        uint       `gorm:"default:3"` // 状态 0禁用 1公开访问 2授权访问 3继承父级
}

// 新建权限
func CreatePermission(ctx *zservice.Context, param ZauthPermissionTable) (*ZauthPermissionTable, *zservice.Error) {

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

	return (&ZauthPermissionTable{
		Name:         param.Name,
		PermissionID: pid,
		ParentID:     param.ParentID,
		Service:      param.Service,
		Action:       param.Action,
		Path:         param.Path,
		State:        param.State,
		Expires:      param.Expires,
	}).Save(ctx)
}

// 同步权限表缓存
func SyncPermissionTableCache(ctx *zservice.Context) *zservice.Error {
	return SyncTableCache(ctx, &[]ZauthPermissionTable{}, func(v any) string {
		return fmt.Sprintf(RK_PermissionInfo, v.(*ZauthPermissionTable).PermissionID)
	})
}

// 获取一个未使用的权限 ID
func GetNewPermissionID(ctx *zservice.Context) (uint, *zservice.Error) {
	return GetNewTableID(ctx, func() uint {
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
	return HasTableValue(ctx, &ZauthPermissionTable{}, fmt.Sprintf(RK_PermissionInfo, id), fmt.Sprintf("permission_id = %v", id))
}

func (z *ZauthPermissionTable) Save(ctx *zservice.Context) (*ZauthPermissionTable, *zservice.Error) {
	rk_info := fmt.Sprintf(RK_PermissionInfo, z.PermissionID)

	// 上锁
	un, e := Redis.Lock(rk_info)
	if e != nil {
		return nil, e
	}
	defer un()

	if z.ID == 0 { // 创建
		if e := Mysql.Create(z).Error; e != nil {
			return nil, zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
		}
	} else { // 更新
		if e := Mysql.Save(z).Error; e != nil {
			return nil, zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
		}
	}

	// 删除缓存
	if e := Redis.Del(rk_info).Err(); e != nil {
		ctx.LogError(e)
	}

	return z, nil
}
