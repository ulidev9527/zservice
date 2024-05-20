package main

import (
	"encoding/json"
	"time"
	"zservice/zservice"
	"zservice/zservice/zglobal"

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
			type MapNNN struct {
				ID   uint
				Name string
				TTT  time.Time
			}

			zservice.TestAction("setAndSetex", func() {

				rk1 := "setAndSetex"
				rk2 := "setAndSetex2"
				v := MapNNN{
					ID:   111,
					Name: "nnnn",
					TTT:  time.Now(),
				}

				if e := db.Set(rk1, string(zservice.JsonMustMarshal(v))).Err(); e != nil {
					zservice.LogError(e)
				}

				if e := db.SetEx(rk2, string(zservice.JsonMustMarshal(v)), zglobal.Time_10Day).Err(); e != nil {
					zservice.LogError(e)
				}

				v2 := &MapNNN{}

				if e := db.Get(rk1).Scan(v2); e != nil {
					zservice.LogError(e)
				} else {
					zservice.LogInfo(v2)
				}

				if e := db.GetScan(rk2, v2); e != nil {
					zservice.LogError(e)
				} else {
					zservice.LogInfo(v2)
				}

			})

			zservice.TestAction("get test", func() {

				rk := "getKey"
				s, e := db.Set(rk, "value").Result()
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
				var maps = MapNNN{
					ID:   10,
					Name: "dddd",
					TTT:  time.Now(),
				}
				zservice.LogInfo(maps)
				zservice.LogInfo(string(zservice.JsonMustMarshal(&maps)))

				if e := db.HSet(rk, zservice.JsonMustUnmarshal_MapAny(zservice.JsonMustMarshal(maps))).Err(); e != nil {
					zservice.LogError(e)
				}
			})

			zservice.TestAction("hgetall", func() {
				rk := "hsetKey"

				mapn := MapNNN{}

				if maps, e := db.HGetAll(rk).Result(); e != nil {
					zservice.LogError(e)
				} else {
					if e := json.Unmarshal(zservice.JsonMustMarshal(maps), &mapn); e != nil {
						zservice.LogError(e)
					}
				}

				zservice.LogInfo(mapn)
			})

		},
	})

	zservice.AddDependService(redisS.ZService)

	zservice.Start()

}
