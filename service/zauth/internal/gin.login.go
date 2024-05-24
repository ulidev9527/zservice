package internal

import (
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice/ex/ginservice"
	"zservice/zservice/zglobal"

	"github.com/gin-gonic/gin"
)

func initGinLogin() {

	Gin.POST("/login", gin_Login)

}

type gin_T_Login struct {
	LoginType     uint   `json:"lt"`  // 登陆类型 1.手机短信 2.账号密码
	Phone         string `json:"p"`   // T1 手机号
	SMSVerifyCode string `json:"svc"` // T1 短信验证码
	LoginName     string `json:"ln"`  // T2 登陆名
	LoginPass     string `json:"lp"`  // T2 密码
}

func gin_Login(ctx *gin.Context) {
	zctx := ginservice.GetCtxEX(ctx)

	req := gin_T_Login{}

	if e := ctx.ShouldBind(&req); e != nil {
		zctx.LogError(e)
		ctx.JSON(200, gin.H{"code": zglobal.Code_ParamsErr})
		return
	}

	at, e := CreateToken(zctx)
	if e != nil {
		zctx.LogError(e)
		ctx.JSON(200, gin.H{"code": e.GetCode()})
		return
	}
	zctx.AuthToken = at.Token

	switch req.LoginType {
	case 1: // 手机号登陆
		res := Logic_LoginByPhone(zctx, &zauth_pb.LoginByPhone_REQ{
			Phone:       req.Phone,
			Expires:     uint32(zglobal.Time_10Day.Seconds()),
			LoginTarget: "zauth",
			VerifyCode:  req.SMSVerifyCode,
		})
		ctx.JSON(200, gin.H{"code": res.Code})
		return
	case 2: // 账号登陆
		res := Logic_LoginByAccount(zctx, &zauth_pb.LoginByAccount_REQ{
			Account:     req.LoginName,
			Password:    req.LoginPass,
			Expires:     uint32(zglobal.Time_10Day.Seconds()),
			LoginTarget: "zauth",
		})
		ctx.JSON(200, gin.H{"code": res.Code})
		return
	default:
		ctx.JSON(200, gin.H{"code": zglobal.Code_ParamsErr})
		return
	}
}
