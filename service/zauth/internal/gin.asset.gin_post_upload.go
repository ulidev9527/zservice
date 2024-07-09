package internal

import (
	"net/http"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zserviceex/ginservice"

	"github.com/gin-gonic/gin"
)

func gin_post_upload(ctx *gin.Context) {

	zctx := GinService.GetCtx(ctx)
	file, e := ctx.FormFile("file")

	if e != nil {
		zctx.LogError(e)
		ctx.JSON(http.StatusOK, gin.H{"code": zservice.Code_Fail})
		return
	}

	if file.Size > 10485760 { // 10MB
		zctx.LogError("file size too big:", file.Size)
		ctx.JSON(http.StatusOK, gin.H{"code": zservice.Code_Reject})
		return
	}

	if bt, e := ginservice.ReadUploadFile(file); e != nil {
		zctx.LogError(e)
		ctx.JSON(http.StatusOK, gin.H{"code": zservice.Code_Fail})
		return
	} else {
		ctx.JSON(http.StatusOK, Logic_UploadAsset(zctx, &zauth_pb.UploadAsset_REQ{
			Name:    file.Filename,
			Expires: zservice.StringToInt64(ctx.PostForm("expires")),
			Data:    bt,
		}))
	}
}
