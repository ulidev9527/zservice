package internal

import (
	"fmt"
	"net/http"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice/ex/ginservice"
	"zservice/zservice/zglobal"

	"github.com/gin-gonic/gin"
)

func gin_post_config_service_uploadEnvConfig(ctx *gin.Context) {
	zctx := ginservice.GetCtxEX(ctx)
	file, e := ctx.FormFile("file")
	if e != nil {
		zctx.LogError(e)
		ctx.JSON(http.StatusOK, gin.H{"code": zglobal.Code_Fail})
		return
	}
	serviceName := ctx.Param("service")
	filePath := fmt.Sprintf(FI_ServiceEnvFile, serviceName)
	ctx.SaveUploadedFile(file, filePath)

	ctx.JSON(http.StatusOK, Logic_ConfigSyncServiceEnvConfig(zctx, &zauth_pb.ConfigSyncServiceEnvConfig_REQ{
		Service:  serviceName,
		FilePath: filePath,
	}))
}
