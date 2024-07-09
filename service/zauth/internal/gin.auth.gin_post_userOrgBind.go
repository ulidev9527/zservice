package internal

import (
	"net/http"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"

	"github.com/gin-gonic/gin"
)

func gin_post_userOrgBind(ctx *gin.Context) {

	zctx := GinService.GetCtx(ctx)

	req := &zauth_pb.UserOrgBind_REQ{}

	if e := ctx.ShouldBindJSON(req); e != nil {
		zctx.LogError(e)
		ctx.JSON(http.StatusOK, gin.H{"code": zservice.Code_Fail})
		return
	}

	ctx.JSON(http.StatusOK, Logic_UserOrgBind(zctx, req))

}
