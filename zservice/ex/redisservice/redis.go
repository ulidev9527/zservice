package redisservice

import (
	"context"
	"fmt"
	"zservice/zservice"

	"github.com/redis/go-redis/v9"
)

type RedisService struct {
	*zservice.ZService
	Redis *GoRedisEX
}

// Redis 配置
type RedisServiceConfig struct {
	Name          string           // 服务名称,仅用于日志显示
	Addr          string           // redis 连接地址
	Pass          string           // redis 连接密码
	KeyPrefix     string           // 前缀
	IgnorePrefix  bool             // 是否忽略前缀 默认 false
	KeyLockPrefix string           // 锁前缀 默认:"zservice:keylock:"
	OnStart       func(*GoRedisEX) // 启动的回调
}

// 创建一个 redis 服务
func NewRedisService(c *RedisServiceConfig) *RedisService {

	if c == nil {
		zservice.LogPanic("RedisServiceConfig is nil")
		return nil
	}

	name := "RedisService"
	if c.Name != "" {
		name = fmt.Sprint(name, "-", c.Name)
	}

	rs := &RedisService{}

	rs.Redis = &GoRedisEX{
		client: redis.NewClient(&redis.Options{
			Addr:     c.Addr,
			Password: c.Pass,
		}),
		keyPrefix:       c.KeyPrefix,
		ignoreKeyPrefix: c.IgnorePrefix,
		keyLockPrefix:   c.KeyLockPrefix,
	}

	// key前缀处理
	if !c.IgnorePrefix && rs.Redis.keyPrefix == "" {
		if rs.Redis.keyPrefix == "" { // 默认使用服务名
			rs.Redis.keyPrefix = fmt.Sprint(zservice.GetServiceName(), ":")
		}
		rs.Redis.keyPrefix = fmt.Sprint(c.KeyPrefix, ":")
	}
	// 锁前缀处理
	if rs.Redis.keyLockPrefix == "" {
		rs.Redis.keyLockPrefix = "zservice:keylock:"
	}

	rs.ZService = zservice.NewService(name, func(s *zservice.ZService) {

		_, e := rs.Redis.client.Info(context.TODO(), "stats").Result()
		if e != nil {
			zservice.LogPanic(e)
		}

		if c.OnStart != nil {
			c.OnStart(rs.Redis)
		}

		s.StartDone()
	})
	return rs
}
