package test

import (
	"testing"
)

///  测试结果 - 相差不大

type IBenchmark_AnyToType_Item interface {
	Update()
}

type Benchmark_AnyToType_Item struct {
}

func (b *Benchmark_AnyToType_Item) Update() {

}

func (b *Benchmark_AnyToType_Item) GetType() int {
	return 1
}

var Benchmark_AnyToType_Item_IMap = map[int]IBenchmark_AnyToType_Item{}
var Benchmark_AnyToType_Item_Map = map[int]*Benchmark_AnyToType_Item{}
var Benchmark_AnyToType_Count = 100000

func init() {

	for i := range Benchmark_AnyToType_Count {
		Benchmark_AnyToType_Item_IMap[i] = &Benchmark_AnyToType_Item{}
	}

	for i := range Benchmark_AnyToType_Count {
		Benchmark_AnyToType_Item_Map[i] = &Benchmark_AnyToType_Item{}
	}

}

func Benchmark_AnyToType1(b *testing.B) {

	for i := 0; i < b.N; i++ {
		for _, item := range Benchmark_AnyToType_Item_IMap {
			iib := item.(*Benchmark_AnyToType_Item)
			iib.Update()
			iib.GetType()
		}
	}

}

func Benchmark_AnyToType2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, item := range Benchmark_AnyToType_Item_Map {
			item.Update()
			item.GetType()
		}
	}

}
