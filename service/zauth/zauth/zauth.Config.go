package zauth

import (
	"encoding/json"
	"fmt"
	"zservice/service/zauth/internal"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/zglobal"
)

// 获取文件配置

// 获取指定文件的配置
// 不传 key 返回所有配置数组
// 一个 key 返回一个对象
// 多个 key 返回数组
func GetFileConfig(ctx *zservice.Context, fileName string, v any, keys ...string) *zservice.Error {
	// 是否有配置，没有拉取配置
	val, has := fileConfigMap.Load(fileName)

	if !has { // 配置中心拉取配置
		req := &zauth_pb.GetFileConfig_REQ{
			FileName: fileName,
		}

		res, e := func() (*zauth_pb.GetFileConfig_RES, error) {
			if zauthInitConfig.ZauthServiceName == zservice.GetServiceName() {
				return internal.Logic_GetFileConfig(ctx, req), nil
			}
			return grpcClient.GetFileConfig(ctx, req)
		}()
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