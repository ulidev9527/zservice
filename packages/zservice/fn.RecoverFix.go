package zservice

import (
	"runtime"
)

func RecoverFix() {
	if e := recover(); e != nil {
		buf := make([]byte, 2048)            // 分配缓冲区
		runtime.Stack(buf, false)            // 获取完整的调用栈
		LogError("panic in", e, string(buf)) // 记录调用栈详情
	}
}
