package ginservice

import (
	"net/http"

	"github.com/ulidev9527/zservice/zservice"

	"github.com/gin-gonic/gin"
)

// 中间件
// CORS跨域中间件
func GinMiddlewareCORS(zs *zservice.ZService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		method := ctx.Request.Method
		origin := ctx.Request.Header.Get("Origin")
		if origin != "" {
			ctx.Header("Access-Control-Allow-Origin", "*") // 可将将 * 替换为指定的域名
			ctx.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			ctx.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization,"+zservice.S_C2S_Token)
			ctx.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			ctx.Header("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			ctx.AbortWithStatus(http.StatusNoContent)
		}
	}
}
