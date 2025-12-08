package test

import (
	"slices"
	"sync"
	"testing"

	"github.com/panjf2000/ants/v2"
)

var __ArrayDelete_Count = 10000
var __ArrayDelete_Pool, _ = ants.NewPool(__ArrayDelete_Count)
var __ArrayDelete_WG sync.WaitGroup
var __ArrayDelete_Lock sync.Mutex

type __ArrayDelete_Test struct {
	Index int
}

func (test *__ArrayDelete_Test) Update() {
	test.Index--
}

func Benchmark_ArrayDelete1(b *testing.B) {

	for i := 0; i < b.N; i++ {
		arr := []*__ArrayDelete_Test{}
		// 给 arr 添加到 1k
		for j := range __ArrayDelete_Count {
			arr = append(arr, &__ArrayDelete_Test{Index: j})
		}

		for {
			if len(arr) == 0 {
				break
			}
			arrLen := 0
			for _, r := range arr {
				r.Update()
				if r.Index > 0 {
					arr[arrLen] = r
					arrLen++
				}
			}
			arr = arr[:arrLen]
		}
	}
}

func Benchmark_ArrayDelete_GO(b *testing.B) {
	for i := 0; i < b.N; i++ {

		arr := []*__ArrayDelete_Test{}
		for j := range __ArrayDelete_Count {
			arr = append(arr, &__ArrayDelete_Test{Index: j})
		}

		for {
			if len(arr) == 0 {
				break
			}
			arrLen := 0

			__ArrayDelete_WG.Add(len(arr))
			for _, r := range arr {

				__ArrayDelete_Pool.Submit(func() {

					r.Update()

					__ArrayDelete_Lock.Lock()
					defer __ArrayDelete_Lock.Unlock()

					if r.Index > 0 {
						arr[arrLen] = r
						arrLen++
					}
					__ArrayDelete_WG.Done()
				})
			}
			__ArrayDelete_WG.Wait()
			arr = arr[:arrLen]
		}
	}
}

func Benchmark_ArrayDelete_Slices1(b *testing.B) {

	for i := 0; i < b.N; i++ {
		arr := []*__ArrayDelete_Test{}
		// 给 arr 添加到 1k
		for j := range __ArrayDelete_Count {
			arr = append(arr, &__ArrayDelete_Test{Index: j})
		}

		for {
			if len(arr) == 0 {
				break
			}

			for i := 0; i < len(arr); i++ {
				arr[i].Update()
				if arr[i].Index <= 0 {
					arr = slices.Delete(arr, i, i+1)
					i--
				}
			}
		}
	}

}

func Benchmark_ArrayDelete_Slices2(b *testing.B) {

	for i := 0; i < b.N; i++ {
		arr := []*__ArrayDelete_Test{}
		// 给 arr 添加到 1k
		for j := range __ArrayDelete_Count {
			arr = append(arr, &__ArrayDelete_Test{Index: j})
		}

		for {
			if len(arr) == 0 {
				break
			}

			for _, r := range arr {
				r.Update()
			}

			arr = slices.DeleteFunc(arr, func(r *__ArrayDelete_Test) bool {
				return r.Index <= 0
			})
		}
	}

}

func Benchmark_ArrayDelete_Slices_GO(b *testing.B) {

	for i := 0; i < b.N; i++ {
		arr := []*__ArrayDelete_Test{}
		// 给 arr 添加到 1k
		for j := range __ArrayDelete_Count {
			arr = append(arr, &__ArrayDelete_Test{Index: j})
		}

		for {
			if len(arr) == 0 {
				break
			}

			__ArrayDelete_WG.Add(len(arr))
			for _, r := range arr {
				__ArrayDelete_Pool.Submit(func() {
					r.Update()
					__ArrayDelete_WG.Done()
				})
			}
			__ArrayDelete_WG.Wait()

			arr = slices.DeleteFunc(arr, func(r *__ArrayDelete_Test) bool {
				return r.Index <= 0
			})
		}
	}

}
