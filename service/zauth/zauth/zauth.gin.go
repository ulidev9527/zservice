package zauth

import (
	"zservice/service/zauth/internal"
	"zservice/zservice"

	"github.com/gin-gonic/gin"
)

// 授权检查
func GinCheckAuthMiddleware(zs *zservice.ZService) gin.HandlerFunc {
	return internal.GinMiddlewareCheckAuth(zs)
}
