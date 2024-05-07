package dbservice

import (
	"context"
	"fmt"
	"zservice/zservice"

	"github.com/redis/go-redis/v9"
)

type RedisService struct {
	*zservice.ZService
	Redis *redis.Client
}

// Redis 配置
type RedisServiceConfig struct {
	Name    string
	Addr    string
	Pass    string
	OnStart func(*redis.Client) // 启动的回调
}

// 创建一个 redis 服务
func NewRedisService(c *RedisServiceConfig) *RedisService {
	if c == nil {
		zservice.LogPanic("ZServiceRESTConfig is nil")
		return nil
	}
	name := "RedisService"
	if name != "" {
		name = fmt.Sprint(name, "-", c.Name)
	}

	rs := &RedisService{}

	rs.Redis = redis.NewClient(&redis.Options{
		Addr:     c.Addr,
		Password: c.Pass,
	})

	rs.ZService = zservice.NewService(name, func(s *zservice.ZService) {

		_, e := rs.Redis.Info(context.TODO(), "stats").Result()
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
