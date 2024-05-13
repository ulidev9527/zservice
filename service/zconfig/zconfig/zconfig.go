package zconfig

import (
	"encoding/json"
	"sync"
	"zservice/service/zconfig/zconfig_pb"
	"zservice/zglobal"
	"zservice/zservice"

	"github.com/joho/godotenv"
	"github.com/nsqio/go-nsq"
	clientv3 "go.etcd.io/etcd/client/v3"
)

type ZConfigConfig struct {
	Etcd     *clientv3.Client
	NsqAddrs []string
	IsNsqd   bool
}

var grpcClient *GrpcClient
var fileConfigMap = &sync.Map{}

// 初始化
func Init(c *ZConfigConfig) {
	grpcClient = NewGrpcClient(c.Etcd)

	NewNsqConsumer_FileConfigChange(&NsqConsumerConfig{
		Addrs:  c.NsqAddrs,
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
// 不传 key 返回所有配置 map
// 一个 key 返回一个对象
// 多个 key 返回数组
func GetFileConfig(fileName string, v any, keys ...string) *zservice.Error {
	// 是否有配置，没有拉取配置
	val, has := fileConfigMap.Load(fileName)

	if !has {
		res, e := grpcClient.GetFileConfig(zservice.NewEmptyContext(), &zconfig_pb.GetFileConfig_REQ{
			FileName: fileName,
		})
		if e != nil {
			return zservice.NewError(e)
		}

		if res.GetCode() != zglobal.Code_SUCC {
			return zservice.NewError("get file config fail").SetCode(res.Code)
		}

		var maps map[string]map[string]any
		ee := json.Unmarshal([]byte(res.Value), &maps)
		if ee != nil {
			return zservice.NewError(ee)
		}

		fileConfigMap.Store(fileName, maps)

		val = maps
	}

	maps := val.(map[string]map[string]any)
	keyCount := len(keys)
	if keyCount == 0 {
		arr := make([]map[string]any, 0)
		for _, _v := range maps {
			arr = append(arr, _v)
		}
		json.Unmarshal(zservice.JsonMustMarshal(arr), v)
		return nil
	} else if keyCount == 1 {
		json.Unmarshal(zservice.JsonMustMarshal(maps[keys[0]]), v)
		return nil
	} else {
		arr := make([]any, 0)
		for _, _v := range keys {
			m := maps[_v]
			if m == nil {
				return zservice.NewError("key not exist", _v).SetCode(zglobal.Code_Zconfig_GetConfigFail)
			}
			arr = append(arr, maps[_v])
		}
		json.Unmarshal(zservice.JsonMustMarshal(arr), v)
		return nil
	}
}
