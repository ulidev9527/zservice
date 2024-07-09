package internal

import (
	"net/http"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"

	"github.com/gin-gonic/gin"
)

func gin_get_GetOrgList(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, Logic_GetOrgList(GinService.GetCtx(ctx), &zauth_pb.GetOrgList_REQ{
		Page:   zservice.StringToUint32(ctx.Query("p")),
		Size:   zservice.StringToUint32(ctx.Query("si")),
		Search: ctx.Query("se"),
	}))
}
