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

				if e := db.SetEX(rk2, string(zservice.JsonMustMarshal(v)), zglobal.Time_1m).Err(); e != nil {
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

				if e := db.SetEX(rk1, string(zservice.JsonMustMarshal(v)), zglobal.Time_1m).Err(); e != nil {
					zservice.LogError(e)
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

			zservice.TestAction("lpush", func() {
				rk := "lpushKey"
				if e := db.LPush(rk, 1, 2, 3, 4).Err(); e != nil {
					zservice.LogError(e)
				}

				if e := db.LPush(rk, zservice.JsonMustMarshalString(&MapNNN{
					ID:   111,
					Name: "nnnn",
					TTT:  time.Now(),
				})).Err(); e != nil {
					zservice.LogError(e)
				}

			})

			zservice.TestAction("lrange", func() {
				rk := "lpushKey"
				s, e := db.LRange(rk, 0, -1).Result()
				if e != nil {
					zservice.LogError(e)
				}
				zservice.LogInfo(s)
			})
			zservice.TestAction("rpop", func() {
				rk := "lpushKey"

				for i := 0; i < 5; i++ {

					s, e := db.RPop(rk).Result()
					if e != nil {
						zservice.LogError(e)
					}
					zservice.LogInfo(s)
					time.Sleep(time.Millisecond * 100)
				}
			})

			zservice.TestAction("scan", func() {

				if keys, index, e := db.Scan(0, "*", 1000).Result(); e != nil {
					zservice.LogError(e)
				} else {
					zservice.LogInfo(keys, index)
				}

				if keys, index, e := db.ScanType(0, "*", 1000, "string").Result(); e != nil {
					zservice.LogError(e)
				} else {
					zservice.LogInfo(keys, index)
				}

			})

			zservice.TestAction("setnx", func() {

				rk := "setnxKey"
				if has, e := db.SetNX(rk, "value", zglobal.Time_1m).Result(); e != nil {
					zservice.LogError(e)
				} else {
					zservice.LogInfo("setnx", has)
				}

			})

			zservice.TestAction("setex", func() {

				rk := "setexKey"

				if e := db.SetEX(rk, "value", zglobal.Time_1m).Err(); e != nil {
					zservice.LogError(e)
				}

				if e := db.SetEX(rk, "value", zglobal.Time_1m).Err(); e != nil {
					zservice.LogError(e)
				}

			})

		},
	})

	zservice.AddDependService(redisS.ZService)

	zservice.Start()

}
