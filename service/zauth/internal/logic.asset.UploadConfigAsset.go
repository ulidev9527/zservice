package internal

import (
	"fmt"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
)

// 同步服务配置
func Logic_UploadConfigAsset(ctx *zservice.Context, in *zauth_pb.UploadConfigAsset_REQ) *zauth_pb.UploadConfigAsset_RES {

	// 参数验证
	if in.Service == "" || len(in.Data) == 0 {
		ctx.LogError("param error", in.Service, in.Parser, len(in.Data))
		return &zauth_pb.UploadConfigAsset_RES{Code: zservice.Code_ParamsErr}
	}

	// 上传
	assetInfoRES := Logic_UploadAsset(ctx, &zauth_pb.UploadAsset_REQ{
		Name: in.Name,
		Data: in.Data,
	})
	if assetInfoRES.Code != zservice.Code_SUCC {
		return &zauth_pb.UploadConfigAsset_RES{Code: assetInfoRES.Code}
	}

	// 解析配置
	parser := ConfigAssetParserMap[in.Parser] // 获取解析器
	if parser == nil {
		ctx.LogError("parser not found", in.Parser)
		return &zauth_pb.UploadConfigAsset_RES{Code: zservice.Code_Zauth_config_ParserNotExist}
	}

	parserMaps, e := parser(fmt.Sprintf(zservice.FI_UploadDir, assetInfoRES.Info.Md5)) // 解析
	if e != nil {
		ctx.LogError(e)
		return &zauth_pb.UploadConfigAsset_RES{Code: zservice.Code_Zauth_config_ParserFail}
	}

	// 解析后上传到资源表
	parserAssetInfoRES := Logic_UploadAsset(ctx, &zauth_pb.UploadAsset_REQ{
		Name: fmt.Sprintf("%s.%s.%d", in.Service, in.Name, in.Parser),
		Data: zservice.JsonMustMarshal(parserMaps),
	})
	if parserAssetInfoRES.Code != zservice.Code_SUCC {
		return &zauth_pb.UploadConfigAsset_RES{Code: parserAssetInfoRES.Code}
	}

	// 存储配置信息
	tab, e := CreateConfigAssetInfo(ctx, &zauth_pb.ConfigAssetInfo{
		Name:          in.Name,
		Service:       in.Service,
		Parser:        in.Parser,
		AssetID:       assetInfoRES.Info.AssetID,
		ParserAssetID: parserAssetInfoRES.Info.AssetID,
	})
	if e != nil {
		ctx.LogError(e)
		return &zauth_pb.UploadConfigAsset_RES{Code: e.GetCode()}
	}

	// 事件推送
	zservice.Go(func() {
		if e := EtcdService.SendEvent(ctx, fmt.Sprintf(zservice.EV_Config_ServiceFileConfigChange, in.Service), in.Name); e != nil {
			ctx.LogError(e)
		}
	})

	return &zauth_pb.UploadConfigAsset_RES{
		Code:            zservice.Code_SUCC,
		AssetInfo:       assetInfoRES.Info,
		ParserAssetInfo: parserAssetInfoRES.Info,
		ConfigInfo:      tab.ToConfigAssetInfo(),
	}
}
