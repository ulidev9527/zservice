package internal

import (
	"net/http"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice/service/ginservice"
	"zservice/zservice/zglobal"

	"github.com/gin-gonic/gin"
)

// 验证验证码
func gin_post_SMS_VerifyCodeVerify(ctx *gin.Context) {

	zctx := ginservice.GetCtxEX(ctx)

	req := gin_T_SMS_VerifyCodeSend{}

	if e := ctx.ShouldBind(&req); e != nil {
		zctx.LogError(e)
		ctx.JSON(http.StatusOK, gin.H{"code": zglobal.Code_ParamsErr})
		return
	}

	ctx.JSON(http.StatusOK, Logic_SMSVerifyCodeVerify(zctx, &zauth_pb.SMSVerifyCodeVerify_REQ{Phone: req.Phone, VerifyCode: req.VerifyCode}))
}
