package internal

import (
	"fmt"
	"path"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/zglobal"

	"github.com/redis/go-redis/v9"
)

// 同步服务配置
func Logic_ConfigSyncServiceFileConfig(ctx *zservice.Context, in *zauth_pb.ConfigSyncServiceFileConfig_REQ) *zauth_pb.Default_RES {

	if in.Service == "" || in.FilePath == "" {
		return &zauth_pb.Default_RES{Code: zglobal.Code_ParamsErr}
	}

	// 验证文件
	if e := parserFileVerify(in.FilePath); e != nil {
		ctx.LogError(e)
		return &zauth_pb.Default_RES{Code: e.GetCode()}
	}

	// 提取文件名
	fileName := path.Base(in.FilePath)

	// 验证 md5
	md5Str := ""
	if str, e := zservice.Md5File(in.FilePath); e != nil {
		ctx.LogError(e)
		return &zauth_pb.Default_RES{Code: zglobal.Code_Zauth_config_GetFileMd5Fail}
	} else {
		md5Str = str
	}
	rk_md5 := fmt.Sprintf(RK_Config_ServiceFileConfigMD5, in.Service, fileName)
	if s, e := Redis.Get(rk_md5).Result(); e != nil {
		if e != redis.Nil {
			ctx.LogError(e)
			return &zauth_pb.Default_RES{Code: zglobal.Code_ErrorBreakoff}
		}
	} else if s == md5Str {
		return &zauth_pb.Default_RES{Code: zglobal.Code_Zauth_config_FileMd5NotChange}
	}

	// 解析文件
	parser := fileParserMap[in.Parser]
	if parser == nil {
		ctx.LogError("parser not found", in.Parser)
		return &zauth_pb.Default_RES{Code: zglobal.Code_Zauth_config_ParserNotExist}
	}

	if maps, e := parser(in.FilePath); e != nil {
		ctx.LogError(e)
		return &zauth_pb.Default_RES{Code: zglobal.Code_Zauth_config_ParserFail}
	} else if e := Redis.HSet(fmt.Sprintf(RK_Config_ServiceFileConfig, in.Service, fileName), maps).Err(); e != nil { // 存储 redis
		ctx.LogError(e)
		return &zauth_pb.Default_RES{Code: zglobal.Code_ErrorBreakoff}
	} else {

		// 存储 md5 信息
		if e := Redis.Set(rk_md5, md5Str).Err(); e != nil {
			ctx.LogError(e)
		}

		// 事件推送
		EV_Send_Config_serviceFileConfigChange(ctx, Etcd, in.Service, fileName)

		return &zauth_pb.Default_RES{Code: zglobal.Code_SUCC}
	}
}
