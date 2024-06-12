package internal

import (
	"net/http"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice/ex/ginservice"

	"github.com/gin-gonic/gin"
)

func gin_get_config_service_envCVonfig(ctx *gin.Context) {

	zctx := ginservice.GetCtxEX(ctx)
	serviceName := ctx.Param("service")
	ctx.JSON(http.StatusOK, Logic_ConfigGetEnvConfig(zctx, &zauth_pb.ConfigGetEnvConfig_REQ{
		Service: serviceName,
	}))
}
