package internal

import (
	"net/http"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice/ex/ginservice"

	"github.com/gin-gonic/gin"
)

func gin_post_loginout(ctx *gin.Context) {
	zctx := ginservice.GetCtxEX(ctx)

	ctx.JSON(http.StatusOK, Logic_Logout(zctx, &zauth_pb.Default_REQ{}))
}
