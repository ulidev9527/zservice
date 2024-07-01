package zservice

import (
	"fmt"
	"os"
)

var Version = "0.1.0"
var ISDebug = false

// 服务
var mainService *ZService

// zservice 初始化
func Init(serviceName, serviceVersion string) {
	fmt.Println("zservice init start")

	initLogger()
	initEnv()

	// 配置初始化环境变量
	if Getenv("ZSERVICE_NAME") == "" {
		Setenv("ZSERVICE_NAME", serviceName)
	}
	Setenv("ZSERVICE_VERSION", serviceVersion)

	// 加载 .env 文件环境变量
	if _, err := os.Stat(".env"); !os.IsNotExist(err) {
		e := LoadFileEnv(".env") // load .env file
		if e != nil {
			LogError("load .env fail:", e)
		}
	}

	// 自定义其它文件配置
	func() {
		arr := GetenvStringSplit("ZSERVICE_FILES_ENV")
		if len(arr) > 0 { // load other env files
			for _, v := range arr {
				e := LoadFileEnv(v)
				if e != nil {
					LogError("load env files fail:", e)
				}
			}
		}
	}()

	// 加载远程环境变量
	if Getenv("ZSERVICE_REMOTE_ENV_ADDR") != "" {
		LogInfo("load remote addr", Getenv("ZSERVICE_REMOTE_ENV_ADDR"))
		e := LoadRemoteEnv(Getenv("ZSERVICE_REMOTE_ENV_ADDR"))
		if e != nil {
			LogError(e)
		}
	}

	mainService = createService(Getenv("ZSERVICE_NAME"), nil)

	LogInfof("run service at:    zservice v%s", Version)

	if Getenv("ZSERVICE_NAME") == "" {
		LogPanic("zservice name is empty, you need run zservice.Init first")
	}
	if Getenv("ZSERVICE_VERSION") == "" {
		LogPanic("zservice version is empty, you need run zservice.Init first")
	}

	LogInfof("run service up:    %s v%s", serviceName, Getenv("ZSERVICE_VERSION"))
	LogInfof("run service name:  %s", Getenv("ZSERVICE_NAME"))

	mainService.name = Getenv("ZSERVICE_NAME")
	mainService.tranceName = mainService.name
	ISDebug = GetenvBool("ZSERVICE_DEBUG")
	LogInfo("run service debug:", BoolToString(ISDebug))
}

// 获取服务名称
func GetServiceName() string {
	return mainService.name
}

// 获取主服务
func GetMainService() *ZService {
	return mainService
}

func Start() *ZService {
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
