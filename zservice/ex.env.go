package zservice

import (
	"fmt"
	"os"
	"sync"

	"github.com/joho/godotenv"
)

// 环境变量缓存
var envCacheMap = &sync.Map{}

func MergeEnv(envs map[string]string) {
	for k, v := range envs {
		envCacheMap.Store(k, v)
	}
}

func initEnv(c *ZServiceConfig) {
	// .env 文件加载
	if _, err := os.Stat(".env"); !os.IsNotExist(err) {
		e := LoadFileEnv(".env") // load .env file
		if e != nil {
			LogError("load .env fail:", e)
		}
	}

	if len(c.EnvFils) > 0 { // load other env files
		for _, v := range c.EnvFils {
			e := LoadFileEnv(v)
			if e != nil {
				LogError("load env files fail:", e)
			}
		}
	}
}

func Getenv(key string) string {
	s, has := envCacheMap.Load(key)
	if has {
		return s.(string)
	}

	s = os.Getenv(key)
	if s == "" {
		return ""
	}

	maps, e := godotenv.Unmarshal(fmt.Sprintf("%s=%s", key, s))
	if e != nil {
		LogError(e)
	}
	MergeEnv(maps)
	return maps[key]
}

func SetEnv(key string, value string) {
	envCacheMap.Store(key, value)
}

func GetenvInt(key string) int {
	return Convert_StringToInt(Getenv(key))
}

func GetenvUInt32(key string) int32 {
	return int32(GetenvInt(key))
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
		return StringSplit(str, split[0], true)
	}

	return StringSplit(str, ",", true)
}

// 加载本地文件环境变量
func LoadFileEnv(envFile string) *Error {
	fi, e := os.Stat(envFile)
	if e != nil {
		return NewError(e)
	}
	if fi.Size() > 1024*1024 {
		return NewError("env file too large")
	}

	data, e := os.ReadFile(envFile)
	if e != nil {
		return NewError(e)
	}

	mpas, e := godotenv.UnmarshalBytes(data)
	if e != nil {
		return NewError(e)
	}
	MergeEnv(mpas)
	return nil
}
