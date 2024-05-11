package main

import (
	"zservice/zservice"

	"zservice/zservice/ex/redisservice"

	"github.com/redis/go-redis/v9"
)

func init() {

	zservice.Init(&zservice.ZServiceConfig{
		Name:    "redis_test",
		Version: "1.0.0",
	})
}
func main() {

	redisS := redisservice.NewRedisService(&redisservice.RedisServiceConfig{
		Addr: zservice.Getenv("REDIS_ADDR"),
		Pass: zservice.Getenv("REDIS_PASS"),
		OnStart: func(db *redis.Client) {

			zservice.TestAction("get test", func() {

				has, e := db.Exists(zservice.TODO(), "key").Result()
				if e != nil {
					zservice.LogError(e)
				}
				zservice.LogInfo(has)

				s, e := db.Get(zservice.TODO(), "key").Result()
				if e != nil {
					zservice.LogError(e)
				}
				zservice.LogInfo(s)
			})

		},
	})

	zservice.AddDependService(redisS.ZService)

	zservice.Start()
	zservice.WaitStop()

}
