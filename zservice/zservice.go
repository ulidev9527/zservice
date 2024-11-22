package zservice

import (
	"fmt"
	"os"
)

var Version = "0.2.1"
var ISDebug = false

// 服务
var mainService *ZService

// zservice 初始化
func Init(opt ZserviceOption) {
	fmt.Println("zservice init start")

	initLogger()
	initEnv()

	// 加载 .env 文件环境变量
	if _, err := os.Stat(".env"); !os.IsNotExist(err) {
		e := LoadFileEnv(".env") // load .env file
		if e != nil {
			LogError("load .env fail:", e)
		}
	}

	mainService = createService(opt)

	LogInfof("run service at:    zservice v%s", Version)
	LogInfof("run service up:    %s v%s", mainService.opt.Name, mainService.opt.Version)
	LogInfof("run service name:  %s", Getenv("ZSERVICE_NAME"))

	mainService.tranceName = mainService.opt.Name
	ISDebug = mainService.opt.Debug
	LogInfo("run service debug:", BoolToString(ISDebug))
}

// 获取服务名称
func GetServiceName() string {
	return mainService.opt.Name
}

// 获取主服务
func GetMainService() *ZService {
	return mainService
}

func Start() *ZService {

	if mainService == nil {
		LogPanic("you need use zservice.Init method first")
		return nil
	}

	return mainService.Start()
}

func Stop() {
	for i := 0; i < len(mainService.dependService); i++ {
		mainService.dependService[i].Stop()
	}
	mainService.Stop()
}

// 添加依赖服务
func AddDependService(s ...*ZService) *ZService {
	return mainService.AddDependService(s...)
}

// 等待启动
func WaitStart() *ZService {
	return mainService.WaitStart()
}

// 等待停止
func WaitStop() {
	for i := 0; i < len(mainService.dependService); i++ {
		mainService.dependService[i].WaitStop()
	}
	mainService.WaitStop()
}
