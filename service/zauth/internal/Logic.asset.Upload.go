package internal

import (
	"fmt"
	"os"
	"path/filepath"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
)

func Logic_UploadAsset(ctx *zservice.Context, in *zauth_pb.UploadAsset_REQ) *zauth_pb.AssetInfo_RES {

	if len(in.Data) > 10485760 { // 10MB
		ctx.LogError("file size too big:", len(in.Data))
		return &zauth_pb.AssetInfo_RES{Code: zservice.Code_Reject}
	}

	md5Str := zservice.Md5Bytes(in.Data)

	savePath := fmt.Sprintf(zservice.FI_UploadDir, md5Str)
	saveDir := filepath.Dir(savePath)

	// 文件夹创建
	if e := os.MkdirAll(saveDir, 0750); e != nil {
		ctx.LogError(e)
		return &zauth_pb.AssetInfo_RES{Code: zservice.Code_Fail}
	}

	// 存储文件
	if _, e := os.Stat(savePath); e != nil {
		if os.IsNotExist(e) {

			if e := os.WriteFile(savePath, in.Data, 0750); e != nil {
				ctx.LogError(e)
				return &zauth_pb.AssetInfo_RES{Code: zservice.Code_Fail}
			}

		} else {
			ctx.LogError(e)
			return &zauth_pb.AssetInfo_RES{Code: zservice.Code_Fail}
		}
	}

	// 创建资源信息
	tab, e := CreateAssetInfo(ctx, &zauth_pb.AssetInfo{
		Name:    in.Name,
		Size:    uint64(len(in.Data)),
		Md5:     md5Str,
		Expires: in.Expires,
	})
	if e != nil {
		ctx.LogError(e)
		return &zauth_pb.AssetInfo_RES{Code: e.GetCode()}
	} else {
		return &zauth_pb.AssetInfo_RES{Code: zservice.Code_SUCC, Info: tab.ToAssetInfo()}
	}
}
