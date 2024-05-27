package zauth

import (
	"net/http"
	"strings"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice"
	"zservice/zservice/ex/ginservice"
	"zservice/zservice/zglobal"

	"github.com/gin-gonic/gin"
)

// 授权检查
func GinCheckAuthMiddleware(zs *zservice.ZService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		zctx := ginservice.GetCtxEX(ctx)
		zctx.AuthSign = zservice.MD5String(ctx.Request.UserAgent()) // 生成签名

		// 授权查询
		if e := CheckAuth(zctx, &zauth_pb.CheckAuth_REQ{
			Auth: string(zservice.JsonMustMarshal([]string{zservice.GetServiceName(), strings.ToLower(ctx.Request.Method), ctx.Request.URL.Path})),
		}); e != nil {

			ctx.JSON(http.StatusOK, &zglobal.Default_RES{
				Code: e.GetCode(),
				Msg:  zctx.TraceID,
			})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
