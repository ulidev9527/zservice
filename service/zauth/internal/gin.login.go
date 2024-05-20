package internal

import (
	"zservice/zservice/ex/ginservice"

	"github.com/gin-gonic/gin"
)

func initGinLogin() {
	Gin.POST("/login", ginservice.GinAuthEXMiddleware(GinService.ZService), func(ctx *gin.Context) {

	})
}
