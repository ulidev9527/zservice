package internal

import (
	"net/http"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zserviceex/ginservice"

	"github.com/gin-gonic/gin"
)

func gin_get_orgListGet(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, Logic_OrgListGet(ginservice.GetCtxEX(ctx), &zauth_pb.OrgListGet_REQ{
		Page:   zservice.StringToUint32(ctx.Query("p")),
		Size:   zservice.StringToUint32(ctx.Query("si")),
		Search: ctx.Query("se"),
	}))
}
