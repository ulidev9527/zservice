package internal

import (
	"zservice/zservice/ex/ginservice"

	"github.com/gin-gonic/gin"
)

var GinService *ginservice.GinService
var Gin *gin.Engine

func InitGin() {
	Gin.POST("/config/:service/uploadFileConfig", gin_post_config_service_uploadFileConfig)
	Gin.POST("/config/:service/uploadEnvConfig", gin_post_config_service_uploadEnvConfig)
	Gin.GET("/config/:service/envConfig", gin_get_config_service_envCVonfig)
	Gin.GET("/config/serviceEnvConfig/:auth", gin_get_config_ServiceEnvConfig_auth)

	Gin.POST("/login", gin_post_login)

	Gin.POST("/sms/verifyCodeSend", gin_post_SMS_VerifyCodeSend)
	Gin.POST("/sms/verifyCodeVerify", gin_post_SMS_VerifyCodeVerify)

	Gin.GET("/permission", gin_get_permission)
	Gin.PUT("/permission", gin_put_permission)
	Gin.POST("/permission", gin_post_permission)
	Gin.POST("/permission/bind", gin_post_permissionBind)

	Gin.GET("/org", gin_get_orgListGet)
	Gin.POST("/org", gin_post_orgCreate)
	Gin.PUT("/org", gin_put_orgUpdate)
}
