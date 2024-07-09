package internal

import (
	"net/http"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"

	"github.com/gin-gonic/gin"
)

// 获取权限
func gin_get_permission(ctx *gin.Context) {

	ctx.JSON(http.StatusOK, Logic_PermissionListGet(GinService.GetCtx(ctx), &zauth_pb.PermissionListGet_REQ{
		Page:   zservice.StringToUint32(ctx.Query("p")),
		Size:   zservice.StringToUint32(ctx.Query("si")),
		Search: ctx.Query("se"),
	}))

}
