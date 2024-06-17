package internal

import (
	"fmt"
	"os"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/zglobal"
)

func Logic_AddAsset(ctx *zservice.Context, in *zauth_pb.AddAsset_REQ) *zauth_pb.AssetInfo_RES {

	md5Str := zservice.Md5Bytes(in.FileBytes)

	savePath := fmt.Sprintf(FI_UploadDir, md5Str)

	// 存储文件
	if _, e := os.Stat(savePath); e != nil {
		if os.IsNotExist(e) {
			if e := os.WriteFile(savePath, in.FileBytes, 0666); e != nil {
				ctx.LogError(e)
				return &zauth_pb.AssetInfo_RES{Code: zglobal.Code_ErrorBreakoff}
			}

		} else {
			ctx.LogError(e)
			return &zauth_pb.AssetInfo_RES{Code: zglobal.Code_ErrorBreakoff}
		}
	}

	// 创建资源信息
	tab, e := AssetCreate(ctx, &zauth_pb.AssetInfo{
		Name:   in.Name,
		Md5:    md5Str,
		Expire: in.Expire,
		Size:   uint64(len(in.FileBytes)),
	})
	if e != nil {
		ctx.LogError(e)
		return &zauth_pb.AssetInfo_RES{Code: e.GetCode()}
	} else {
		return &zauth_pb.AssetInfo_RES{Code: zglobal.Code_SUCC, Info: &zauth_pb.AssetInfo{
			Name:   tab.Name,
			Md5:    tab.MD5,
			Token:  tab.Token,
			Expire: tab.Expire,
			Size:   tab.Size,
		}}
	}

}
