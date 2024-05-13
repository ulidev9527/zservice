package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"zservice/zservice"
)

func init() {
	zservice.Init(&zservice.ZServiceConfig{
		Name:    "test",
		Version: "1.0.0",
	})
}

func main() {

	arr1 := []int{1, 2, 3, 456}
	GetEle(arr1)

	type CC struct {
		ID         string `json:"id"`
		Name       string `json:"name"`
		Desc       string `json:"desc"`
		Icon       string `json:"icon"`
		LimitCount uint32 `json:"limit_count"`
	}

	arr2 := []CC{}
	GetEle(arr2)
	GetEle(&arr2)

	str := `[ {"id":"5"},{"id":"8"},{"def_container_id":0,"desc":"初始化创建","icon":"","id":"1","limit_count":1,"name":"角色背包","type":0},{"id":"4"},{"id":"7"},{"id":"10"},{"id":"3"},{"desc":"用户ID","id":"51","name":"用户ID"},{"id":"52","name":"用户昵称"},{"id":"53","name":"用户头像"},{"desc":"常见的货币","icon":"gold","id":"54","name":"金币"},{"def_container_id":0,"desc":"","icon":"","id":"2","limit_count":1,"name":"","type":0},{"id":"9"},{"id":"6"} ]`

	// c := &CC{}
	e := json.Unmarshal([]byte(str), &arr2)
	if e != nil {
		zservice.LogError(e)
	} else {
		zservice.LogInfo(arr2)
	}

}

func GetEle(v any) {

	vt := reflect.TypeOf(v).Elem()
	fmt.Println("Array element type:", vt)
}
