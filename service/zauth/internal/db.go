package internal

import (
	"zservice/zglobal"
	"zservice/zservice"
	"zservice/zservice/ex/gormservice"
	"zservice/zservice/ex/redisservice"

	"gorm.io/gorm"
)

var (
	MysqlService *gormservice.GormMysqlService
	Mysql        *gorm.DB
	RedisService *redisservice.RedisService
	Redis        *redisservice.GoRedisEX
)

func InitMysql() {

	Mysql.AutoMigrate(&ZauthAccountTable{})
	Mysql.AutoMigrate(&ZauthOrgTable{})
	Mysql.AutoMigrate(&ZauthPermissionTable{})
	Mysql.AutoMigrate(&ZauthAccountOrgBindTable{})
	Mysql.AutoMigrate(&ZauthPermissionBindTable{})
}
func InitRedis() {

	// if e := Mysql.Raw(`
	// WITH RECURSIVE cte(id) AS (
	// 	SELECT g_id FROM account_group_bind_tables WHERE uid=?
	// 	UNION ALL SELECT
	// 	agt.g_id FROM cte JOIN account_group_tables agt ON cte.id = agt.id
	// ) SELECT DISTINCT id FROM cte WHERE id > 0;
	// `, 1001).Find(&[]struct{}{}).Error; e != nil {
	// 	zservice.LogError(e)
	// }
}

// 同步表缓存
func SyncTableCache(ctx *zservice.Context, tabArr any, getRK func(v any) string) *zservice.Error {

	startCount := 0
	limitCount := 200 // 每次同步 200 条
	allCount := 0     // 所有查询到的数据
	errorCount := 0   // 错误数量
	for {
		// 查数据库
		if e := Mysql.Limit(limitCount).Order("created_at ASC").Offset(startCount).Find(&tabArr).Error; e != nil {
			return zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
		}

		arr := tabArr.([]any)

		// 更新缓存
		for _, v := range arr {
			allCount++
			rk_info := getRK(v)
			un, e := Redis.Lock(rk_info)
			if e != nil {
				ctx.LogError(e)
				errorCount++
			}
			defer un()

			if e := Redis.HSet(rk_info, &v).Err(); e != nil {
				ctx.LogError(e)
				errorCount++
			} else {
				Redis.Expire(rk_info, zglobal.Time_10Day) // 设置过期时间
			}
		}

		startCount += limitCount // 更新查询起点

		if len(arr) < limitCount { // 同步完成
			break
		}
	}

	if allCount > 0 && errorCount > 0 {
		if errorCount < allCount {
			return zservice.NewErrorf("SyncOrgTableCache has Error, A:%v E:%v", allCount, errorCount).SetCode(zglobal.Code_Zauth_SyncCacheIncomplete)
		} else {
			return zservice.NewError("SyncOrgTableCache Fail").SetCode(zglobal.Code_Zauth_SyncCacheErr)
		}
	} else {
		return nil
	}
}

// 查询表中是否有指定值
func HasTableValue(ctx *zservice.Context, tab any, rk string, sqlWhere string) (bool, *zservice.Error) {

	if has, e := Redis.Exists(rk).Result(); e != nil {
		return false, zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
	} else if has > 0 {
		return true, nil
	}

	// 验证数据库中是否存在
	count := int64(0)
	if e := Mysql.Model(&tab).Where(sqlWhere).Count(&count).Error; e != nil {
		return false, zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
	}

	return count > 0, nil

}

// 获取一个新的ID
func GetNewTableID(
	ctx *zservice.Context,
	genID func() uint,
	verifyFN func(ctx *zservice.Context, id uint) (bool, *zservice.Error),
	handleErr func(e *zservice.Error) *zservice.Error,
) (uint, *zservice.Error) {
	forCount := 0
	orgID := uint(0)
	for {
		if forCount > 10 {
			return 0, handleErr(zservice.NewError("gen id count max fail").SetCode(zglobal.Code_Zauth_GenIDCountMaxErr))
		}
		orgID = genID()

		if has, e := verifyFN(ctx, orgID); e != nil {
			return 0, handleErr(e)
		} else if has {
			forCount++
			continue
		} else {
			break
		}
	}

	return orgID, nil
}
