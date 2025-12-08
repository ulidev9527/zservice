package zservice

import "time"

// 帧更新
func RunFrameUpdate(fps int, updateFunc func(dt int64)) {
	Go(func() {
		frameDuration := time.Second / time.Duration(fps)
		lastTime := time.Now()

		for {
			currentTime := time.Now()
			deltaTime := currentTime.Sub(lastTime).Milliseconds() // 计算时间增量（秒）
			lastTime = currentTime

			updateFunc(deltaTime) // 执行每帧逻辑

			// 控制帧率
			elapsed := time.Since(currentTime)
			if elapsed < frameDuration {
				time.Sleep(frameDuration - elapsed)
			}
		}
	})
}
