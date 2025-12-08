package goredisservice

import (
	"time"
	"zserviceapps/packages/zservice"

	"github.com/redis/go-redis/v9"
)

// 数据服务，结合：github.com/redis/go-redis/v9 和 gorm.io/gorm 配合使用
type Service struct {
	zservice  *zservice.ZService
	GoRedisEx *GoRedisEX
}

type ServiceOption struct {
	MaxIdleConns    int            // 最大空闲连接数 default: 10
	MaxOpenConns    int            // 最大连接数 default: 30
	ConnMaxLifetime float32        // 连接最大生命周期 default: 300s
	RedisAddr       string         // redis 地址 填入地址才会启用 Redis 功能
	RedisPass       string         // redis 密码
	OnStart         func(*Service) // 启动的回调
}

// 同步初始化参数
func syncDefaultOption(opt ServiceOption) ServiceOption {
	if opt.MaxIdleConns == 0 {
		opt.MaxIdleConns = 10
	}
	if opt.MaxOpenConns == 0 {
		opt.MaxOpenConns = 30
	}
	if opt.ConnMaxLifetime == 0 {
		opt.ConnMaxLifetime = 3
	}
	return opt
}

func NewService(opt ServiceOption) *Service {
	opt = syncDefaultOption(opt)
	ser := &Service{}

	ser.zservice = zservice.NewService(zservice.ServiceOptions{
		Name: "GoRedis-" + opt.RedisAddr,
		OnStart: func(_ *zservice.ZService) {
			ser.GoRedisEx = &GoRedisEX{
				Client: redis.NewClient(&redis.Options{
					Addr:         opt.RedisAddr,
					Password:     opt.RedisPass,
					MaxIdleConns: opt.MaxIdleConns,
					PoolSize:     opt.MaxOpenConns,
					PoolTimeout:  time.Duration(opt.ConnMaxLifetime) * time.Second,
				}),
			}

			for {
				_, e := ser.GoRedisEx.Info(ser.zservice.NewContext(), "stats").Result()
				if e != nil {
					ser.zservice.LogError("has error, waiting 5s", e)
					time.Sleep(time.Second * 5)
					continue
				}

				break
			}

			if opt.OnStart != nil {
				opt.OnStart(ser)
			}
		},
	})

	return ser
}

func (ser *Service) GetZService() *zservice.ZService { return ser.zservice }
