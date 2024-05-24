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
	TargetType   uint       // 外部ID类型 0无效 1组织 2账号
	TargetID     uint       // 外部ID
	PermissionID uint       // 权限ID
	Expires      *time.Time // 过期时间
	Allow        bool       // 状态 0禁止访问 1允许访问
}

// 权限绑定
func PermissionBind(ctx *zservice.Context, targetType uint, targetID uint, permissionID uint, Expires *time.Time, Allow bool) (*ZauthPermissionBindTable, *zservice.Error) {

	// 验证参数是否正确
	switch targetType {
	case 1:
		// 组织验证
		if has, e := HasOrgByID(ctx, targetID); e != nil {
			return nil, e
		} else if !has {
			return nil, zservice.NewError("otherID invalid:", targetID).SetCode(zglobal.Code_Zauth_PermissionBind_TargetIDErr)
		}
	case 2:
		// 账号验证
		if has, e := HasAccountByID(ctx, targetID); e != nil {
			return nil, e
		} else if !has {
			return nil, zservice.NewError("otherID invalid:", targetID).SetCode(zglobal.Code_Zauth_PermissionBind_TargetIDErr)
		}
	default:
		return nil, zservice.NewError("otherIDType invalid:", targetType).SetCode(zglobal.Code_Zauth_PermissionBind_TargetTypeErr)
	}

	// 权限验证
	if has, e := HasPermissionByID(ctx, permissionID); e != nil {
		return nil, e
	} else if !has {
		return nil, zservice.NewError("permissionID invalid:", permissionID).SetCode(zglobal.Code_Zauth_PermissionBind_PermissionIDErr)
	}

	// 是否有绑定
	if has, e := HasPermissionBind(ctx, targetType, targetID, permissionID); e != nil {
		return nil, e
	} else if has {
		return nil, zservice.NewError("already bind").SetCode(zglobal.Code_Zauth_PermissionBind_Already_Bind)
	}

	// 绑定
	bind := &ZauthPermissionBindTable{
		TargetType:   targetType,
		TargetID:     targetID,
		PermissionID: permissionID,
		Expires:      Expires,
		Allow:        Allow,
	}

	if e := bind.Save(ctx); e != nil {
		return nil, e
	}

	return bind, nil

}

// 是否有权限绑定
func HasPermissionBind(ctx *zservice.Context, targetType uint, targetID uint, permissionID uint) (bool, *zservice.Error) {
	return dbhelper.HasTableValue(ctx,
		&ZauthPermissionBindTable{},
		fmt.Sprintf(RK_PermissionBindInfo, targetType, targetID, permissionID),
		fmt.Sprintf("target_type = %d AND target_id = %d AND permission_id = %d", targetType, targetID, permissionID),
	)
}

// 获取权限绑定
func GetPermissionBind(ctx *zservice.Context, targetType uint, targetID uint, permissionID uint) (*ZauthPermissionBindTable, *zservice.Error) {
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
	if z.Expires == nil {
		return false
	}
	return z.Expires.Before(time.Now())
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
