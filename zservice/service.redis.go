package zservice

import (
	"fmt"

	"github.com/redis/go-redis/v9"
)

type RedisService struct {
	*ZService
	Redis *redis.Client
}

// Redis 配置
type RedisServiceConfig struct {
	Name          string
	Addr          string
	Pass          string
	OnBeforeStart func(*RedisService) // 启动前的回调
}

// 创建一个 redis 服务
func NewRedisService(c *RedisServiceConfig) *RedisService {
	if c == nil {
		LogError("ZServiceRESTConfig is nil")
		return nil
	}
	name := c.Name
	if name == "" {
		name = "DEF"
	}

	r := redis.NewClient(&redis.Options{
		Addr:     c.Addr,
		Password: c.Pass,
	})

	zs := NewService(&ZServiceConfig{
		Name: fmt.Sprint("RedisService-", name),
		OnBeforeStart: func(s *ZService) {
			if c.OnBeforeStart != nil {
				c.OnBeforeStart(&RedisService{
					ZService: s,
					Redis:    r,
				})
			}
		},
	})

	rs := &RedisService{
		ZService: zs,
		Redis:    r,
	}
	return rs
}

func (rs *RedisService) Get(ctx *ZContext, key string) (v string, e error) {
	return rs.Redis.Get(ctx, key).Result()
}
