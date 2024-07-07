package internal

import (
	"net/http"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zserviceex/ginservice"

	"github.com/gin-gonic/gin"
)

// 权限绑定
func gin_post_permissionBind(ctx *gin.Context) {

	zctx := ginservice.GetCtxEX(ctx)

	req := &zauth_pb.PermissionBind_REQ{}
	if e := ctx.ShouldBindJSON(req); e != nil {
		zctx.LogError(e)
		ctx.JSON(http.StatusOK, gin.H{"code": zservice.Code_Fail})
		return
	}

	ctx.JSON(http.StatusOK, Logic_PermissionBind(zctx, req))
}
