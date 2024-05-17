package main

import (
	"fmt"
	"reflect"
	"zservice/zservice"
)

func init() {
	zservice.Init("test", "1.0.0")
}

func main() {

	defer func() {
		zservice.LogInfo("1+1 === 0")
	}()
	zservice.LogInfo("111")
	if 1+1 > 0 {
		zservice.LogInfo("222")
		defer func() {
			zservice.LogInfo("1+1 > 0")
		}()
		zservice.LogInfo("333")
	}
	zservice.LogInfo("444")

}

func GetEle(v any) {

	vt := reflect.TypeOf(v).Elem()
	fmt.Println("Array element type:", vt)
}
