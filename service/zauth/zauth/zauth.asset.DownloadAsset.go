package zauth

import (
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
)

func DownloadAsset(ctx *zservice.Context, in *zauth_pb.DownloadAsset_REQ) *zauth_pb.AssetInfo_RES {
	if res, e := grpcClient.DownloadAsset(ctx, in); e != nil {
		ctx.LogError(e)
		return &zauth_pb.AssetInfo_RES{Code: zservice.Code_500Err}
	} else {
		return res
	}
}
