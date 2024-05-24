package zhelper

import (
	"zservice/zservice"
	"zservice/zservice/ex/redisservice"
	"zservice/zservice/zglobal"

	"gorm.io/gorm"
)

type DBHelper struct {
	Redis *redisservice.GoRedisEX
	Mysql *gorm.DB
}

func NewDBHelper(Redis *redisservice.GoRedisEX, Mysql *gorm.DB) *DBHelper {
	return &DBHelper{Redis: Redis, Mysql: Mysql}
}

// 同步表缓存
func (db *DBHelper) SyncTableCache(ctx *zservice.Context, tabArr any, getRK func(v any) string) *zservice.Error {

	startCount := 0
	limitCount := 200 // 每次同步 200 条
	allCount := 0     // 所有查询到的数据
	errorCount := 0   // 错误数量
	for {
		// 查数据库
		if e := db.Mysql.Limit(limitCount).Order("created_at ASC").Offset(startCount).Find(&tabArr).Error; e != nil {
			return zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
		}

		arr := tabArr.([]any)

		// 更新缓存
		for _, v := range arr {
			allCount++
			rk_info := getRK(v)
			un, e := db.Redis.Lock(rk_info)
			if e != nil {
				ctx.LogError(e)
				errorCount++
			}
			defer un()

			if e := db.Redis.Set(rk_info, string(zservice.JsonMustMarshal(v))).Err(); e != nil {
				ctx.LogError(e)
				errorCount++
			} else {
				db.Redis.Expire(rk_info, zglobal.Time_10Day) // 设置过期时间
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
func (db *DBHelper) HasTableValue(ctx *zservice.Context, tab any, rk string, sqlWhere string) (bool, *zservice.Error) {

	if has, e := db.Redis.Exists(rk).Result(); e != nil {
		return false, zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
	} else if has > 0 {
		return true, nil
	}

	// 验证数据库中是否存在
	count := int64(0)
	if e := db.Mysql.Model(&tab).Where(sqlWhere).Count(&count).Error; e != nil {
		return false, zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
	}

	return count > 0, nil

}

// 获取一个新的ID
func (db *DBHelper) GetNewTableID(
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

// 获取指定值
func (db *DBHelper) GetTableValue(ctx *zservice.Context, tab any, rk string, sqlWhere string) *zservice.Error {
	// 读缓存
	if has, e := db.Redis.Exists(rk).Result(); e != nil {
		return zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
	} else if has > 0 {
		if e := db.Redis.GetScan(rk, &tab); e != nil {
			return zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
		}

		return nil
	}

	// 验证数据库中是否存在
	count := int64(0)
	wh := db.Mysql.Model(&tab).Where(sqlWhere)

	if e := wh.Count(&count).Error; e != nil {
		return zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
	}
	if count == 0 {
		return zservice.NewError("not found").SetCode(zglobal.Code_DB_NotFound)
	}
	if e := db.Mysql.Model(&tab).Where(sqlWhere).First(&tab).Error; e != nil {
		return zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
	}

	// 更新缓存
	if e := db.Redis.Set(rk, string(zservice.JsonMustMarshal(tab))).Err(); e != nil {
		return zservice.NewError(e).SetCode(zglobal.Code_ErrorBreakoff)
	}

	return nil
}
