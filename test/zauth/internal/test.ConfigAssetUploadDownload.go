package internal

import (
	"fmt"
	"os"
	"zservice/service/zauth/zauth"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
)

func Test_ConfigAssetUploadDownload() *zservice.ZService {
	ctx := zservice.NewContext()
	// 测试配置文件上传下载
	return zservice.NewService("config asset up test", func(z *zservice.ZService) {
		defer z.StartDone()

		filePath := "Config_Item.xlsx"
		bt, e := os.ReadFile(fmt.Sprint("static/", filePath))
		if e != nil {
			ctx.LogError(e)
			return
		}

		upRES := zauth.UploadConfigAsset(ctx, &zauth_pb.UploadConfigAsset_REQ{
			Name:   "Config_Item.xlsx",
			Parser: zservice.E_ConfigAsset_Parser_Excel,
			Data:   bt,
		})
		if upRES.Code != zservice.Code_SUCC {
			ctx.LogError("UploadAsset Fail:", upRES)
			return
		}

		downRES := zauth.DownloadAsset(ctx, &zauth_pb.DownloadAsset_REQ{
			AssetID: upRES.ConfigInfo.ParserAssetID,
		})
		if downRES.Code != zservice.Code_SUCC {
			ctx.LogError("DownloadAsset Fail:", upRES)
			return
		}

		ctx.LogInfo("DownloadAsset Success:", downRES)

		if e := os.WriteFile(fmt.Sprintf("static/down.%s.txt", downRES.Info.Name), downRES.Info.Data, 0750); e != nil {
			ctx.LogError(e)
			return
		}

		if config, e := GetConfig_Item(ctx, 1); e != nil {
			ctx.LogError(e)
		} else {
			ctx.LogInfo(config)
		}

	})
}
