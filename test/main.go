package main

import (
	"time"
	"zservice/zservice"
	"zservice/zservice/zglobal"
)

func init() {
	zservice.Init("test", "1.0.0")
}

func main() {

	a := time.Now()
	b := time.Now().Add(zglobal.Time_10Day)

	if a.Before(b) {
		zservice.LogError(1)
	} else {
		zservice.LogError(2)
	}

}
