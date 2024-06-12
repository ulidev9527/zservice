package internal

import (
	"net/http"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice/ex/ginservice"
	"zservice/zservice/zglobal"

	"github.com/gin-gonic/gin"
)

type gin_T_SMS_VerifyCodeSend struct {
	Phone      string `json:"phone"`
	VerifyCode string `json:"vc"`
}

// 发送验证码
func gin_post_SMS_VerifyCodeSend(ctx *gin.Context) {

	zctx := ginservice.GetCtxEX(ctx)

	req := gin_T_SMS_VerifyCodeSend{}

	if e := ctx.ShouldBind(&req); e != nil {
		zctx.LogError(e)
		ctx.JSON(http.StatusOK, gin.H{"code": zglobal.Code_ParamsErr})
		return
	}

	ctx.JSON(http.StatusOK, Logic_SMSVerifyCodeSend(zctx, &zauth_pb.SMSVerifyCodeSend_REQ{Phone: req.Phone, VerifyCode: req.VerifyCode}))

}
