package main

import (
	"time"
	"zservice/zservice"
)

func init() {
	zservice.Init("test", "1.0.0")
}

func main() {

	zservice.LogInfo(time.Until(time.Now().Add(time.Second * 60)).Seconds())

}
