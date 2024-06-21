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
	BT       []byte `json:"bt"`
}

func main() {

	a := uint32(100)

	t1 := time.Now()

	t2 := t1.Add(time.Second * time.Duration(a))

	zservice.LogInfo(time.Until(t2))

}
