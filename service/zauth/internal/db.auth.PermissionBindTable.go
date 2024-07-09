package internal

import (
	"fmt"
	"time"
	"zservice/zservice"
	"zservice/zserviceex/dbservice"

	"gorm.io/gorm"
)

// 账号权限绑定表
type PermissionBindTable struct {
	gorm.Model
	TargetType   uint32        // 外部ID类型 0无效 1组织 2账号
	TargetID     uint32        // 外部ID
	PermissionID uint32        // 权限ID
	Expires      zservice.Time // 过期时间
	State        uint32        // 状态 0禁止访问 1允许访问
}

// 是否有权限绑定
func HasPermissionBind(ctx *zservice.Context, targetType uint32, targetID uint32, permissionID uint32) (bool, *zservice.Error) {
	return DBService.HasTableValue(ctx, dbservice.HasTableValueOption{
		Tab:      PermissionBindTable{},
		RK:       fmt.Sprintf(RK_PermissionBindInfo, targetType, targetID, permissionID),
		SQLConds: []any{"target_type = ? AND target_id = ? AND permission_id = ?", targetType, targetID, permissionID},
	})
}

// 获取权限绑定
func GetPermissionBind(ctx *zservice.Context, targetType uint32, targetID uint32, permissionID uint32) (*PermissionBindTable, *zservice.Error) {
	tab := &PermissionBindTable{}
	if e := DBService.GetTableValue(ctx, dbservice.GetTableValueOption{
		Tab:      tab,
		RK:       fmt.Sprintf(RK_PermissionBindInfo, targetType, targetID, permissionID),
		SQLConds: []any{"target_type = ? AND target_id = ? AND permission_id = ?", targetType, targetID, permissionID},
	}); e != nil {
		return nil, e.AddCaller()
	}

	return tab, nil
}

// 是否过期
func (z *PermissionBindTable) IsExpired() bool {
	if z.Expires.IsZero() {
		return false
	}
	return z.Expires.After(time.Now())
}

// 存储
func (z *PermissionBindTable) Save(ctx *zservice.Context) *zservice.Error {

	rk_info := fmt.Sprintf(RK_PermissionBindInfo, z.TargetType, z.TargetID, z.PermissionID)

	// 上锁
	un, e := Redis.Lock(rk_info)
	if e != nil {
		return e.AddCaller()
	}
	defer un()

	if e := Gorm.Save(&z).Error; e != nil {
		return zservice.NewError(e)
	}

	// 删缓存
	zservice.Go(func() {
		if e := Redis.Del(rk_info).Err(); e != nil {
			zservice.LogError(zservice.Code_Redis_DelFail, e)
		}
	})

	return nil
}