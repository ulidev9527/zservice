package zservice

import "time"

func TestAction(name string, cb func()) {
	LogInfo("TestAction start", name)
	t := time.Now()

	cb()

	LogInfo("TestAction end", name, time.Since(t))
}
