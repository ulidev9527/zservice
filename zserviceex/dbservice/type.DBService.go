package dbservice

import (
	"fmt"
	"time"
	"zservice/zservice"
)

// 数据服务，结合：github.com/redis/go-redis/v9 和 gorm.io/gorm 配合使用
type DBService struct {
	*zservice.ZService
	option DBServiceOption // 配置
	Gorm   *GormEX         // 数据库 gorm
	Redis  *GoRedisEX      // 缓存 redis
}

type DBServiceOption struct {
	Name        string           // 服务名称 仅用于日志显示，如果主服务中有多个DBService，建议配置，如果只有一个DBService，可以忽略此配置
	GORMType    string           // 数据库类型 目前仅支持 mysql
	GORMName    string           // 数据库名称
	GORMAddr    string           // 数据库地址 填入地址才会启用 DB 功能
	GORMUser    string           // 数据库用户名
	GORMPass    string           // 数据库密码
	RedisAddr   string           // redis 地址 填入地址才会启用 Redis 功能
	RedisPass   string           // redis 密码
	RedisPrefix string           // redis 前缀 默认使用 zservice.Init 中的 serviceName
	Debug       bool             // 是否开启 debug
	OnStart     func(*DBService) // 启动的回调
}

// 检查是否启动完成
func checkStartDone(dbs *DBService) {

	if dbs.GetState() != 1 {
		dbs.LogWarn("checkStartDone: service not start")
		return
	}

	count := 0
	start := 0

	if dbs.option.GORMAddr != "" {
		count++
		if dbs.Gorm != nil {
			start++
		}
	}

	if dbs.option.RedisAddr != "" {
		count++
		if dbs.Redis != nil {
			start++
		}
	}
	if count == 0 {
		dbs.LogError("checkStartDone: db start in 0")
		return
	}

	if count == start {
		if dbs.option.OnStart != nil {
			dbs.option.OnStart(dbs)
		}
		dbs.StartDone()
	}
}

func NewDBService(opt DBServiceOption) *DBService {

	dbs := &DBService{
		option: opt,
	}
	name := "dbservice"
	if opt.Name != "" {
		name = fmt.Sprintf("%v-%v", name, opt.Name)
	}
	zs := zservice.NewService(name, func(s *zservice.ZService) {
		count := 0
		// gorm
		if opt.GORMAddr != "" {
			dbs.LogInfo("Init DB")
			count++
			dbs.Gorm = NewGormEX(opt)
			zservice.Go(func() { checkStartDone(dbs) })
		}

		// redis
		if opt.RedisAddr != "" {
			dbs.LogInfo("Init Redis")
			count++
			dbs.Redis = NewGoRedisEX(opt)
			zservice.Go(func() { checkStartDone(dbs) })
		}

		if count == 0 {
			dbs.LogWarn("dbservice start: not option in start")
			s.StartDone()
		}

	})

	dbs.ZService = zs
	return dbs
}

// 是否是未找到错误
func (dbs *DBService) IsNotFoundErr(e error) bool {
	return dbs.Gorm != nil && dbs.Gorm.IsNotFoundErr(e) || dbs.Redis != nil && dbs.Redis.IsNotFoundErr(e)
}

// 同步表缓存
func (dbs *DBService) SyncTableCache(ctx *zservice.Context, tabArr any, getRK func(v any) string) *zservice.Error {

	startCount := 0
	limitCount := 200 // 每次同步 200 条
	allCount := 0     // 所有查询到的数据
	errorCount := 0   // 错误数量
	for {
		// 查数据库
		if e := dbs.Gorm.Limit(limitCount).Order("created_at ASC").Offset(startCount).Find(&tabArr).Error; e != nil {
			return zservice.NewError(e)
		}

		arr := tabArr.([]any)

		// 更新缓存
		for _, v := range arr {
			allCount++
			rk_info := getRK(v)
			un, e := dbs.Redis.Lock(rk_info)
			if e != nil {
				ctx.LogError(e)
				errorCount++
			}
			defer un()

			if e := dbs.Redis.Set(rk_info, string(zservice.JsonMustMarshal(v))).Err(); e != nil {
				ctx.LogError(e)
				errorCount++
			} else {
				dbs.Redis.Expire(rk_info, zservice.Time_10Day) // 设置过期时间
			}
		}

		startCount += limitCount // 更新查询起点

		if len(arr) < limitCount { // 同步完成
			break
		}
	}

	if allCount > 0 && errorCount > 0 {
		if errorCount < allCount {
			return zservice.NewErrorf("SyncTableCache has Error, A:%v E:%v", allCount, errorCount).SetCode(zservice.Code_SyncCacheIncomplete)
		} else {
			return zservice.NewError("SyncTableCache Fail").SetCode(zservice.Code_SyncCacheErr)
		}
	} else {
		return nil
	}
}

