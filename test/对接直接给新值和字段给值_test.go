package test

import "testing"

type Benchmark_SetValue struct {
	Value int
}

var Benchmark_SetValueList []Benchmark_SetValue
var Benchmark_SetValue_LoopCount = 100000

func init() {
	Benchmark_SetValueList = make([]Benchmark_SetValue, Benchmark_SetValue_LoopCount)
	for i := 0; i < Benchmark_SetValue_LoopCount; i++ {
		Benchmark_SetValueList[i] = Benchmark_SetValue{}
	}
}

func Benchmark_SetValue_New(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for i := 0; i < Benchmark_SetValue_LoopCount; i++ {
			Benchmark_SetValueList[i] = Benchmark_SetValue{}
		}
	}
}

func Benchmark_SetValue_Field(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for i := 0; i < Benchmark_SetValue_LoopCount; i++ {
			Benchmark_SetValueList[i].Value = 0
		}
	}
}
