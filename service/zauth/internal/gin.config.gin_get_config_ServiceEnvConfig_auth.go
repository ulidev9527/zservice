package internal

import (
	"net/http"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice/service/ginservice"

	"github.com/gin-gonic/gin"
)

func gin_get_config_ServiceEnvConfig_auth(ctx *gin.Context) {
	zctx := ginservice.GetCtxEX(ctx)
	auth := ctx.Param("auth")

	ctx.JSON(http.StatusOK, Logic_ConfigGetServiceEnvConfig(zctx, &zauth_pb.ConfigGetServiceEnvConfig_REQ{
		Auth: auth,
	}))
}
