package zauth

import (
	"encoding/json"
	"fmt"
	"sync"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/zglobal"
)

// 获取指定文件的配置
// 不传 key 返回所有配置数组
// 一个 key 返回一个对象
// 多个 key 返回数组
func ConfigGetFileConfig(ctx *zservice.Context, fileName string, item any, keys ...string) *zservice.Error {
	// 是否有配置，没有拉取配置
	configMap := &sync.Map{}
	if c, has := fileConfigCache.Load(fileName); !has {
		res, e := grpcClient.ConfigGetFileConfig(ctx, &zauth_pb.ConfigGetFileConfig_REQ{
			FileName: fileName,
			Service:  zservice.GetServiceName(),
		})
		if e != nil {
			return zservice.NewError(e)
		}

		if res.Code != zglobal.Code_SUCC {
			return zservice.NewError("get config fail:", fileName).SetCode(res.Code)
		}

		maps := &map[string]string{}
		if e := json.Unmarshal([]byte(res.Value), maps); e != nil {
			return zservice.NewError(e)
		}

		for k, v := range *maps {
			configMap.Store(k, v)
		}

		fileConfigCache.Store(fileName, configMap)
	} else {
		configMap = c.(*sync.Map)
	}

	switch len(keys) {
	case 0:
		str := ""
		configMap.Range(func(key, value any) bool {
			str += value.(string) + ","
			return true
		})

		str = str[0 : len(str)-1] // 去掉尾部逗号 `,`
		str = fmt.Sprintf("[ %s ]", str)
		if e := json.Unmarshal([]byte(str), item); e != nil {
			return zservice.NewError(e).SetCode(zglobal.Code_NotFound)
		}
		return nil

	case 1:
		if str, has := configMap.Load(keys[0]); !has {
			return zservice.NewError("not found").SetCode(zglobal.Code_NotFound)
		} else {
			if e := json.Unmarshal([]byte(str.(string)), item); e != nil {
				return zservice.NewError(e)
			}
			return nil
		}
	default:
		str := ""
		for _, key := range keys {
			s, has := configMap.Load(key)
			if !has {
				continue
			}
			str += s.(string) + ","
		}

		str = str[0 : len(str)-1] // 去掉尾部逗号 `,`
		str = fmt.Sprintf("[ %s ]", str)
		if e := json.Unmarshal([]byte(str), item); e != nil {
			return zservice.NewError(e).SetCode(zglobal.Code_NotFound)
		}
		return nil
	}
}
