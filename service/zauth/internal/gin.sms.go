package internal

import "github.com/gin-gonic/gin"

func initGinSms() {
	Gin.POST("/sms/verifyCodeSend", gin_SMS_SendVerifyCode)
}

// 发送验证码
func gin_SMS_SendVerifyCode(ctx *gin.Context) {

}
