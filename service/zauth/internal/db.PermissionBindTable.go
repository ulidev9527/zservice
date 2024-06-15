package internal

import (
	"fmt"
	"time"
	"zservice/zservice"
	"zservice/zservice/ex/gormservice"
	"zservice/zservice/zglobal"
)

// 账号权限绑定表
type PermissionBindTable struct {
	gormservice.Model
	TargetType   uint32 // 外部ID类型 0无效 1组织 2账号
	TargetID     uint32 // 外部ID
	PermissionID uint32 // 权限ID
	Expires      uint64 // 过期时间
	State        uint32 // 状态 0禁止访问 1允许访问
}

// 是否有权限绑定
func HasPermissionBind(ctx *zservice.Context, targetType uint32, targetID uint32, permissionID uint32) (bool, *zservice.Error) {
	return dbhelper.HasTableValue(ctx,
		&PermissionBindTable{},
		fmt.Sprintf(RK_PermissionBindInfo, targetType, targetID, permissionID),
		fmt.Sprintf("target_type = %d AND target_id = %d AND permission_id = %d", targetType, targetID, permissionID),
	)
}

// 获取权限绑定
func GetPermissionBind(ctx *zservice.Context, targetType uint32, targetID uint32, permissionID uint32) (*PermissionBindTable, *zservice.Error) {
	tab := &PermissionBindTable{}
	if e := dbhelper.GetTableValue(ctx,
		tab,
		fmt.Sprintf(RK_PermissionBindInfo, targetType, targetID, permissionID),
		fmt.Sprintf("target_type = %d AND target_id = %d AND permission_id = %d", targetType, targetID, permissionID),
	); e != nil {
		return nil, e
	}

	return tab, nil
}

// 是否过期
func (z *PermissionBindTable) IsExpired() bool {
	if z.Expires == 0 {
		return false
	}
	return time.Now().Unix() < int64(z.Expires)
}

// 存储
func (z *PermissionBindTable) Save(ctx *zservice.Context) *zservice.Error {

	rk_info := fmt.Sprintf(RK_PermissionBindInfo, z.TargetType, z.TargetID, z.PermissionID)

	// 上锁
	un, e := Redis.Lock(rk_info)
	if e != nil {
		return zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
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
