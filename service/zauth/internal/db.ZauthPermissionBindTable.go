package internal

import (
	"fmt"
	"time"
	"zservice/zservice"
	"zservice/zservice/zglobal"

	"gorm.io/gorm"
)

// 账号权限绑定表
type ZauthPermissionBindTable struct {
	gorm.Model
	TargetType   uint32 // 外部ID类型 0无效 1组织 2账号
	TargetID     uint32 // 外部ID
	PermissionID uint32 // 权限ID
	Expires      uint32 // 过期时间
	State        uint32 // 状态 0禁止访问 1允许访问
}

// 是否有权限绑定
func HasPermissionBind(ctx *zservice.Context, targetType uint32, targetID uint32, permissionID uint32) (bool, *zservice.Error) {
	return dbhelper.HasTableValue(ctx,
		&ZauthPermissionBindTable{},
		fmt.Sprintf(RK_PermissionBindInfo, targetType, targetID, permissionID),
		fmt.Sprintf("target_type = %d AND target_id = %d AND permission_id = %d", targetType, targetID, permissionID),
	)
}

// 获取权限绑定
func GetPermissionBind(ctx *zservice.Context, targetType uint32, targetID uint32, permissionID uint32) (*ZauthPermissionBindTable, *zservice.Error) {
	tab := &ZauthPermissionBindTable{}
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
func (z *ZauthPermissionBindTable) IsExpired() bool {
	if z.Expires == 0 {
		return false
	}
	return time.Now().Unix() < int64(z.Expires)
}

// 存储
func (z *ZauthPermissionBindTable) Save(ctx *zservice.Context) *zservice.Error {

	if z.TargetType == 0 || z.TargetID == 0 || z.PermissionID == 0 {
		return zservice.NewError("param error").SetCode(zglobal.Code_ParamsErr)
	}

	rk_info := fmt.Sprintf(RK_PermissionBindInfo, z.TargetType, z.TargetID, z.PermissionID)

	// 上锁
	un, e := Redis.Lock(rk_info)
	if e != nil {
		return zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
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
