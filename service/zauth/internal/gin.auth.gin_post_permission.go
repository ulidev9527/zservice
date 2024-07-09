package internal

import (
	"net/http"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"

	"github.com/gin-gonic/gin"
)

// 创建权限
func gin_post_permission(ctx *gin.Context) {

	zctx := GinService.GetCtx(ctx)
	req := &zauth_pb.PermissionInfo{}
	if e := ctx.ShouldBindJSON(req); e != nil {
		zctx.LogError(e)
		ctx.JSON(http.StatusOK, gin.H{"code": zservice.Code_Fail})
		return
	}

	ctx.JSON(http.StatusOK, Logic_PermissionCreate(zctx, req))

}
