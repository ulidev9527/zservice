package main

import (
	"fmt"
	"zservice/zservice"
)

func init() {
	zservice.Init(&zservice.ZServiceConfig{
		Name:    "test",
		Version: "1.0.0",
	})
}

func main() {
	zservice.TestAction("test1", func() {
		for i := 0; i < 200000; i++ {
			zservice.Convert_IntToString(zservice.RandomIntRange(100000, 999999))
		}
	})
	zservice.TestAction("test2", func() {
		for i := 0; i < 200000; i++ {
			zservice.RandomIntRange(100000, 999999)
			// if len(strconv.Itoa(c)) != 6 {
			// 	zservice.LogError(c)
			// }
		}
	})
	zservice.TestAction("test3", func() {
		for i := 0; i < 200000; i++ {
			fmt.Sprintf("%d", zservice.RandomIntRange(100000, 999999))
			// if len(strconv.Itoa(c)) != 6 {
			// 	zservice.LogError(c)
			// }
		}
	})
	zservice.TestAction("test3", func() {
		for i := 0; i < 200000; i++ {
			zservice.Sprint("%d", zservice.RandomIntRange(100000, 999999))
			// if len(strconv.Itoa(c)) != 6 {
			// 	zservice.LogError(c)
			// }
		}
	})

}
