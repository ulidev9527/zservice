package zservice

import (
	_ "embed"
	"os"

	"github.com/joho/godotenv"
)

//go:embed version
var Version string

// 服务
var mainService *ZService

type ZServiceConfig struct {
	Name          string   // 服务名称
	Version       string   // 服务版本
	EnvFils       []string // 环境变量文件
	RemoteEnvAddr string   // 远程环境变量地址
	RemoteEnvAuth string   // 远程环境变量鉴权码
	RemoteType    string   // 远程环境变量类型 http/https/grpc
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

	// .env 文件加载
	godotenv.Load()         // load .env file
	if len(c.EnvFils) > 0 { // load other env files
		godotenv.Load(c.EnvFils...)
	}

	// 远程环境变量加载
	if c.RemoteEnvAddr != "" {
		switch c.RemoteType {

		}
		body, e := Get(NewContext(""), c.RemoteEnvAddr, &map[string]any{"auth": c.RemoteEnvAuth}, nil)

		if e != nil {
			mainService.LogPanic(e)
		}

		envMaps, e := godotenv.UnmarshalBytes(body)
		if e != nil {
			mainService.LogPanic(e)
		}

		for k, v := range envMaps {
			e := os.Setenv(k, v)
			if e != nil {
				mainService.LogPanic(e)
			}
		}
	}
}

// 获取服务名称
func GetName() string {
	return mainService.name
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
