package zservice

import "time"

const (
	Time_0     = time.Duration(0)      // 0 秒
	Time_10ms  = time.Millisecond * 10 // 10 毫秒
	Time_1s    = time.Second           // 1s
	Time_1m    = time.Minute           // 1 分钟
	Time_10m   = time.Minute * 10      // 10 分钟
	Time_10Day = time.Hour * 24 * 10   // 10 天
	Time_3Day  = time.Hour * 24 * 3    // 3 天

)
