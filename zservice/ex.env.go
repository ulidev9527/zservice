package zservice

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"zservice/service/zauth/zauth_pb"

	"github.com/joho/godotenv"
)

// 环境变量缓存
var envCacheMap = &sync.Map{}

func init() {
	// 获取所有环境变量
	strArr := os.Environ()
	for _, v := range strArr {
		arr := strings.Split(v, "=")
		SetEnv(arr[0], arr[1])
	}
}

// 合并环境变量
func MergeEnv(envs map[string]string) {
	for k, v := range envs {
		SetEnv(k, v)
	}
}

// 获取环境变量，key 不限制大小写
func Getenv(key string) string {
	key = strings.ToUpper(key)
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
	key = strings.ToUpper(key)
	if key == "ZSERVICE_VERSION" && Getenv(key) != "" {
		return
	}
	envCacheMap.Store(key, value)
}

func GetenvInt(key string) int {
	return StringToInt(Getenv(key))
}

func GetenvUInt32(key string) int32 {
	return int32(GetenvInt(key))
}
func GetenvBool(key string) bool {
	return StringToBoolean(Getenv(key))
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

// 加载远程环境变量
func LoadRemoteEnv(addr string) *Error {

	if addr == "" {
		return NewError("no remote env addr")
	}

	body, e := Get(NewContext(), addr, nil, nil)
	if e != nil {
		return e
	}

	res := zauth_pb.ConfigGetServiceEnvConfig_RES{}
	if e := json.Unmarshal(body, &res); e != nil {
		return NewError(e)
	}

	return LoadStringEnv(res.Value)
}

// 加载字符串中的环境变量
func LoadStringEnv(envStr string) *Error {
	envMaps, e := godotenv.Unmarshal(envStr)
	if e != nil {
		return NewError(e)
	}
	MergeEnv(envMaps)
	return nil

}
