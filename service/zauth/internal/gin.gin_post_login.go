package internal

import (
	"net/http"
	"zservice/service/zauth/zauth_pb"
	"zservice/zservice/ex/ginservice"
	"zservice/zservice/zglobal"

	"github.com/gin-gonic/gin"
)

type gin_T_Login struct {
	LoginType     uint   `json:"lt"`  // 登陆类型 1.手机短信 2.账号密码
	Phone         string `json:"p"`   // T1 手机号
	SMSVerifyCode string `json:"svc"` // T1 短信验证码
	LoginName     string `json:"ln"`  // T2 登陆名
	LoginPass     string `json:"lp"`  // T2 密码
}

// 登陆
func gin_post_login(ctx *gin.Context) {
	zctx := ginservice.GetCtxEX(ctx)

	req := gin_T_Login{}

	if e := ctx.ShouldBind(&req); e != nil {
		zctx.LogError(e)
		ctx.JSON(http.StatusOK, gin.H{"code": zglobal.Code_ParamsErr})
		return
	}

	switch req.LoginType {
	case 1: // 手机号登陆
		res := Logic_LoginByPhone(zctx, &zauth_pb.LoginByPhone_REQ{
			Phone:      req.Phone,
			Expires:    uint32(zglobal.Time_10Day.Seconds()),
			VerifyCode: req.SMSVerifyCode,
		})

		if res.Code == zglobal.Code_SUCC {
			ginservice.SyncHeader(ctx)
		}

		ctx.JSON(http.StatusOK, gin.H{"code": res.Code})
		return
	case 2: // 账号登陆
		res := Logic_LoginByName(zctx, &zauth_pb.LoginByUser_REQ{
			User:     req.LoginName,
			Password: req.LoginPass,
			Expires:  uint32(zglobal.Time_10Day.Seconds()),
		})

		if res.Code == zglobal.Code_SUCC {
			ginservice.SyncHeader(ctx)
		}

		ctx.JSON(http.StatusOK, gin.H{"code": res.Code})
		return
	default:
		ctx.JSON(http.StatusOK, gin.H{"code": zglobal.Code_ParamsErr})
		return
	}
}
