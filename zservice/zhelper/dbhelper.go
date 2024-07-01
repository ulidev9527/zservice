package zhelper

import (
	"errors"
	"time"
	"zservice/zservice"
	"zservice/zservice/ex/redisservice"
	"zservice/zservice/zglobal"

	"gorm.io/gorm"
)

type DBHelper struct {
	redis *redisservice.GoRedisEX
	mysql *gorm.DB
}

func NewDBHelper(Redis *redisservice.GoRedisEX, Mysql *gorm.DB) *DBHelper {
	return &DBHelper{redis: Redis, mysql: Mysql}
}

// 同步表缓存
func (db *DBHelper) SyncTableCache(ctx *zservice.Context, tabArr any, getRK func(v any) string) *zservice.Error {

	startCount := 0
	limitCount := 200 // 每次同步 200 条
	allCount := 0     // 所有查询到的数据
	errorCount := 0   // 错误数量
	for {
		// 查数据库
		if e := db.mysql.Limit(limitCount).Order("created_at ASC").Offset(startCount).Find(&tabArr).Error; e != nil {
			return zservice.NewError(e)
		}

		arr := tabArr.([]any)

		// 更新缓存
		for _, v := range arr {
			allCount++
			rk_info := getRK(v)
			un, e := db.redis.Lock(rk_info)
			if e != nil {
				ctx.LogError(e)
				errorCount++
			}
			defer un()

			if e := db.redis.Set(rk_info, string(zservice.JsonMustMarshal(v))).Err(); e != nil {
				ctx.LogError(e)
				errorCount++
			} else {
				db.redis.Expire(rk_info, zglobal.Time_10Day) // 设置过期时间
			}
		}

		startCount += limitCount // 更新查询起点

		if len(arr) < limitCount { // 同步完成
			break
		}
	}

	if allCount > 0 && errorCount > 0 {
		if errorCount < allCount {
			return zservice.NewErrorf("SyncOrgTableCache has Error, A:%v E:%v", allCount, errorCount).SetCode(zglobal.Code_SyncCacheIncomplete)
		} else {
			return zservice.NewError("SyncOrgTableCache Fail").SetCode(zglobal.Code_SyncCacheErr)
		}
	} else {
		return nil
	}
}

// 查询表中是否有指定值
func (db *DBHelper) HasTableValue(ctx *zservice.Context, tab any, rk string, sqlWhere string) (bool, *zservice.Error) {

	if has, e := db.redis.Exists(rk).Result(); e != nil {
		return false, zservice.NewError(e)
	} else if has > 0 {
		return true, nil
	}

	// 验证数据库中是否存在
	count := int64(0)
	if e := db.mysql.Model(&tab).Where(sqlWhere).Count(&count).Error; e != nil {
		return false, zservice.NewError(e)
	}

	return count > 0, nil

}

// 获取一个新的ID
func (db *DBHelper) GetNewTableID(
	ctx *zservice.Context,
	genID func() uint32,
	verifyFN func(ctx *zservice.Context, id uint32) (bool, *zservice.Error),
) (uint32, *zservice.Error) {
	forCount := 0
	orgID := uint32(0)
	for {
		if forCount > 10 {
			return 0, zservice.NewError("gen id count max fail").SetCode(zglobal.Code_GenIDCountMaxErr)
		}
		orgID = genID()

		if has, e := verifyFN(ctx, orgID); e != nil {
			return 0, e
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
// 注意，如果没找到数据回返回：zglobal.Code_NotFound
func (db *DBHelper) GetTableValue(ctx *zservice.Context, tab any, rk string, sqlWhere string, expires ...time.Duration) *zservice.Error {
	// 读缓存
	if e := db.redis.GetScan(rk, &tab); e != nil {
		if e.GetCode() != zglobal.Code_NotFound {
			return e.AddCaller()
		}
	}

	// 查库
	if e := db.mysql.First(tab, sqlWhere).Error; e != nil {
		if errors.Is(e, gorm.ErrRecordNotFound) {
			return zservice.NewError(e).SetCode(zglobal.Code_NotFound)
		} else {
			return zservice.NewError(e)
		}
	}

	// 更新缓存
	zservice.Go(func() {
		if e := func() error {
			if len(expires) > 0 {
				return db.redis.SetEX(rk, string(zservice.JsonMustMarshal(tab)), expires[0]).Err()
			} else {
				return db.redis.SetEX(rk, string(zservice.JsonMustMarshal(tab)), zglobal.Time_3Day).Err()
			}
		}(); e != nil {
			ctx.LogError(e)
		}
	})

	return nil
}
