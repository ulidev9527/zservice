package internal

import (
	"fmt"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/zglobal"
)

// 获取文件配置
func Logic_ConfigGetFileConfig(ctx *zservice.Context, in *zauth_pb.ConfigGetFileConfig_REQ) *zauth_pb.ConfigGetFileConfig_RES {
	fKey := fmt.Sprintf(RK_Config_ServiceFileConfig, in.Service, in.FileName) // 缓存 key

	if maps, e := Redis.HGetAll(fKey).Result(); e != nil {
		ctx.LogError(e)
		return &zauth_pb.ConfigGetFileConfig_RES{Code: zglobal.Code_NotFound}
	} else {
		return &zauth_pb.ConfigGetFileConfig_RES{Code: zglobal.Code_SUCC, Value: string(zservice.JsonMustMarshal(maps))}
	}

}
