package test

import (
	"sync"
	"testing"
	"time"

	"github.com/ulidev9527/zservice/zservice"
)

func Benchmark_GOCoroutine_golang(t *testing.B) {
	wg := sync.WaitGroup{}
	for i := 0; i < t.N; i++ {
		wg.Add(1)
		zservice.Go(func() {
			time.Sleep(time.Millisecond * 10) // 休眠 10毫秒
			wg.Done()
		})
	}
	wg.Wait()
}

func Benchmark_GOCoroutine_ants(t *testing.B) {
	wg := sync.WaitGroup{}
	for i := 0; i < t.N; i++ {
		wg.Add(1)
		zservice.GO_ants(func() {
			time.Sleep(time.Millisecond * 10) // 休眠 10毫秒
			wg.Done()
		})
	}
	wg.Wait()
}
