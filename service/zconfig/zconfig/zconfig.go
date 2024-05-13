package zconfig

import (
	"encoding/json"
	"fmt"
	"sync"
	"zservice/service/zconfig/zconfig_pb"
	"zservice/zglobal"
	"zservice/zservice"

	"github.com/joho/godotenv"
	"github.com/nsqio/go-nsq"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type ZConfigConfig struct {
	Etcd            *clientv3.Client
	NsqConsumerAddr string // nsq consumer addr
	IsNsqd          bool
}

var grpcClient *GrpcClient
var fileConfigMap = &sync.Map{}

// 初始化
func Init(c *ZConfigConfig) {
	grpcClient = NewGrpcClient(c.Etcd)

	NewNsqConsumer_FileConfigChange(&NsqConsumerConfig{
		Addr:   c.NsqConsumerAddr,
		IsNsqd: c.IsNsqd,
		OnMessage: func(m *nsq.Message) error {
			fileName := string(m.Body)
			zservice.LogInfo("Update config ", fileName)
			fileConfigMap.Delete(fileName)
			return nil
		},
	})
}

// 加载远程环境变量
func LoadRemoteEnv(addr string, auth string) *zservice.Error {

	if addr == "" {
		return zservice.NewError("no addr")
	}

	body, e := zservice.Get(zservice.NewEmptyContext(), addr, &map[string]any{"auth": auth}, nil)
	if e != nil {
		return e
	}

	envMaps, _e := godotenv.UnmarshalBytes(body)
	if _e != nil {
		return zservice.NewError(_e)
	}

	zservice.MergeEnv(envMaps)
	return nil
}

// 获取指定文件的配置
// 不传 key 返回所有配置数组
// 一个 key 返回一个对象
// 多个 key 返回数组
func GetFileConfig(fileName string, v any, keys ...string) *zservice.Error {
	// 是否有配置，没有拉取配置
	val, has := fileConfigMap.Load(fileName)

	if !has { // 配置中心拉取配置
		res, e := grpcClient.GetFileConfig(zservice.NewEmptyContext(), &zconfig_pb.GetFileConfig_REQ{
			FileName: fileName,
		})
		if e != nil {
			return zservice.NewError(e)
		}

		if res.GetCode() != zglobal.Code_SUCC {
			return zservice.NewError("get file config fail").SetCode(res.Code)
		}

		var maps map[string]string
		ee := json.Unmarshal([]byte(res.Value), &maps)
		if ee != nil {
			return zservice.NewError(ee)
		}

		fileConfigMap.Store(fileName, maps)

		val = maps
	}

	// 开始解析
	maps := val.(map[string]string)
	keyCount := len(keys)

	if keyCount == 1 { // 解析一个配置
		str := maps[keys[0]]
		json.Unmarshal([]byte(str), v)
		return nil
	}

	useMap := map[string]string{}

	if keyCount == 0 {
		useMap = maps
	} else {
		for _, _v := range keys {
			str := maps[_v]
			if str == "" {
				return zservice.NewError("key not exist", _v).SetCode(zglobal.Code_Zconfig_GetConfigFail)
			}
			useMap[_v] = str
		}
	}

	// 数组模式解析
	jStr := ""
	for _, _v := range useMap {
		jStr += _v + ","
	}
	jStr = jStr[0 : len(jStr)-1]
	jStr = fmt.Sprintf("[ %s ]", jStr)
	zservice.LogInfo(jStr)
	if e := json.Unmarshal([]byte(jStr), v); e != nil {
		return zservice.NewError(e).SetCode(zglobal.Code_Zconfig_GetConfigFail)
	}
	return nil
}