type HasTableValueOption struct {
	Tab      any
	RK       string
	SQLConds []any
}

// 查询表中是否有指定值
func (dbs *DBService) HasTableValue(ctx *zservice.Context, opt HasTableValueOption) (bool, *zservice.Error) {

	if has, e := dbs.Redis.Exists(opt.RK).Result(); e != nil {
		return false, zservice.NewError(e)
	} else if has > 0 {
		return true, nil
	}

	dbSql := dbs.Gorm.Model(opt.Tab)

	if len(opt.SQLConds) > 0 {
		dbSql = dbSql.Where(opt.SQLConds[0], opt.SQLConds[1:]...)
	}

	// 验证数据库中是否存在
	count := int64(0)

	if e := dbSql.Count(&count).Error; e != nil {
		return false, zservice.NewError(e)
	}

	return count > 0, nil

}

// 获取一个新的ID
func (dbs *DBService) GetNewTableID(
	ctx *zservice.Context,
	genID func() uint32,
	verifyFN func(ctx *zservice.Context, id uint32) (bool, *zservice.Error),
) (uint32, *zservice.Error) {
	forCount := 0
	orgID := uint32(0)
	for {
		if forCount > 10 {
			return 0, zservice.NewError("gen id count max fail").SetCode(zservice.Code_GenIDCountMaxErr)
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

type GetTableValueOption struct {
	Tab      any            // 需要查询的表结构体对象
	RK       string         // redis缓存Key
	SQLConds []any          // 查询条件, 使用 Gorm 中的 Where 条件，不要自己拼接，避免 sql 注入
	Order    string         // 排序
	Expires  *time.Duration // 缓存过期时间,默认：10天, 0 表示不过期
}

// 获取指定值
// 注意，如果没找到数据回返回：zservice.Code_NotFound
func (dbs *DBService) GetTableValue(ctx *zservice.Context, opt GetTableValueOption) *zservice.Error {
	// 读缓存
	if e := dbs.Redis.GetScan(opt.RK, opt.Tab); e != nil {
		if e.GetCode() != zservice.Code_NotFound {
			return e.AddCaller()
		}
	} else {
		return nil
	}

	dbSql := dbs.Gorm.Model(opt.Tab)
	if opt.Order != "" {
		dbSql = dbSql.Order(opt.Order)
	}
	if len(opt.SQLConds) > 0 {
		dbSql = dbSql.Where(opt.SQLConds[0], opt.SQLConds[1:]...)
	}

	// 查库
	if e := dbSql.First(opt.Tab).Error; e != nil {
		if dbs.IsNotFoundErr(e) {
			return zservice.NewError(e).SetCode(zservice.Code_NotFound)
		} else {
			return zservice.NewError(e)
		}
	}

	// 更新缓存
	zservice.Go(func() {
		if e := func() error {
			if opt.Expires == nil {
				return dbs.Redis.SetEX(opt.RK, string(zservice.JsonMustMarshal(opt.Tab)), zservice.Time_10Day).Err()
			} else if opt.Expires.Abs().Milliseconds() == 0 {
				return dbs.Redis.Set(opt.RK, string(zservice.JsonMustMarshal(opt.Tab))).Err()
			} else {
				return dbs.Redis.SetEX(opt.RK, string(zservice.JsonMustMarshal(opt.Tab)), zservice.Time_10Day).Err()
			}
		}(); e != nil {
			ctx.LogError(e) // 缓存更新失败
		}
	})

	return nil
}

// 存储表数据
func (dbs *DBService) SaveTableValue(ctx *zservice.Context, tab any, rk string) *zservice.Error {

	unlock := func() {}

	if rk != "" {
		un, e := dbs.Redis.Lock(rk)
		if e != nil {
			return e.AddCaller()
		} else {
			unlock = un
		}
	}

	defer unlock()

	if e := dbs.Gorm.Save(tab).Error; e != nil {
		return zservice.NewError(e)
	}

	// 删除缓存
	if rk != "" {
		zservice.Go(func() {
			if e := dbs.Redis.Del(rk).Err(); e != nil {
				ctx.LogError(e)
			}
		})
	}

	return nil

}
