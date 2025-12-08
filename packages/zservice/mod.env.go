package zservice

import (
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/joho/godotenv"

	_ "net/http/pprof"
)

var envCacheMap = &sync.Map{}                              // 环境变量缓存
type WatchEnvChangeHandle func(key, newVal, oldVal string) // 环境变量改变监听
var (
	__watchEnv_handleMap_lock = &sync.RWMutex{}
	__watchEnv_handleMap      = map[string][]WatchEnvChangeHandle{} // 监听环境变量的映射表
)

func init() {
	// 获取所有环境变量
	strArr := os.Environ()
	for _, v := range strArr {
		arr := strings.Split(v, "=")
		Setenv(arr[0], arr[1])
	}

	// 加载 .env 文件环境变量
	if _, e := os.Stat(".env"); e == nil {
		// load .env file
		if e := LoadFileEnv(".env"); e != nil {
			LogError("load .env fail:", e)
		}

		// 监听文件变化
		Go(func() {

			WatchEnvChange("", func(key, newVal, oldVal string) {
				LogInfo("Env Change:", key, oldVal, "To", newVal)
			})

			md5Str, e := Md5File(".env")
			if e != nil {
				LogPanic(e)
			}
			for {
				time.Sleep(time.Second)
				if str, e := Md5File(".env"); e != nil {
					LogError(".env file watch error, stop watch,", e)
					break
				} else if md5Str != str {
					md5Str = str
					if e := LoadFileEnv(".env"); e != nil {
						LogError("load .env fail:", e)
					}
				}
			}
		})

	}

}

// 合并环境变量
func MergeEnv(envs map[string]string) {
	for k, v := range envs {
		Setenv(k, v)
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
		LogError("can`t find env:", key)
		return ""
	}

	maps, e := godotenv.Unmarshal(fmt.Sprintf("%s=%s", key, s))
	if e != nil {
		LogError(e)
	}
	MergeEnv(maps)
	return maps[key]
}

func Setenv(key string, value string) {
	key = strings.ToUpper(key)

	old := ""
	// 未改变的不更新
	if s, ok := envCacheMap.Load(key); ok {
		old = s.(string)
		if old == value {
			return
		}
	}

	envCacheMap.Store(key, value)

	__watchEnv_handleMap_lock.RLock()
	defer __watchEnv_handleMap_lock.RUnlock()

	for _, cb := range __watchEnv_handleMap[""] {
		cb(key, value, old)
	}
	for _, cb := range __watchEnv_handleMap[key] {
		cb(key, value, old)
	}
}

func GetenvInt(key string) int {
	return StringToInt(Getenv(key))
}
func GetenvUInt(key string) uint {
	return uint(GetenvInt(key))
}
func GetenvUint8(key string) uint8 {
	return uint8(GetenvInt(key))
}
func GetenvInt32(key string) int32 {
	return int32(GetenvInt(key))
}
func GetenvInt64(key string) int64 {
	return int64(GetenvInt(key))
}
func GetenvFloat32(key string) float32 {
	return StringToFloat32(Getenv(key))
}
func GetenvFloat64(key string) float64 {
	return StringToFloat64(Getenv(key))
}
func GetenvUInt32(key string) uint32 {
	return uint32(GetenvInt(key))
}
func GetenvUInt64(key string) uint64 {
	return uint64(GetenvInt(key))
}

func GetenvBool(key string) bool {
	return StringToBoolean(Getenv(key))
}

// 转数组
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

// 转 int32数组
func GetenvInt32Split(key string, split ...string) []int32 {
	arr := GetenvStringSplit(key, split...)
	newArr := []int32{}
	for _, v := range arr {
		newArr = append(newArr, StringToInt32(v))
	}

	return newArr
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

// 加载字符串中的环境变量
func LoadStringEnv(envStr string) *Error {
	envMaps, e := godotenv.Unmarshal(envStr)
	if e != nil {
		return NewError(e)
	}
	MergeEnv(envMaps)
	return nil

}

// 监听环境变量
// 首次调用会自动触发一次
// @key: 空字符串表示监听全部
func WatchEnvChange(key string, cb WatchEnvChangeHandle) {
	key = strings.ToUpper(key)
	__watchEnv_handleMap_lock.Lock()
	defer __watchEnv_handleMap_lock.Unlock()

	if key == "" {
		envCacheMap.Range(func(key, value any) bool {
			cb(key.(string), value.(string), "")
			return true
		})
	} else if val, ok := envCacheMap.Load(key); ok {
		cb(key, val.(string), "")
	}

	__watchEnv_handleMap[key] = append(__watchEnv_handleMap[key], cb)

}
