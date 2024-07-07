package internal

import (
	"fmt"
	"os"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
)

func Logic_DownloadAsset(ctx *zservice.Context, in *zauth_pb.DownloadAsset_REQ) *zauth_pb.AssetInfo_RES {

	assetTab := &AssetTable{}

	if in.AssetID != "" {
		if tab, e := GetAssetByID(ctx, in.AssetID); e != nil {
			ctx.LogError(e)
			return &zauth_pb.AssetInfo_RES{Code: e.GetCode()}
		} else {
			assetTab = tab
		}
	} else {
		if caTab, e := GetConfigAssetInfo(ctx, in.Service, in.Name); e != nil {
			ctx.LogError(e)
			return &zauth_pb.AssetInfo_RES{Code: e.GetCode()}
		} else if tab, e := GetAssetByID(ctx, caTab.ParserAssetID); e != nil {
			ctx.LogError(e)
			return &zauth_pb.AssetInfo_RES{Code: e.GetCode()}
		} else {
			assetTab = tab
		}
	}
	result := &zauth_pb.AssetInfo_RES{
		Code: zservice.Code_SUCC,
		Info: assetTab.ToAssetInfo(),
	}

	if file, e := os.Open(fmt.Sprintf(zservice.FI_UploadDir, assetTab.MD5)); e != nil {
		ctx.LogError(e)
		return &zauth_pb.AssetInfo_RES{Code: zservice.Code_Fail}
	} else {
		defer file.Close()
		result.Info.Data = make([]byte, assetTab.Size)
		if i, e := file.Read(result.Info.Data); e != nil {
			ctx.LogError(i, e)
			return &zauth_pb.AssetInfo_RES{Code: zservice.Code_Fail}
		}
		return result
	}
}
