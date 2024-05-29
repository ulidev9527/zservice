package internal

import (
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice/ex/ginservice"
	"zservice/zservice/zglobal"

	"github.com/gin-gonic/gin"
)

func initGinSms() {
	Gin.POST("/sms/verifyCodeSend", gin_SMS_SendVerifyCode)
	Gin.POST("/sms/verifyCodeVerify", gin_SMS_VerifyCodeVerify)
}

type gin_T_SMS_SendVerifyCode struct {
	Phone      string `json:"phone"`
	VerifyCode string `json:"vc"`
}

// 发送验证码
func gin_SMS_SendVerifyCode(ctx *gin.Context) {

	zctx := ginservice.GetCtxEX(ctx)

	req := gin_T_SMS_SendVerifyCode{}

	if e := ctx.ShouldBind(&req); e != nil {
		zctx.LogError(e)
		ctx.JSON(200, gin.H{"code": zglobal.Code_ParamsErr})
		return
	}

	ctx.JSON(200, Logic_SMSVerifyCodeSend(zctx, &zauth_pb.SMSVerifyCodeSend_REQ{Phone: req.Phone, VerifyCode: req.VerifyCode}))

}

// 验证验证码
func gin_SMS_VerifyCodeVerify(ctx *gin.Context) {

	zctx := ginservice.GetCtxEX(ctx)

	req := gin_T_SMS_SendVerifyCode{}

	if e := ctx.ShouldBind(&req); e != nil {
		zctx.LogError(e)
		ctx.JSON(200, gin.H{"code": zglobal.Code_ParamsErr})
		return
	}

	ctx.JSON(200, Logic_SMSVerifyCodeVerify(zctx, &zauth_pb.SMSVerifyCodeVerify_REQ{Phone: req.Phone, VerifyCode: req.VerifyCode}))
}
