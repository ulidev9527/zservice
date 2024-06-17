package internal

import (
	"net/http"
	"os"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/ex/ginservice"
	"zservice/zservice/zglobal"

	"github.com/gin-gonic/gin"
)

func gin_post_upload(ctx *gin.Context) {

	zctx := ginservice.GetCtxEX(ctx)
	file, e := ctx.FormFile("file")

	if e != nil {
		zctx.LogError(e)
		ctx.JSON(http.StatusOK, gin.H{"code": zglobal.Code_ErrorBreakoff})
		return
	}

	if file.Size > 30*1024*1024 { // 最大30M
		zctx.LogError("file size too big")
		ctx.JSON(http.StatusOK, gin.H{"code": zglobal.Code_Reject})
		return
	}

	filePath := zservice.GetTempFilepath()

	if e := ctx.SaveUploadedFile(file, filePath); e != nil {
		zctx.LogError(e)
		ctx.JSON(http.StatusOK, gin.H{"code": zglobal.Code_ErrorBreakoff})
		return
	}
	defer os.Remove(filePath)

	if bt, e := os.ReadFile(filePath); e != nil {
		zctx.LogError(e)
		ctx.JSON(http.StatusOK, gin.H{"code": zglobal.Code_ErrorBreakoff})
		return
	} else {
		ctx.JSON(http.StatusOK, Logic_AddAsset(zctx, &zauth_pb.AddAsset_REQ{
			Name:      file.Filename,
			FileBytes: bt,
		}))
	}

}
