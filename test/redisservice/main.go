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

				rk := "getKey"
				s, e := db.Set(rk, "value", 0).Result()
				if e != nil {
					zservice.LogError(e)
				}
				zservice.LogInfo(s)

				has, e := db.Exists(rk).Result()
				if e != nil {
					zservice.LogError(e)
				}
				zservice.LogInfo(has)

				s, e = db.Get(rk).Result()
				if e != nil {
					zservice.LogError(e)
				}
				zservice.LogInfo(s)
			})

			zservice.TestAction("hset test", func() {
				rk := "hsetKey"
				var maps = &struct {
					ID   uint
					Name string
				}{
					ID:   10,
					Name: "dddd",
				}
				zservice.LogInfo(maps)
				zservice.LogInfo(string(zservice.JsonMustMarshal(&maps)))

				if e := db.HSet(rk, zservice.JsonMustUnmarshal_MapAny(zservice.JsonMustMarshal(maps))).Err(); e != nil {
					zservice.LogError(e)
				}

			})

		},
	})

	zservice.AddDependService(redisS.ZService)

	zservice.Start()

}
