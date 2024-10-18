package test

import (
	"testing"

	"github.com/ulidev9527/zservice/zservice"
)

func Benchmark_RandomInt(t *testing.B) {
	for i := 0; i < t.N; i++ {
		zservice.RandomInt(10)
	}
}

func Benchmark_RandomString(t *testing.B) {
	for i := 0; i < t.N; i++ {
		zservice.RandomString(32)
	}
}

func Benchmark_RandomMd5(t *testing.B) {
	for i := 0; i < t.N; i++ {
		zservice.RandomMD5()
	}
}

func Benchmark_RandomMd5_XID(t *testing.B) {
	for i := 0; i < t.N; i++ {
		zservice.RandomMD5_XID()
	}
}
func Benchmark_RandomMd5_XID_Random(t *testing.B) {
	for i := 0; i < t.N; i++ {
		zservice.RandomMD5_XID_Random()
	}
}

func Test_RandomString(t *testing.T) {
	for i := 0; i < 10; i++ {
		t.Log(zservice.RandomString(32))
	}
}
