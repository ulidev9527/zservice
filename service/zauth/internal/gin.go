package internal

import (
	"zservice/zservice/ex/ginservice"

	"github.com/gin-gonic/gin"
)

var GinService *ginservice.GinService
var Gin *gin.Engine

func InitGin() {
	initGinConfig()
	initGinLogin()
	initGinSms()
	initGinPermission()
}
