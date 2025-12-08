package test

import (
	"sync"
	"testing"
	"zserviceapps/packages/zservice"
)

func init() {
	// 插入 100w 条数据

}
func TestTypeAssertionAndTypeCreationPerformance(t *testing.T) {
	// 创建一个类型断言的测试

}

func BenchmarkTypeAssertion(b *testing.B) {

	var testVector3Pool = sync.Pool{
		New: func() any { return &zservice.Vector3{} },
	}

	for i := 0; i < 1000000; i++ {
		testVector3Pool.Put(&zservice.Vector3{})
	}

	var a = 0
	// 创建一个类型断言的测试
	for i := 0; i < b.N; i++ {
		vector3 := testVector3Pool.Get().(*zservice.Vector3)
		testVector3Pool.Put(vector3)
		a++
	}
	b.Logf("a: %d", a)
}

func BenchmarkTypeCreation(b *testing.B) {

	var testVector3Pool = sync.Pool{
		New: func() any { return &zservice.Vector3{} },
	}

	// 创建一个类型创建的测试
	var a = 0
	for i := 0; i < b.N; i++ {
		vector3 := &zservice.Vector3{}
		testVector3Pool.Put(vector3)
		a++
	}
	b.Logf("a: %d", a)
}
