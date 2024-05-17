package main

import (
	"zservice/zservice"

	"zservice/zservice/ex/redisservice"
)

func init() {

	zservice.Init("redis_test", "1.0.0")
}
func main() {

	redisS := redisservice.NewRedisService(&redisservice.RedisServiceConfig{
		Addr: zservice.Getenv("REDIS_ADDR"),
		Pass: zservice.Getenv("REDIS_PASS"),
		OnStart: func(db *redisservice.GoRedisEX) {

			zservice.TestAction("get test", func() {

				has, e := db.Exists("key").Result()
				if e != nil {
					zservice.LogError(e)
				}
				zservice.LogInfo(has)

				s, e := db.Get("key").Result()
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
