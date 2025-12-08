package zservice

import "time"

var Version = "1.0.0"

// 服务
var __mainService *ZService

// 获取服务上下文
func GetMainCtx() *Context {
	return GetMainService().launchCtx
}

// 获取服务名称
func GetServiceName() string {
	return GetMainService().option.Name
}

// 获取主服务
func GetMainService() *ZService {
	initMainService()
	return __mainService
}

func Start() *ZService {

	ser := GetMainService()

	ser.Start()
	ser.WaitStart()
	ser.LogInfof("[[[[[[ %s start done %s]]]]]]", GetServiceName(), time.Since(ser.createTime))

	return ser
}

func Stop() {
	if __mainService == nil {
		return
	}
	__mainService.Stop()
}

// 添加依赖服务
func AddDependService(s ...*ZService) *ZService {
	return __mainService.AddDependService(s...)
}

// 等待启动
func WaitStart() *ZService {
	return GetMainService().WaitStart()
}

// 等待停止
func WaitStop() {
	if __mainService == nil {
		return
	}
	__mainService.WaitStop()
}
