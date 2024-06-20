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
	Name      string           // 服务名称,仅用于日志显示
	Addr      string           // redis 连接地址
	Pass      string           // redis 连接密码
	KeyPrefix string           // 前缀
	OnStart   func(*GoRedisEX) // 启动的回调
}

// 创建一个 redis 服务
func NewRedisService(c *RedisServiceConfig) *RedisService {

	if c == nil {
		zservice.LogPanic("RedisServiceConfig is nil")
		return nil
	}

	keyPrefix := c.KeyPrefix
	if keyPrefix == "" {
		keyPrefix = zservice.GetServiceName()
	}

	name := fmt.Sprint("RedisService-", c.Addr, "-", keyPrefix)

	rs := &RedisService{}

	rs.Redis = &GoRedisEX{
		client: redis.NewClient(&redis.Options{
			Addr:     c.Addr,
			Password: c.Pass,
		}),
		keyPrefix:     fmt.Sprint(keyPrefix, ":"),
		keyLockPrefix: fmt.Sprint("__zserviceKeyLock:", keyPrefix, ":"),
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

func IsNilErr(e error) bool {
	return e == redis.Nil
}
