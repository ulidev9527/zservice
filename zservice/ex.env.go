package zservice

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

// 环境变量缓存
var envCacheMap = map[string]string{}

func ClearEnvCache() {
	for k := range envCacheMap {
		delete(envCacheMap, k)
	}

}

func initEnv(c *ZServiceConfig) {
	ClearEnvCache()
	// .env 文件加载
	if _, err := os.Stat(".env"); !os.IsNotExist(err) {
		e := godotenv.Load() // load .env file
		if e != nil {
			LogError("load .env fail:", e)
		}
	}

	if len(c.EnvFils) > 0 { // load other env files
		e := godotenv.Load(c.EnvFils...)
		if e != nil {
			LogError("load env files fail:", e)
		}
	}

	// 远程环境变量加载
	if c.RemoteEnvAddr != "" {
		LoadRemoteEnv(c.RemoteEnvAddr, c.RemoteEnvAuth)
	}
}

func Getenv(key string) string {
	s := envCacheMap[key]
	if s != "" {
		return s
	}
	s = os.Getenv(key)
	if s == "" {
		return s
	}

	maps, e := godotenv.Unmarshal(fmt.Sprintf("%s=%s", key, s))
	if e != nil {
		LogError(e)
	}
	for k, v := range maps {
		envCacheMap[k] = v
	}

	return Getenv(key)
}

func GetenvInt(key string) int {
	return Convert_StringToInt(Getenv(key))
}

func GetenvBool(key string) bool {
	return Convert_StringToBoolean(Getenv(key))
}

// json
func GetenvStringSplit(key string, split ...string) []string {
	str := Getenv(key)
	if str == "" {
		return []string{}
	}

	if len(split) > 0 {
		return strings.Split(str, split[0])
	}

	return strings.Split(str, ",")
}

// 加载远程环境变量
func LoadRemoteEnv(addr string, auth string) {
	body, e := Get(NewContext(""), addr, &map[string]any{"auth": auth}, nil)

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
			continue
		}
		envCacheMap[k] = v
	}
}
