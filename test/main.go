package main

import (
	"zservice/zservice"
)

func init() {
	zservice.Init("test", "1.0.0")
}

func main() {

	e1 := zservice.NewError("123").SetCode(100)

	zservice.LogError(e1)

	e2 := zservice.NewErrore(e1)
	zservice.LogError(e2)

	e2.AddCaller()
	zservice.LogError(e2)
	e3 := getE()
	e3.AddCaller()
	zservice.LogError(e3)

}
