package test

import (
	"testing"
)

// Benchmark_pb对象创建和使用对象池效果对比-12                                	28893456	        40.42 ns/op	      48 B/op	       1 allocs/op
// Benchmark_pb对象创建和使用对象池效果对比-12                                	28845735	        42.82 ns/op	      48 B/op	       1 allocs/op
// Benchmark_pb对象创建和使用对象池效果对比-12                                	29103068	        40.32 ns/op	      48 B/op	       1 allocs/op

func Benchmark_pb对象创建和使用对象池效果对比(b *testing.B) {
	// for i := 0; i < b.N; i++ {
	// 	msg := &pb.Battle_Attr{}
	// 	msg.Reset()
	// 	// b.Log(b.N)
	// }
}

// pkg: zserviceapps/test
// cpu: Intel(R) Core(TM) i7-8750H CPU @ 2.20GHz
// Benchmark_pb对象使用对象池效果对比-12                          	56942002	        19.65 ns/op	       0 B/op	       0 allocs/op
// Benchmark_pb对象使用对象池效果对比-12                          	58975476	        20.08 ns/op	       0 B/op	       0 allocs/op
// Benchmark_pb对象使用对象池效果对比-12                          	60620725	        19.54 ns/op	       0 B/op	       0 allocs/op

func Benchmark_pb对象使用对象池效果对比(b *testing.B) {
	// pool := &sync.Pool{
	// 	New: func() any {
	// 		return &pb.Battle_Attr{}
	// 	},
	// }

	// for i := 0; i < b.N; i++ {
	// 	msg := pool.Get().(*pb.Battle_Attr)
	// 	msg.Reset()
	// 	pool.Put(msg)
	// 	// b.Log(b.N)
	// }
}
