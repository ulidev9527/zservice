package internal

import (
	"zservice/zserviceex/ginservice"

	"github.com/gin-gonic/gin"
)

var GinService *ginservice.GinService
var Gin *gin.Engine

func InitGin() {
	Gin.POST("/asset/configAsset", gin_post_configAsset)
	Gin.POST("/asset/upload", gin_post_upload)
	Gin.POST("/asset/zzzzString", gin_post_zzzzString)

	Gin.POST("/auth/login", gin_post_login)
	Gin.POST("/auth/logout", gin_post_loginout)

	Gin.POST("/sms/verifyCodeSend", gin_post_SMS_VerifyCodeSend)
	Gin.POST("/sms/verifyCodeVerify", gin_post_SMS_VerifyCodeVerify)

	Gin.GET("/permission/permission", gin_get_permission)
	Gin.PUT("/permission/permission", gin_put_permission)
	Gin.POST("/permission/permission", gin_post_permission)
	Gin.POST("/permission/bind", gin_post_permissionBind)

	Gin.GET("/permission/org", gin_get_GetOrgList)
	Gin.POST("/permission/org", gin_post_orgCreate)
	Gin.PUT("/permission/org", gin_put_orgUpdate)

	Gin.POST("/permission/userOrgBind", gin_post_userOrgBind)

}
