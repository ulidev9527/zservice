package internal

import (
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
)

// 获取文件配置
func Logic_DownloadConfigAsset(ctx *zservice.Context, in *zauth_pb.DownloadConfigAsset_REQ) *zauth_pb.AssetInfo_RES {

	// 信息获取
	tab, e := GetConfigAssetInfo(ctx, in.Service, in.Name)
	if e != nil {
		ctx.LogError(e)
		return &zauth_pb.AssetInfo_RES{Code: zservice.Code_NotFound}
	}

	return Logic_DownloadAsset(ctx, &zauth_pb.DownloadAsset_REQ{AssetID: tab.AssetID})

}
