package test

import (
	"testing"
)

type Benchmark_Map_100wObj struct {
}

func (obj *Benchmark_Map_100wObj) Run() {

}

var benchmark_Map_100w_maps = make(map[int]*Benchmark_Map_100wObj)
var benchmark_Arr_100w_arr = make([]*Benchmark_Map_100wObj, 1000000)

func init() {
	for i := 0; i < 1000000; i++ {
		benchmark_Map_100w_maps[i] = &Benchmark_Map_100wObj{}
		benchmark_Arr_100w_arr[i] = &Benchmark_Map_100wObj{}
	}
}

func Benchmark_Map_100w(b *testing.B) {

	// 测试结果
	// Benchmark_Map_100w-12    	1000000000	         0.01536 ns/op	       0 B/op	       0 allocs/op
	// PASS
	// ok  	zserviceapps/test	0.879s

	for _, n := range benchmark_Map_100w_maps {
		n.Run()
	}
}

func Benchmark_Arr_100w(b *testing.B) {
	// 测试结果
	// Benchmark_Arr_100w-12    	1000000000	         0.0003346 ns/op	       0 B/op	       0 allocs/op
	// PASS
	// ok  	zserviceapps/test	0.766s

	for _, n := range benchmark_Arr_100w_arr {
		n.Run()
	}
}
