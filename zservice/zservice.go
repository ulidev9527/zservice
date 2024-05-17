package zservice

import (
	"os"
)

var Version = "0.1.0"

// 服务
var mainService *ZService

// zservice 初始化
func Init(serviceName, serviceVersion string) {
	// 配置初始化环境变量
	SetEnv(ENV_ZSERVICE_NAME, serviceName)
	SetEnv(ENV_ZSERVICE_VERSION, serviceVersion)

	// 加载 .env 文件环境变量
	if _, err := os.Stat(".env"); !os.IsNotExist(err) {
		e := LoadFileEnv(".env") // load .env file
		if e != nil {
			LogError("load .env fail:", e)
		}
	}

	// 自定义其它文件配置
	func() {
		arr := GetenvStringSplit(ENV_ZSERVICE_FILES_ENV)
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
	if Getenv(ENV_ZSERVICE_REMOTE_ENV_ADDR) != "" {
		e := LoadRemoteEnv(Getenv(ENV_ZSERVICE_REMOTE_ENV_ADDR), Getenv(ENV_ZSERVICE_REMOTE_ENV_AUTH))
		if e != nil {
			LogError(e)
		}
	}

	LogInfof("run service at:   zservice v%s", Version)

	if Getenv(ENV_ZSERVICE_NAME) == "" {
		LogPanic("zservice name is empty, you need run zservice.Init first")
	}
	if Getenv(ENV_ZSERVICE_VERSION) == "" {
		LogPanic("zservice version is empty, you need run zservice.Init first")
	}

	LogInfof("run service up:   %s v%s", serviceName, Getenv(ENV_ZSERVICE_VERSION))
	LogInfof("reg service name: %s", Getenv(ENV_ZSERVICE_NAME))

	mainService = createService(Getenv(ENV_ZSERVICE_NAME), nil)
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
	return mainService.start()
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
