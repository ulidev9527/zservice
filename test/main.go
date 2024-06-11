package main

import (
	"time"
	"zservice/zservice"
)

func init() {
	zservice.Init("test", "1.0.0")
}

type TT struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	NickName string `json:"nickName"`
}

func main() {
	var num uint32 = 100
	negativeNum := -30

	// 将负数转换为补码形式

	num -= uint32(-negativeNum)

	zservice.LogInfo(num)

	zservice.LogInfo(time.Until(time.Now().Add(time.Second * 60)).Seconds())

}
