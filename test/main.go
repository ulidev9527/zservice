package main

import (
	"encoding/json"
	"os"
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
	bt, e := os.ReadFile("main.go")
	if e != nil {
		zservice.LogError(e)
		return
	}

	fi, e := os.Stat("main.go1")
	if e != nil {
		os.IsNotExist(e)
		zservice.LogError(e)
	} else {
		zservice.LogInfo(fi.Name())
	}

	zservice.LogInfo(zservice.MD5String(string(bt)))
	zservice.LogInfo(zservice.Md5File("main.go"))
	zservice.LogInfo(zservice.Md5Bytes(bt))

	t := TT{
		ID:       "1",
		Name:     "2",
		NickName: "3",
		BT:       bt,
	}

	s := string(zservice.JsonMustMarshal(t))

	zservice.LogInfo(s)

	a := TT{}
	json.Unmarshal([]byte(s), &a)
	zservice.LogInfo(a)

}
