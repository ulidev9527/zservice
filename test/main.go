package main

import (
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

	tp := "2"
	switch tp {
	case "1", "2":
		zservice.LogInfo("12")
	case "boo", "kkk":
		zservice.LogInfo("abc")
	}

}
