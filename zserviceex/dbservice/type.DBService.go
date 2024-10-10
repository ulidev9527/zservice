package dbservice

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ulidev9527/zservice/zservice"
)

// 数据服务，结合：github.com/redis/go-redis/v9 和 gorm.io/gorm 配合使用
type DBService struct {
	*zservice.ZService
	option DBServiceOption // 配置
	Gorm   *GormEX         // 数据库 gorm
	Redis  *GoRedisEX      // 缓存 redis
}

type DBServiceOption struct {
	Name            string           // 服务名称 仅用于日志显示，如果主服务中有多个DBService，建议配置，如果只有一个DBService，可以忽略此配置
	DBType          string           // 数据库 类型 目前支持 mysql/postgres
	DBName          string           // 数据库 名称
	DBHost          string           // 数据库 地址 填入地址才会启用 Gorm 功能
	DBPort          int              // 数据库 端口
	DBUser          string           // 数据库 用户名
	DBPass          string           // 数据库 密码
	DBParams        string           // 数据库 额外参数
	MaxIdleConns    int              // 最大空闲连接数 default: 10
	MaxOpenConns    int              // 最大连接数 default: 30
	ConnMaxLifetime float32          // 连接最大生命周期 default: 300s
	RedisAddr       string           // redis 地址 填入地址才会启用 Redis 功能
	RedisPass       string           // redis 密码
	RedisPrefix     string           // redis 前缀 默认使用 zservice.Init 中的 serviceName
	Debug           bool             // 是否开启 debug
	OnStart         func(*DBService) // 启动的回调
}

func NewDBService(opt DBServiceOption) *DBService {

	dbs := &DBService{
		option: opt,
	}

	// 处理配置
	name := "dbservice"
	if opt.Name != "" {
		name = fmt.Sprintf("%v-%v", name, opt.Name)
	}
	if opt.MaxIdleConns == 0 {
		opt.MaxIdleConns = 10
	}

	if opt.MaxOpenConns == 0 {
		opt.MaxOpenConns = 30
	}
	if opt.ConnMaxLifetime == 0 {
		opt.ConnMaxLifetime = 3
	}
	zs := zservice.NewService(name, func(s *zservice.ZService) {

		// gorm
		if opt.DBHost != "" {
			dbs.Gorm = NewGormEX(opt)
		}

		// redis
		if opt.RedisAddr != "" {
			dbs.Redis = NewGoRedisEX(opt)
		}
		s.StartDone()

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
			return zservice.NewErrorf("SyncTableCache has Error, A:%v E:%v", allCount, errorCount)
		} else {
			return zservice.NewError("SyncTableCache Fail")
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
			return 0, zservice.NewError("gen id count max fail")
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
	Tab        any                                                          // 需要查询的表结构体对象
	RK         string                                                       // redis缓存Key
	SQLConds   []any                                                        // 查询条件, 使用 Gorm 中的 Where 条件，不要自己拼接，避免 sql 注入
	Order      string                                                       // 排序
	Expires    *time.Duration                                               // 缓存过期时间,默认：10天, 0 表示不过期
	RedisGetFN func(*zservice.Context, any, string) (bool, *zservice.Error) // 获取缓存数据，需要对数据进行重新处理, 返回处理状态或者错误消息，处理后的数据为空表示处理失败，会继续进行数据查询
	RedisSetFN func(*zservice.Context, any) string                          // 存储缓存数据，返回处理后的字符串，字符串将存储到 redis
}

// 获取指定值
// 注意，如果没找到数据回返回：zservice.Code_NotFound
func (dbs *DBService) GetTableValue(ctx *zservice.Context, opt GetTableValueOption) *zservice.Error {
	// 读缓存
	if opt.RK != "" { // 需要读缓存
		if s, e := dbs.Redis.Get(opt.RK).Result(); e != nil {
			if !dbs.Redis.IsNotFoundErr(e) {
				return zservice.NewError(e).SetCode(zservice.Code_Fatal)
			}
		} else {
			if opt.RedisGetFN != nil { // 使用提供的数据处理方法
				if ok, e := opt.RedisGetFN(ctx, opt.Tab, s); e != nil {
					return e.AddCaller()
				} else if ok { // 处理成功，直接返回，失败回去数据库进行查询
					return nil
				}

				// 处理失败，数据库查询
			} else { // 使用默认的数据处理方法
				if e := json.Unmarshal([]byte(s), opt.Tab); e != nil { // 处理失败，删除缓存，继续数据查询
					ctx.LogError(e, opt)
				} else {
					return nil
				}
			}

			// 执行到此次代表缓存查询的数据有问题，删除缓存
			zservice.Go(func() {
				if e := dbs.Redis.Del(opt.RK); e != nil {
					ctx.LogError(e)
				}
			})
		}
	}

	// 数据库查询

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
		if opt.RK == "" {
			return
		}
		if e := func() error {
			val := ""
			if opt.RedisSetFN != nil {
				val = opt.RedisSetFN(ctx, opt.Tab)
			} else {
				val = string(zservice.JsonMustMarshal(opt.Tab))
			}
			if opt.Expires == nil {
				return dbs.Redis.SetEX(opt.RK, val, zservice.Time_10Day).Err()
			} else if opt.Expires.Abs().Milliseconds() == 0 {
				return dbs.Redis.Set(opt.RK, val).Err()
			} else {
				return dbs.Redis.SetEX(opt.RK, val, zservice.Time_10Day).Err()
			}
		}(); e != nil {
			ctx.LogError(e) // 缓存更新失败
		}
	})

	return nil
}

type SaveTableValueOption struct {
	Tab any    // 表
	RK  string // 缓存的 Redis
}

// 存储表数据
func (dbs *DBService) SaveTableValue(ctx *zservice.Context, opt SaveTableValueOption) *zservice.Error {

	unlock := func() {}

	if opt.RK != "" {
		un, e := dbs.Redis.Lock(opt.RK)
		if e != nil {
			return e.AddCaller()
		} else {
			unlock = un
		}
	}

	defer unlock()

	if e := dbs.Gorm.Save(opt.Tab).Error; e != nil {
		return zservice.NewError(e)
	}

	// 删除缓存
	if opt.RK != "" {
		zservice.Go(func() {
			if e := dbs.Redis.Del(opt.RK).Err(); e != nil {
				ctx.LogError(e)
			}
		})
	}

	return nil

}
