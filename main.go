package main

import (
	"zserviceapps/packages/zservice"
)

const (
	A_1 = iota
)

func main() {
	zservice.NewService(zservice.ServiceOptions{
		Name:    "main",
		Version: "1.0.0",
	})

	for i := range 10 {
		zservice.LogInfo(i)
	}

}
