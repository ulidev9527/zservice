package zservice

import (
	_ "embed"
)

//go:embed version
var Version string

// 服务
var mainService *ZService

type ZServiceConfig struct {
	Name    string   // 服务名称
	Version string   // 服务版本
	EnvFils []string // 环境变量文件
}

func init() {
}

// zservice 初始化
func Init(c *ZServiceConfig) {

	LogInfof("zservice v%s", Version)
	LogInfof("%s v%s", c.Name, c.Version)

	if c.Name == "" {
		LogPanic("zservice name is empty")
	}
	if c.Version == "" {
		LogPanic("zservice version is empty")
	}

	mainService = createService(c.Name, nil)
	initEnv(c)
}

// 获取服务名称
func GetServiceName() string {
	return mainService.name
}

// 获取主服务
func GetMainService() *ZService {
	return mainService
}

func Start() {
	mainService.start()
}

func Stop() {
	for i := 0; i < len(mainService.dependService); i++ {
		mainService.dependService[i].Stop()
	}
	mainService.Stop()
}

// 添加依赖服务
func AddDependService(s *ZService) {
	mainService.AddDependService(s)
}

// 等待启动
func WaitStart() {
	mainService.WaitStart()
}

// 等待停止
func WaitStop() {
	mainService.WaitStop()
}
