package zservice

import "sync"

var wg = sync.WaitGroup{} // 线程阻断
var mu = sync.Mutex{}     // 互斥锁

var WaitCounter = 0 // 阻断计数
// 添加一个异步阻断
func WaitGroupAdd() {
	WaitCounter++
	wg.Add(1)
}

// 取消一个异步阻断
func WaitGroupRemove() {
	WaitCounter--
	wg.Done()
}

// 等待异步完成
func WaitGroup() {
	wg.Wait()
}

// 互斥锁上锁
func MutexLock() {
	mu.Lock()
}

// // 互斥锁尝试上锁 不建议使用
// func MutexTryLock() bool {
// 	return mu.TryLock()
// }

// 互斥锁解锁
func MutexUnlock() {
	mu.Unlock()
}
